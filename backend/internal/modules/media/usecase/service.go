package usecase

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"path"
	"strings"
	"time"

	"crm-starterkit/internal/core/storage"
	"crm-starterkit/internal/modules/media/domain"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("media not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrNoStorage    = errors.New("storage unavailable")
)

type Service struct {
	repo    domain.Repository
	storage storage.Driver
}

func NewService(repo domain.Repository, storageDriver storage.Driver) *Service {
	return &Service{repo: repo, storage: storageDriver}
}

func (s *Service) Upload(ctx context.Context, ownerID, filename, contentType string, reader io.Reader, size int64) (*domain.Media, error) {
	if s.storage == nil {
		return nil, ErrNoStorage
	}
	if filename == "" {
		return nil, ErrInvalidInput
	}
	id := uuid.NewString()
	key := path.Join("media", ownerID, id+"_"+path.Base(filename))
	obj, err := s.storage.Put(ctx, key, reader, size, contentType)
	if err != nil {
		return nil, err
	}
	media := &domain.Media{
		ID: id, OwnerID: ownerID, Key: obj.Key, URL: obj.URL, ContentType: contentType,
		Size: obj.Size, CreatedAt: time.Now().UTC(),
	}
	if strings.HasPrefix(contentType, "image/") {
		thumb, err := s.GenerateThumbnail(ctx, media)
		if err == nil && thumb != nil {
			media.ThumbnailKey = thumb.ThumbnailKey
			media.ThumbnailURL = thumb.ThumbnailURL
		}
	}
	if err := s.repo.Create(ctx, media); err != nil {
		return nil, err
	}
	return media, nil
}

func (s *Service) GenerateThumbnail(ctx context.Context, media *domain.Media) (*domain.Media, error) {
	if s.storage == nil || media == nil {
		return nil, ErrNoStorage
	}
	rc, err := s.storage.Get(ctx, media.Key)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	img, _, err := image.Decode(rc)
	if err != nil {
		return media, nil
	}
	const max = 200
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w > h {
		h = h * max / w
		w = max
	} else {
		w = w * max / h
		h = max
	}
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, b, draw.Over, nil)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}
	thumbKey := media.Key + ".thumb.jpg"
	obj, err := s.storage.Put(ctx, thumbKey, &buf, int64(buf.Len()), "image/jpeg")
	if err != nil {
		return nil, err
	}
	media.ThumbnailKey = obj.Key
	media.ThumbnailURL = obj.URL
	return media, nil
}

func (s *Service) Delete(ctx context.Context, id, ownerID string) error {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if ownerID != "" && m.OwnerID != ownerID {
		return ErrNotFound
	}
	if s.storage != nil {
		_ = s.storage.Delete(ctx, m.Key)
		if m.ThumbnailKey != "" {
			_ = s.storage.Delete(ctx, m.ThumbnailKey)
		}
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) ListMine(ctx context.Context, ownerID string) ([]domain.Media, error) {
	return enrichFilenames(s.repo.ListByOwner(ctx, ownerID))
}

func (s *Service) ListAll(ctx context.Context) ([]domain.Media, error) {
	return enrichFilenames(s.repo.ListAll(ctx))
}

func enrichFilenames(items []domain.Media, err error) ([]domain.Media, error) {
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].Filename = path.Base(items[i].Key)
	}
	return items, nil
}
