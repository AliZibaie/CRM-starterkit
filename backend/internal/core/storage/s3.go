package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"crm-starterkit/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Driver struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewS3Driver(cfg config.StorageConfig) (*S3Driver, error) {
	if cfg.S3.Bucket == "" {
		return nil, fmt.Errorf("AWS_S3_BUCKET is required for s3 storage")
	}

	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.S3.Region),
	}
	if cfg.S3.AccessKeyID != "" && cfg.S3.SecretAccessKey != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.S3.AccessKeyID, cfg.S3.SecretAccessKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
		}
		o.UsePathStyle = cfg.S3.ForcePathStyle
	})

	publicURL := strings.TrimRight(cfg.S3.PublicBaseURL, "/")
	if publicURL == "" {
		publicURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.S3.Bucket, cfg.S3.Region)
	}

	return &S3Driver{
		client:    client,
		bucket:    cfg.S3.Bucket,
		publicURL: publicURL,
	}, nil
}

func (d *S3Driver) Name() string {
	return "s3"
}

func (d *S3Driver) Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*Object, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	if size > 0 {
		input.ContentLength = aws.Int64(size)
	}
	if _, err := d.client.PutObject(ctx, input); err != nil {
		return nil, err
	}
	return &Object{
		Key:         key,
		ContentType: contentType,
		Size:        size,
		URL:         d.URL(key),
	}, nil
}

func (d *S3Driver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := d.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (d *S3Driver) Delete(ctx context.Context, key string) error {
	_, err := d.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (d *S3Driver) URL(key string) string {
	return fmt.Sprintf("%s/%s", d.publicURL, strings.TrimLeft(key, "/"))
}
