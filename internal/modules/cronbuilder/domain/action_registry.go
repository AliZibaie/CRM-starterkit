package domain

type ActionRegistry interface {
	List() []ActionInfo
	Has(actionType string) bool
}
