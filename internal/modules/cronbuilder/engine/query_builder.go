package engine

import (
	"fmt"
	"strings"

	"crm-starterkit/internal/modules/cronbuilder/domain"
)

type BuiltQuery struct {
	WhereSQL string
	Args     []any
}

func BuildCondition(table string, node domain.ConditionNode, intro domain.SchemaIntrospector) (*BuiltQuery, error) {
	if !intro.IsTableAllowed(table) {
		return nil, fmt.Errorf("table not allowed")
	}
	sql, args, err := buildNode(table, node, intro)
	if err != nil {
		return nil, err
	}
	if sql == "" {
		return &BuiltQuery{WhereSQL: "1=1", Args: nil}, nil
	}
	return &BuiltQuery{WhereSQL: sql, Args: args}, nil
}

func buildNode(table string, node domain.ConditionNode, intro domain.SchemaIntrospector) (string, []any, error) {
	logic := strings.ToUpper(strings.TrimSpace(node.Logic))
	if logic == "" {
		logic = "AND"
	}
	if logic != "AND" && logic != "OR" {
		return "", nil, fmt.Errorf("invalid logic")
	}
	if len(node.Children) > 0 {
		parts := make([]string, 0, len(node.Children))
		var args []any
		for _, child := range node.Children {
			s, a, err := buildNode(table, child, intro)
			if err != nil {
				return "", nil, err
			}
			if s == "" {
				continue
			}
			parts = append(parts, "("+s+")")
			args = append(args, a...)
		}
		if len(parts) == 0 {
			return "", nil, nil
		}
		return strings.Join(parts, " "+logic+" "), args, nil
	}
	if node.Field == "" || node.Operator == "" {
		return "", nil, nil
	}
	if !intro.IsColumnAllowed(table, node.Field) {
		return "", nil, fmt.Errorf("column not allowed: %s", node.Field)
	}
	col := QuoteIdent(node.Field)
	switch node.Operator {
	case domain.OpEq:
		return col + " = ?", []any{node.Value}, nil
	case domain.OpNeq:
		return col + " <> ?", []any{node.Value}, nil
	case domain.OpGt:
		return col + " > ?", []any{node.Value}, nil
	case domain.OpGte:
		return col + " >= ?", []any{node.Value}, nil
	case domain.OpLt:
		return col + " < ?", []any{node.Value}, nil
	case domain.OpLte:
		return col + " <= ?", []any{node.Value}, nil
	case domain.OpContains:
		return col + " LIKE ?", []any{fmt.Sprintf("%%%v%%", node.Value)}, nil
	case domain.OpIsNull:
		return col + " IS NULL", nil, nil
	case domain.OpNotNull:
		return col + " IS NOT NULL", nil, nil
	default:
		return "", nil, fmt.Errorf("unsupported operator")
	}
}
