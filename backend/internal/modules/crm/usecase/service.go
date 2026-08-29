package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/crm/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Service struct {
	repo domain.Repository
	bus  *event.Bus
}

func NewService(repo domain.Repository, bus *event.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

func (s *Service) CreateLead(ctx context.Context, name, email, phone, ownerID, source string) (*domain.Lead, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	lead := &domain.Lead{
		ID: uuid.NewString(), Name: name, Email: strings.TrimSpace(email), Phone: strings.TrimSpace(phone),
		Status: "new", OwnerID: ownerID, Source: source, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateLead(ctx, lead); err != nil {
		return nil, err
	}
	if s.bus != nil {
		s.bus.Publish(domain.EventLeadCreated, map[string]any{"lead_id": lead.ID, "owner_id": ownerID, "email": lead.Email})
	}
	return lead, nil
}

func (s *Service) ConvertLeadToCustomer(ctx context.Context, leadID, actorID string) (*domain.Contact, error) {
	lead, err := s.repo.FindLead(ctx, leadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	now := time.Now().UTC()
	contact := &domain.Contact{
		ID: uuid.NewString(), LeadID: lead.ID, Name: lead.Name, Email: lead.Email, Phone: lead.Phone,
		OwnerID: lead.OwnerID, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateContact(ctx, contact); err != nil {
		return nil, err
	}
	lead.Status = "converted"
	lead.UpdatedAt = now
	_ = s.repo.UpdateLead(ctx, lead)
	_ = s.LogActivity(ctx, "lead", lead.ID, "convert", "Lead converted to contact", actorID)
	return contact, nil
}

func (s *Service) MoveDealStage(ctx context.Context, dealID, stage, actorID string) (*domain.Deal, error) {
	deal, err := s.repo.FindDeal(ctx, dealID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	deal.Stage = stage
	deal.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateDeal(ctx, deal); err != nil {
		return nil, err
	}
	_ = s.LogActivity(ctx, "deal", deal.ID, "stage", "Moved to "+stage, actorID)
	if s.bus != nil {
		s.bus.Publish(domain.EventDealStageMoved, map[string]any{"deal_id": deal.ID, "stage": stage})
	}
	return deal, nil
}

func (s *Service) AssignOwner(ctx context.Context, entityType, entityID, ownerID, actorID string) error {
	switch entityType {
	case "lead":
		lead, err := s.repo.FindLead(ctx, entityID)
		if err != nil {
			return ErrNotFound
		}
		lead.OwnerID = ownerID
		lead.UpdatedAt = time.Now().UTC()
		if err := s.repo.UpdateLead(ctx, lead); err != nil {
			return err
		}
	case "deal":
		deal, err := s.repo.FindDeal(ctx, entityID)
		if err != nil {
			return ErrNotFound
		}
		deal.OwnerID = ownerID
		deal.UpdatedAt = time.Now().UTC()
		if err := s.repo.UpdateDeal(ctx, deal); err != nil {
			return err
		}
	case "contact":
		contact, err := s.repo.FindContact(ctx, entityID)
		if err != nil {
			return ErrNotFound
		}
		contact.OwnerID = ownerID
		contact.UpdatedAt = time.Now().UTC()
		if err := s.repo.UpdateContact(ctx, contact); err != nil {
			return err
		}
	default:
		return ErrInvalidInput
	}
	return s.LogActivity(ctx, entityType, entityID, "assign", "Assigned owner "+ownerID, actorID)
}

func (s *Service) LogActivity(ctx context.Context, entityType, entityID, typ, body, actorID string) error {
	a := &domain.Activity{
		ID: uuid.NewString(), EntityType: entityType, EntityID: entityID, Type: typ, Body: body, ActorID: actorID, CreatedAt: time.Now().UTC(),
	}
	return s.repo.CreateActivity(ctx, a)
}

func (s *Service) CreateTaskReminder(ctx context.Context, title, entityType, entityID, ownerID string, dueAt, reminderAt *time.Time) (*domain.Task, error) {
	if title == "" || ownerID == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	task := &domain.Task{
		ID: uuid.NewString(), Title: title, EntityType: entityType, EntityID: entityID, OwnerID: ownerID,
		Status: "open", DueAt: dueAt, ReminderAt: reminderAt, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	if s.bus != nil {
		s.bus.Publish("cronbuilder.schedule_request", map[string]any{
			"source": "crm.task.reminder", "task_id": task.ID, "owner_id": ownerID, "reminder_at": reminderAt,
		})
	}
	return task, nil
}

func (s *Service) ProcessDueReminders(ctx context.Context) error {
	tasks, err := s.repo.ListDueReminders(ctx, time.Now().UTC())
	if err != nil {
		return err
	}
	for i := range tasks {
		t := tasks[i]
		if s.bus != nil {
			s.bus.Publish(domain.EventTaskReminder, map[string]any{
				"task_id": t.ID, "owner_id": t.OwnerID, "title": t.Title,
			})
		}
		t.Status = "reminded"
		t.UpdatedAt = time.Now().UTC()
		t.ReminderAt = nil
		_ = s.repo.UpdateTask(ctx, &t)
	}
	return nil
}

func (s *Service) Customer360(ctx context.Context, contactID string) (*domain.Customer360, error) {
	contact, err := s.repo.FindContact(ctx, contactID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	orders, _ := s.repo.LoadOrdersByEmail(ctx, contact.Email)
	tickets, _ := s.repo.LoadTicketsByEmail(ctx, contact.Email)
	acts, _ := s.repo.ListActivities(ctx, "contact", contact.ID)
	deals, _ := s.repo.ListDealsByContact(ctx, contact.ID)
	tags, _ := s.repo.ListTags(ctx, "contact", contact.ID)
	return &domain.Customer360{
		Contact: contact, Orders: orders, Tickets: tickets, Activities: acts, Deals: deals, Tags: tags,
	}, nil
}

func (s *Service) ListDealsByPipeline(ctx context.Context, pipelineID string) ([]domain.Deal, error) {
	if pipelineID == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.ListDealsByPipeline(ctx, pipelineID)
}

func (s *Service) ListActivities(ctx context.Context, entityType, entityID string) ([]domain.Activity, error) {
	if entityType == "" || entityID == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.ListActivities(ctx, entityType, entityID)
}

func (s *Service) CreateContact(ctx context.Context, name, email, phone, company, ownerID string) (*domain.Contact, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	c := &domain.Contact{
		ID: uuid.NewString(), Name: name, Email: strings.TrimSpace(email), Phone: strings.TrimSpace(phone),
		Company: strings.TrimSpace(company), OwnerID: ownerID, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateContact(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) UpdateContact(ctx context.Context, id, name, email, phone, company string) (*domain.Contact, error) {
	c, err := s.repo.FindContact(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	c.Name = strings.TrimSpace(name)
	c.Email = strings.TrimSpace(email)
	c.Phone = strings.TrimSpace(phone)
	c.Company = strings.TrimSpace(company)
	c.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateContact(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteContact(ctx context.Context, id string) error {
	if _, err := s.repo.FindContact(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.DeleteContact(ctx, id)
}

func (s *Service) ListLeads(ctx context.Context) ([]domain.Lead, error) { return s.repo.ListLeads(ctx) }
func (s *Service) ListContacts(ctx context.Context) ([]domain.Contact, error) {
	return s.repo.ListContacts(ctx)
}

func (s *Service) CreateDeal(ctx context.Context, title, contactID, pipelineID, stage, ownerID string, amount int64) (*domain.Deal, error) {
	now := time.Now().UTC()
	d := &domain.Deal{
		ID: uuid.NewString(), Title: title, ContactID: contactID, PipelineID: pipelineID, Stage: stage,
		OwnerID: ownerID, Amount: amount, CreatedAt: now, UpdatedAt: now,
	}
	return d, s.repo.CreateDeal(ctx, d)
}

func (s *Service) UpsertPipeline(ctx context.Context, name, stagesJSON string) (*domain.Pipeline, error) {
	p := &domain.Pipeline{ID: uuid.NewString(), Name: name, StagesJSON: stagesJSON, CreatedAt: time.Now().UTC()}
	return p, s.repo.UpsertPipeline(ctx, p)
}

func (s *Service) ListPipelines(ctx context.Context) ([]domain.Pipeline, error) {
	return s.repo.ListPipelines(ctx)
}

func (s *Service) UpsertTag(ctx context.Context, name string) (*domain.Tag, error) {
	t := &domain.Tag{ID: uuid.NewString(), Name: name}
	return t, s.repo.UpsertTag(ctx, t)
}

func (s *Service) AttachTag(ctx context.Context, entityType, entityID, tagID string) error {
	return s.repo.AttachTag(ctx, entityType, entityID, tagID)
}

func (s *Service) UpsertSegment(ctx context.Context, name, rulesJSON string) (*domain.Segment, error) {
	seg := &domain.Segment{ID: uuid.NewString(), Name: name, RulesJSON: rulesJSON, CreatedAt: time.Now().UTC()}
	return seg, s.repo.UpsertSegment(ctx, seg)
}

func (s *Service) ListSegments(ctx context.Context) ([]domain.Segment, error) {
	return s.repo.ListSegments(ctx)
}

func (s *Service) BindEvents(bus *event.Bus) {
	if bus == nil {
		return
	}
	bus.Subscribe("cronbuilder.action", func(payload any) {
		_ = s.ProcessDueReminders(context.Background())
	})
}
