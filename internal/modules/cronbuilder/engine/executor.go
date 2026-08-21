package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/cronbuilder/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Executor struct {
	db     *gorm.DB
	intro  domain.SchemaIntrospector
	actions domain.ActionRegistry
	logs   domain.LogRepository
	bus    *event.Bus
}

func NewExecutor(db *gorm.DB, intro domain.SchemaIntrospector, actions domain.ActionRegistry, logs domain.LogRepository, bus *event.Bus) *Executor {
	return &Executor{db: db, intro: intro, actions: actions, logs: logs, bus: bus}
}

func (e *Executor) Run(ctx context.Context, job *domain.JobDefinition) (*domain.JobExecutionLog, error) {
	started := time.Now().UTC()
	log := &domain.JobExecutionLog{
		ID:        uuid.NewString(),
		JobID:     job.ID,
		StartedAt: started,
	}
	affected, err := e.execute(ctx, job)
	log.FinishedAt = time.Now().UTC()
	log.AffectedRows = affected
	if err != nil {
		log.Success = false
		log.ErrorMessage = err.Error()
	} else {
		log.Success = true
	}
	_ = e.logs.Create(ctx, log)
	return log, err
}

func (e *Executor) execute(ctx context.Context, job *domain.JobDefinition) (int64, error) {
	if !e.actions.Has(job.Action.Type) {
		return 0, fmt.Errorf("unknown action")
	}
	built, err := BuildCondition(job.TableName, job.Condition, e.intro)
	if err != nil {
		return 0, err
	}
	table := QuoteIdent(job.TableName)
	switch job.Action.Type {
	case "update":
		sets, ok := job.Action.Params["set"].(map[string]any)
		if !ok {
			raw, _ := json.Marshal(job.Action.Params["set"])
			_ = json.Unmarshal(raw, &sets)
		}
		if len(sets) == 0 {
			return 0, fmt.Errorf("update requires set params")
		}
		parts := make([]string, 0, len(sets))
		args := make([]any, 0, len(sets)+len(built.Args))
		for col, val := range sets {
			if !e.intro.IsColumnAllowed(job.TableName, col) {
				return 0, fmt.Errorf("column not allowed: %s", col)
			}
			parts = append(parts, QuoteIdent(col)+" = ?")
			args = append(args, val)
		}
		args = append(args, built.Args...)
		sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(parts, ", "), built.WhereSQL)
		res := e.db.WithContext(ctx).Exec(sql, args...)
		return res.RowsAffected, res.Error
	case "notify", "sms", "email", "webhook":
		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, built.WhereSQL)
		var count int64
		if err := e.db.WithContext(ctx).Raw(sql, built.Args...).Scan(&count).Error; err != nil {
			return 0, err
		}
		if e.bus != nil {
			e.bus.PublishAsync("cronbuilder.action", map[string]any{
				"job_id":   job.ID,
				"action":   job.Action.Type,
				"table":    job.TableName,
				"affected": count,
				"params":   job.Action.Params,
			})
		}
		return count, nil
	default:
		return 0, fmt.Errorf("unsupported action")
	}
}
