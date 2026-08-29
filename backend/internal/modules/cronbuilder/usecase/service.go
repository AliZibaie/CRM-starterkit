package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/modules/cronbuilder/domain"
	"crm-starterkit/internal/modules/cronbuilder/engine"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var (
	ErrJobNotFound   = errors.New("job not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInvalidCron   = errors.New("invalid cron expression")
	ErrInvalidAction = errors.New("invalid action")
	ErrInvalidTable  = errors.New("invalid table")
)

type Service struct {
	jobs     domain.JobRepository
	logs     domain.LogRepository
	intro    domain.SchemaIntrospector
	actions  domain.ActionRegistry
	executor *engine.Executor
	scheduler *engine.Scheduler
}

func NewService(
	jobs domain.JobRepository,
	logs domain.LogRepository,
	intro domain.SchemaIntrospector,
	actions domain.ActionRegistry,
	executor *engine.Executor,
	scheduler *engine.Scheduler,
) *Service {
	return &Service{jobs: jobs, logs: logs, intro: intro, actions: actions, executor: executor, scheduler: scheduler}
}

func (s *Service) ListTables() []domain.TableMeta {
	return s.intro.ListTables()
}

func (s *Service) ListColumns(table string) ([]domain.ColumnMeta, error) {
	return s.intro.ListColumns(table)
}

func (s *Service) ListActions() []domain.ActionInfo {
	return s.actions.List()
}

func (s *Service) BuildCondition(table string, node domain.ConditionNode) (*engine.BuiltQuery, error) {
	return engine.BuildCondition(table, node, s.intro)
}

func (s *Service) CreateCronJob(ctx context.Context, in domain.CreateJobInput) (*domain.JobDefinition, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || strings.TrimSpace(in.CronExpr) == "" {
		return nil, ErrInvalidInput
	}
	if !s.intro.IsTableAllowed(in.TableName) {
		return nil, ErrInvalidTable
	}
	if !s.actions.Has(in.Action.Type) {
		return nil, ErrInvalidAction
	}
	if _, err := cron.ParseStandard(in.CronExpr); err != nil {
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		if _, err2 := parser.Parse(in.CronExpr); err2 != nil {
			return nil, ErrInvalidCron
		}
	}
	if _, err := engine.BuildCondition(in.TableName, in.Condition, s.intro); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	job := &domain.JobDefinition{
		ID: uuid.NewString(), Name: name, TableName: in.TableName,
		Condition: in.Condition, Action: in.Action, CronExpr: in.CronExpr,
		IsActive: in.IsActive, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.jobs.Create(ctx, job); err != nil {
		return nil, err
	}
	if s.scheduler != nil {
		_ = s.scheduler.Register(job)
	}
	return job, nil
}

func (s *Service) GetJob(ctx context.Context, id string) (*domain.JobDefinition, error) {
	job, err := s.jobs.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}
	return job, nil
}

func (s *Service) UpdateJob(ctx context.Context, id string, in domain.CreateJobInput) (*domain.JobDefinition, error) {
	job, err := s.jobs.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}
	name := strings.TrimSpace(in.Name)
	if name == "" || strings.TrimSpace(in.CronExpr) == "" {
		return nil, ErrInvalidInput
	}
	if !s.intro.IsTableAllowed(in.TableName) {
		return nil, ErrInvalidTable
	}
	if !s.actions.Has(in.Action.Type) {
		return nil, ErrInvalidAction
	}
	if _, err := cron.ParseStandard(in.CronExpr); err != nil {
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		if _, err2 := parser.Parse(in.CronExpr); err2 != nil {
			return nil, ErrInvalidCron
		}
	}
	if _, err := engine.BuildCondition(in.TableName, in.Condition, s.intro); err != nil {
		return nil, err
	}
	job.Name = name
	job.TableName = in.TableName
	job.Condition = in.Condition
	job.Action = in.Action
	job.CronExpr = in.CronExpr
	job.IsActive = in.IsActive
	job.UpdatedAt = time.Now().UTC()
	if err := s.jobs.Update(ctx, job); err != nil {
		return nil, err
	}
	if s.scheduler != nil {
		_ = s.scheduler.Register(job)
	}
	return job, nil
}

func (s *Service) DeleteJob(ctx context.Context, id string) error {
	if _, err := s.jobs.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrJobNotFound
		}
		return err
	}
	return s.jobs.Delete(ctx, id)
}

func (s *Service) ListJobs(ctx context.Context) ([]domain.JobDefinition, error) {
	return s.jobs.List(ctx)
}

func (s *Service) RunJobNow(ctx context.Context, id string) (*domain.JobExecutionLog, error) {
	job, err := s.jobs.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}
	return s.executor.Run(ctx, job)
}

func (s *Service) JobExecutionLogs(ctx context.Context, jobID string, limit int) ([]domain.JobExecutionLog, error) {
	if _, err := s.jobs.FindByID(ctx, jobID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}
	return s.logs.ListByJob(ctx, jobID, limit)
}
