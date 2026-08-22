package repository

import (
	"context"
	"encoding/json"
	"time"

	"crm-starterkit/internal/modules/cronbuilder/domain"

	"gorm.io/gorm"
)

type JobModel struct {
	ID            string    `gorm:"primaryKey;size:36"`
	Name          string    `gorm:"size:150;not null"`
	TargetTable   string    `gorm:"column:table_name;size:100;not null"`
	ConditionJSON string    `gorm:"type:text;not null"`
	ActionJSON    string    `gorm:"type:text;not null"`
	CronExpr      string    `gorm:"size:100;not null"`
	IsActive      bool      `gorm:"not null;default:true"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

func (JobModel) TableName() string { return "cron_jobs" }

type LogModel struct {
	ID           string    `gorm:"primaryKey;size:36"`
	JobID        string    `gorm:"size:36;index;not null"`
	StartedAt    time.Time `gorm:"not null"`
	FinishedAt   time.Time `gorm:"not null"`
	Success      bool      `gorm:"not null"`
	AffectedRows int64     `gorm:"not null;default:0"`
	ErrorMessage string    `gorm:"type:text"`
}

func (LogModel) TableName() string { return "cron_job_logs" }

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&JobModel{}, &LogModel{})
}

func (s *Store) Create(ctx context.Context, job *domain.JobDefinition) error {
	m, err := toJobModel(job)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Create(m).Error
}

func (s *Store) Update(ctx context.Context, job *domain.JobDefinition) error {
	m, err := toJobModel(job)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Save(m).Error
}

func (s *Store) FindByID(ctx context.Context, id string) (*domain.JobDefinition, error) {
	var m JobModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toJobDomain(&m)
}

func (s *Store) List(ctx context.Context) ([]domain.JobDefinition, error) {
	var rows []JobModel
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.JobDefinition, 0, len(rows))
	for i := range rows {
		j, err := toJobDomain(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *j)
	}
	return out, nil
}

func (s *Store) ListActive(ctx context.Context) ([]domain.JobDefinition, error) {
	var rows []JobModel
	if err := s.db.WithContext(ctx).Where("is_active = ?", true).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.JobDefinition, 0, len(rows))
	for i := range rows {
		j, err := toJobDomain(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *j)
	}
	return out, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(&JobModel{}, "id = ?", id).Error
}

func (s *Store) CreateLog(ctx context.Context, log *domain.JobExecutionLog) error {
	return s.db.WithContext(ctx).Create(&LogModel{
		ID: log.ID, JobID: log.JobID, StartedAt: log.StartedAt, FinishedAt: log.FinishedAt,
		Success: log.Success, AffectedRows: log.AffectedRows, ErrorMessage: log.ErrorMessage,
	}).Error
}

func (s *Store) ListByJob(ctx context.Context, jobID string, limit int) ([]domain.JobExecutionLog, error) {
	if limit <= 0 {
		limit = 50
	}
	var rows []LogModel
	if err := s.db.WithContext(ctx).Where("job_id = ?", jobID).Order("started_at desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.JobExecutionLog, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.JobExecutionLog{
			ID: r.ID, JobID: r.JobID, StartedAt: r.StartedAt, FinishedAt: r.FinishedAt,
			Success: r.Success, AffectedRows: r.AffectedRows, ErrorMessage: r.ErrorMessage,
		})
	}
	return out, nil
}

type LogStore struct{ *Store }

func (s *Store) AsLogRepo() domain.LogRepository { return &LogStore{s} }
func (l *LogStore) Create(ctx context.Context, log *domain.JobExecutionLog) error {
	return l.Store.CreateLog(ctx, log)
}

func toJobModel(job *domain.JobDefinition) (*JobModel, error) {
	cond, err := json.Marshal(job.Condition)
	if err != nil {
		return nil, err
	}
	act, err := json.Marshal(job.Action)
	if err != nil {
		return nil, err
	}
	return &JobModel{
		ID: job.ID, Name: job.Name, TargetTable: job.TableName,
		ConditionJSON: string(cond), ActionJSON: string(act),
		CronExpr: job.CronExpr, IsActive: job.IsActive,
		CreatedAt: job.CreatedAt, UpdatedAt: job.UpdatedAt,
	}, nil
}

func toJobDomain(m *JobModel) (*domain.JobDefinition, error) {
	var cond domain.ConditionNode
	var act domain.JobAction
	if err := json.Unmarshal([]byte(m.ConditionJSON), &cond); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(m.ActionJSON), &act); err != nil {
		return nil, err
	}
	return &domain.JobDefinition{
		ID: m.ID, Name: m.Name, TableName: m.TargetTable,
		Condition: cond, Action: act, CronExpr: m.CronExpr, IsActive: m.IsActive,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}
