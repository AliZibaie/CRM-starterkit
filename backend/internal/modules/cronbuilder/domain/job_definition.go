package domain

import "time"

type ConditionOp string

const (
	OpEq       ConditionOp = "eq"
	OpNeq      ConditionOp = "neq"
	OpGt       ConditionOp = "gt"
	OpGte      ConditionOp = "gte"
	OpLt       ConditionOp = "lt"
	OpLte      ConditionOp = "lte"
	OpContains ConditionOp = "contains"
	OpIsNull   ConditionOp = "is_null"
	OpNotNull  ConditionOp = "is_not_null"
)

type ConditionNode struct {
	Logic    string           `json:"logic,omitempty"`
	Field    string           `json:"field,omitempty"`
	Operator ConditionOp      `json:"operator,omitempty"`
	Value    any              `json:"value,omitempty"`
	Children []ConditionNode  `json:"children,omitempty"`
}

type JobAction struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params,omitempty"`
}

type JobDefinition struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	TableName   string        `json:"table_name"`
	Condition   ConditionNode `json:"condition"`
	Action      JobAction     `json:"action"`
	CronExpr    string        `json:"cron_expr"`
	IsActive    bool          `json:"is_active"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type JobExecutionLog struct {
	ID           string    `json:"id"`
	JobID        string    `json:"job_id"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
	Success      bool      `json:"success"`
	AffectedRows int64     `json:"affected_rows"`
	ErrorMessage string    `json:"error_message"`
}

type ColumnMeta struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TableMeta struct {
	Name    string       `json:"name"`
	Columns []ColumnMeta `json:"columns,omitempty"`
}

type ActionInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type CreateJobInput struct {
	Name      string
	TableName string
	Condition ConditionNode
	Action    JobAction
	CronExpr  string
	IsActive  bool
}
