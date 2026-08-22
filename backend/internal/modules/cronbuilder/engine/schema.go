package engine

import (
	"fmt"
	"strings"

	"crm-starterkit/internal/modules/cronbuilder/domain"
)

type WhitelistIntrospector struct {
	tables map[string][]domain.ColumnMeta
}

func NewWhitelistIntrospector() *WhitelistIntrospector {
	return &WhitelistIntrospector{
		tables: map[string][]domain.ColumnMeta{
			"users": {
				{Name: "id", Type: "string"}, {Name: "email", Type: "string"},
				{Name: "first_name", Type: "string"}, {Name: "last_name", Type: "string"},
				{Name: "is_active", Type: "bool"}, {Name: "created_at", Type: "timestamp"},
			},
			"roles": {
				{Name: "id", Type: "string"}, {Name: "name", Type: "string"},
				{Name: "description", Type: "string"}, {Name: "created_at", Type: "timestamp"},
			},
			"announcements": {
				{Name: "id", Type: "string"}, {Name: "title", Type: "string"},
				{Name: "status", Type: "string"}, {Name: "created_at", Type: "timestamp"},
			},
			"banners": {
				{Name: "id", Type: "string"}, {Name: "title", Type: "string"},
				{Name: "position", Type: "string"}, {Name: "is_active", Type: "bool"},
				{Name: "starts_at", Type: "timestamp"}, {Name: "ends_at", Type: "timestamp"},
			},
			"tickets": {
				{Name: "id", Type: "string"}, {Name: "subject", Type: "string"},
				{Name: "status", Type: "string"}, {Name: "priority", Type: "string"},
				{Name: "assigned_to", Type: "string"}, {Name: "created_at", Type: "timestamp"},
			},
		},
	}
}

func (w *WhitelistIntrospector) ListTables() []domain.TableMeta {
	out := make([]domain.TableMeta, 0, len(w.tables))
	for name := range w.tables {
		out = append(out, domain.TableMeta{Name: name})
	}
	return out
}

func (w *WhitelistIntrospector) ListColumns(table string) ([]domain.ColumnMeta, error) {
	cols, ok := w.tables[table]
	if !ok {
		return nil, fmt.Errorf("table not allowed: %s", table)
	}
	return cols, nil
}

func (w *WhitelistIntrospector) IsTableAllowed(table string) bool {
	_, ok := w.tables[table]
	return ok
}

func (w *WhitelistIntrospector) IsColumnAllowed(table, column string) bool {
	cols, ok := w.tables[table]
	if !ok {
		return false
	}
	for _, c := range cols {
		if c.Name == column {
			return true
		}
	}
	return false
}

type DefaultActionRegistry struct {
	actions []domain.ActionInfo
}

func NewDefaultActionRegistry() *DefaultActionRegistry {
	return &DefaultActionRegistry{
		actions: []domain.ActionInfo{
			{Type: "update", Description: "Update matching rows with given fields"},
			{Type: "notify", Description: "Publish notification event for matching rows"},
			{Type: "sms", Description: "Queue SMS for matching rows"},
			{Type: "email", Description: "Queue email for matching rows"},
			{Type: "webhook", Description: "Call webhook with matching row payload"},
		},
	}
}

func (r *DefaultActionRegistry) List() []domain.ActionInfo {
	return append([]domain.ActionInfo(nil), r.actions...)
}

func (r *DefaultActionRegistry) Has(actionType string) bool {
	for _, a := range r.actions {
		if a.Type == actionType {
			return true
		}
	}
	return false
}

func QuoteIdent(name string) string {
	cleaned := strings.ReplaceAll(name, `"`, "")
	cleaned = strings.ReplaceAll(cleaned, "`", "")
	cleaned = strings.ReplaceAll(cleaned, ";", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	return cleaned
}
