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
	ID          string
	Name        string
	TableName   string
	Condition   ConditionNode
	Action      JobAction
	CronExpr    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type JobExecutionLog struct {
	ID           string
	JobID        string
	StartedAt    time.Time
	FinishedAt   time.Time
	Success      bool
	AffectedRows int64
	ErrorMessage string
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
