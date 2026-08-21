package domain

import "context"

type JobRepository interface {
	Create(ctx context.Context, job *JobDefinition) error
	Update(ctx context.Context, job *JobDefinition) error
	FindByID(ctx context.Context, id string) (*JobDefinition, error)
	List(ctx context.Context) ([]JobDefinition, error)
	ListActive(ctx context.Context) ([]JobDefinition, error)
	Delete(ctx context.Context, id string) error
}

type LogRepository interface {
	Create(ctx context.Context, log *JobExecutionLog) error
	ListByJob(ctx context.Context, jobID string, limit int) ([]JobExecutionLog, error)
}
