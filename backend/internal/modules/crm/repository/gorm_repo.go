package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/crm/domain"

	"gorm.io/gorm"
)

type LeadModel struct {
	ID, Name, Email, Phone, Status, OwnerID, Source string
	CreatedAt, UpdatedAt                            time.Time
}

func (LeadModel) TableName() string { return "crm_leads" }

type ContactModel struct {
	ID, LeadID, UserID, Name, Email, Phone, Company, OwnerID string
	CreatedAt, UpdatedAt                                      time.Time
}

func (ContactModel) TableName() string { return "crm_contacts" }

type PipelineModel struct {
	ID, Name, StagesJSON string
	CreatedAt            time.Time
}

func (PipelineModel) TableName() string { return "crm_pipelines" }

type DealModel struct {
	ID, Title, ContactID, PipelineID, Stage, OwnerID string
	Amount                                           int64
	CreatedAt, UpdatedAt                             time.Time
}

func (DealModel) TableName() string { return "crm_deals" }

type ActivityModel struct {
	ID, EntityType, EntityID, Type, Body, ActorID string
	CreatedAt                                     time.Time
}

func (ActivityModel) TableName() string { return "crm_activities" }

type TaskModel struct {
	ID, Title, EntityType, EntityID, OwnerID, Status string
	DueAt, ReminderAt                                *time.Time
	CreatedAt, UpdatedAt                             time.Time
}

func (TaskModel) TableName() string { return "crm_tasks" }

type TagModel struct {
	ID, Name string
}

func (TagModel) TableName() string { return "crm_tags" }

type EntityTagModel struct {
	EntityType string `gorm:"primaryKey;size:50"`
	EntityID   string `gorm:"primaryKey;size:36"`
	TagID      string `gorm:"primaryKey;size:36"`
}

func (EntityTagModel) TableName() string { return "crm_entity_tags" }

type SegmentModel struct {
	ID, Name, RulesJSON string
	CreatedAt           time.Time
}

func (SegmentModel) TableName() string { return "crm_segments" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(
		&LeadModel{}, &ContactModel{}, &PipelineModel{}, &DealModel{},
		&ActivityModel{}, &TaskModel{}, &TagModel{}, &EntityTagModel{}, &SegmentModel{},
	)
}

func (s *Store) CreateLead(ctx context.Context, l *domain.Lead) error {
	return s.db.WithContext(ctx).Create(&LeadModel{ID: l.ID, Name: l.Name, Email: l.Email, Phone: l.Phone, Status: l.Status, OwnerID: l.OwnerID, Source: l.Source, CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt}).Error
}
func (s *Store) UpdateLead(ctx context.Context, l *domain.Lead) error {
	return s.db.WithContext(ctx).Save(&LeadModel{ID: l.ID, Name: l.Name, Email: l.Email, Phone: l.Phone, Status: l.Status, OwnerID: l.OwnerID, Source: l.Source, CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt}).Error
}
func (s *Store) FindLead(ctx context.Context, id string) (*domain.Lead, error) {
	var m LeadModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Lead{ID: m.ID, Name: m.Name, Email: m.Email, Phone: m.Phone, Status: m.Status, OwnerID: m.OwnerID, Source: m.Source, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}
func (s *Store) ListLeads(ctx context.Context) ([]domain.Lead, error) {
	var rows []LeadModel
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Lead, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Lead{ID: m.ID, Name: m.Name, Email: m.Email, Phone: m.Phone, Status: m.Status, OwnerID: m.OwnerID, Source: m.Source, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return out, nil
}

func (s *Store) CreateContact(ctx context.Context, c *domain.Contact) error {
	return s.db.WithContext(ctx).Create(&ContactModel{ID: c.ID, LeadID: c.LeadID, UserID: c.UserID, Name: c.Name, Email: c.Email, Phone: c.Phone, Company: c.Company, OwnerID: c.OwnerID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}).Error
}
func (s *Store) UpdateContact(ctx context.Context, c *domain.Contact) error {
	return s.db.WithContext(ctx).Save(&ContactModel{ID: c.ID, LeadID: c.LeadID, UserID: c.UserID, Name: c.Name, Email: c.Email, Phone: c.Phone, Company: c.Company, OwnerID: c.OwnerID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}).Error
}
func (s *Store) FindContact(ctx context.Context, id string) (*domain.Contact, error) {
	var m ContactModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Contact{ID: m.ID, LeadID: m.LeadID, UserID: m.UserID, Name: m.Name, Email: m.Email, Phone: m.Phone, Company: m.Company, OwnerID: m.OwnerID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}
func (s *Store) FindContactByEmail(ctx context.Context, email string) (*domain.Contact, error) {
	var m ContactModel
	if err := s.db.WithContext(ctx).First(&m, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &domain.Contact{ID: m.ID, LeadID: m.LeadID, UserID: m.UserID, Name: m.Name, Email: m.Email, Phone: m.Phone, Company: m.Company, OwnerID: m.OwnerID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}
func (s *Store) ListContacts(ctx context.Context) ([]domain.Contact, error) {
	var rows []ContactModel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Contact, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Contact{ID: m.ID, LeadID: m.LeadID, UserID: m.UserID, Name: m.Name, Email: m.Email, Phone: m.Phone, Company: m.Company, OwnerID: m.OwnerID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return out, nil
}

func (s *Store) UpsertPipeline(ctx context.Context, p *domain.Pipeline) error {
	return s.db.WithContext(ctx).Save(&PipelineModel{ID: p.ID, Name: p.Name, StagesJSON: p.StagesJSON, CreatedAt: p.CreatedAt}).Error
}
func (s *Store) ListPipelines(ctx context.Context) ([]domain.Pipeline, error) {
	var rows []PipelineModel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Pipeline, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Pipeline{ID: m.ID, Name: m.Name, StagesJSON: m.StagesJSON, CreatedAt: m.CreatedAt})
	}
	return out, nil
}

func (s *Store) CreateDeal(ctx context.Context, d *domain.Deal) error {
	return s.db.WithContext(ctx).Create(&DealModel{ID: d.ID, Title: d.Title, ContactID: d.ContactID, PipelineID: d.PipelineID, Stage: d.Stage, OwnerID: d.OwnerID, Amount: d.Amount, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}).Error
}
func (s *Store) UpdateDeal(ctx context.Context, d *domain.Deal) error {
	return s.db.WithContext(ctx).Save(&DealModel{ID: d.ID, Title: d.Title, ContactID: d.ContactID, PipelineID: d.PipelineID, Stage: d.Stage, OwnerID: d.OwnerID, Amount: d.Amount, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}).Error
}
func (s *Store) FindDeal(ctx context.Context, id string) (*domain.Deal, error) {
	var m DealModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Deal{ID: m.ID, Title: m.Title, ContactID: m.ContactID, PipelineID: m.PipelineID, Stage: m.Stage, OwnerID: m.OwnerID, Amount: m.Amount, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}
func (s *Store) ListDealsByContact(ctx context.Context, contactID string) ([]domain.Deal, error) {
	var rows []DealModel
	if err := s.db.WithContext(ctx).Where("contact_id = ?", contactID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Deal, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Deal{ID: m.ID, Title: m.Title, ContactID: m.ContactID, PipelineID: m.PipelineID, Stage: m.Stage, OwnerID: m.OwnerID, Amount: m.Amount, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return out, nil
}

func (s *Store) ListDealsByPipeline(ctx context.Context, pipelineID string) ([]domain.Deal, error) {
	var rows []DealModel
	if err := s.db.WithContext(ctx).Where("pipeline_id = ?", pipelineID).Order("updated_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Deal, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Deal{ID: m.ID, Title: m.Title, ContactID: m.ContactID, PipelineID: m.PipelineID, Stage: m.Stage, OwnerID: m.OwnerID, Amount: m.Amount, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return out, nil
}

func (s *Store) DeleteContact(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(&ContactModel{}, "id = ?", id).Error
}

func (s *Store) CreateActivity(ctx context.Context, a *domain.Activity) error {
	return s.db.WithContext(ctx).Create(&ActivityModel{ID: a.ID, EntityType: a.EntityType, EntityID: a.EntityID, Type: a.Type, Body: a.Body, ActorID: a.ActorID, CreatedAt: a.CreatedAt}).Error
}
func (s *Store) ListActivities(ctx context.Context, entityType, entityID string) ([]domain.Activity, error) {
	var rows []ActivityModel
	if err := s.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Activity, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Activity{ID: m.ID, EntityType: m.EntityType, EntityID: m.EntityID, Type: m.Type, Body: m.Body, ActorID: m.ActorID, CreatedAt: m.CreatedAt})
	}
	return out, nil
}

func (s *Store) CreateTask(ctx context.Context, t *domain.Task) error {
	return s.db.WithContext(ctx).Create(&TaskModel{ID: t.ID, Title: t.Title, EntityType: t.EntityType, EntityID: t.EntityID, OwnerID: t.OwnerID, Status: t.Status, DueAt: t.DueAt, ReminderAt: t.ReminderAt, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}).Error
}
func (s *Store) UpdateTask(ctx context.Context, t *domain.Task) error {
	return s.db.WithContext(ctx).Save(&TaskModel{ID: t.ID, Title: t.Title, EntityType: t.EntityType, EntityID: t.EntityID, OwnerID: t.OwnerID, Status: t.Status, DueAt: t.DueAt, ReminderAt: t.ReminderAt, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}).Error
}
func (s *Store) ListTasksByOwner(ctx context.Context, ownerID string) ([]domain.Task, error) {
	var rows []TaskModel
	if err := s.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Task, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Task{ID: m.ID, Title: m.Title, EntityType: m.EntityType, EntityID: m.EntityID, OwnerID: m.OwnerID, Status: m.Status, DueAt: m.DueAt, ReminderAt: m.ReminderAt, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return out, nil
}
func (s *Store) ListDueReminders(ctx context.Context, before time.Time) ([]domain.Task, error) {
	var rows []TaskModel
	if err := s.db.WithContext(ctx).Where("status = ? AND reminder_at IS NOT NULL AND reminder_at <= ?", "open", before).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Task, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Task{ID: m.ID, Title: m.Title, EntityType: m.EntityType, EntityID: m.EntityID, OwnerID: m.OwnerID, Status: m.Status, DueAt: m.DueAt, ReminderAt: m.ReminderAt, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return out, nil
}

func (s *Store) UpsertTag(ctx context.Context, t *domain.Tag) error {
	return s.db.WithContext(ctx).Save(&TagModel{ID: t.ID, Name: t.Name}).Error
}
func (s *Store) AttachTag(ctx context.Context, entityType, entityID, tagID string) error {
	return s.db.WithContext(ctx).Save(&EntityTagModel{EntityType: entityType, EntityID: entityID, TagID: tagID}).Error
}
func (s *Store) ListTags(ctx context.Context, entityType, entityID string) ([]string, error) {
	var names []string
	err := s.db.WithContext(ctx).Table("crm_tags").
		Select("crm_tags.name").
		Joins("JOIN crm_entity_tags ON crm_entity_tags.tag_id = crm_tags.id").
		Where("crm_entity_tags.entity_type = ? AND crm_entity_tags.entity_id = ?", entityType, entityID).
		Pluck("crm_tags.name", &names).Error
	return names, err
}

func (s *Store) UpsertSegment(ctx context.Context, seg *domain.Segment) error {
	return s.db.WithContext(ctx).Save(&SegmentModel{ID: seg.ID, Name: seg.Name, RulesJSON: seg.RulesJSON, CreatedAt: seg.CreatedAt}).Error
}
func (s *Store) ListSegments(ctx context.Context) ([]domain.Segment, error) {
	var rows []SegmentModel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Segment, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Segment{ID: m.ID, Name: m.Name, RulesJSON: m.RulesJSON, CreatedAt: m.CreatedAt})
	}
	return out, nil
}

func (s *Store) LoadOrdersByEmail(ctx context.Context, email string) ([]map[string]any, error) {
	var rows []map[string]any
	err := s.db.WithContext(ctx).Raw(`
		SELECT o.id, o.status, o.total_amount, o.currency, o.created_at
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE u.email = ?
		ORDER BY o.created_at DESC
		LIMIT 50
	`, email).Scan(&rows).Error
	if err != nil {
		return []map[string]any{}, nil
	}
	return rows, nil
}

func (s *Store) LoadTicketsByEmail(ctx context.Context, email string) ([]map[string]any, error) {
	var rows []map[string]any
	err := s.db.WithContext(ctx).Raw(`
		SELECT t.id, t.subject, t.status, t.priority, t.created_at
		FROM tickets t
		JOIN users u ON u.id = t.requester_id
		WHERE u.email = ?
		ORDER BY t.created_at DESC
		LIMIT 50
	`, email).Scan(&rows).Error
	if err != nil {
		return []map[string]any{}, nil
	}
	return rows, nil
}
