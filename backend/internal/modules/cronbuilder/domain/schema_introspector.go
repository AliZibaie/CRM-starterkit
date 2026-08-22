package domain

type SchemaIntrospector interface {
	ListTables() []TableMeta
	ListColumns(table string) ([]ColumnMeta, error)
	IsTableAllowed(table string) bool
	IsColumnAllowed(table, column string) bool
}
