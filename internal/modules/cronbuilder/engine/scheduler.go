package engine

import (
	"context"
	"sync"

	"crm-starterkit/internal/modules/cronbuilder/domain"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron     *cron.Cron
	executor *Executor
	mu       sync.Mutex
	entries  map[string]cron.EntryID
}

func NewScheduler(executor *Executor) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		executor: executor,
		entries:  make(map[string]cron.EntryID),
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() context.Context {
	return s.cron.Stop()
}

func (s *Scheduler) Register(job *domain.JobDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.entries[job.ID]; ok {
		s.cron.Remove(id)
		delete(s.entries, job.ID)
	}
	if !job.IsActive {
		return nil
	}
	jobCopy := *job
	entryID, err := s.cron.AddFunc(job.CronExpr, func() {
		_, _ = s.executor.Run(context.Background(), &jobCopy)
	})
	if err != nil {
		return err
	}
	s.entries[job.ID] = entryID
	return nil
}

func (s *Scheduler) Unregister(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.entries[jobID]; ok {
		s.cron.Remove(id)
		delete(s.entries, jobID)
	}
}

func (s *Scheduler) LoadAll(jobs []domain.JobDefinition) error {
	for i := range jobs {
		if err := s.Register(&jobs[i]); err != nil {
			return err
		}
	}
	return nil
}
