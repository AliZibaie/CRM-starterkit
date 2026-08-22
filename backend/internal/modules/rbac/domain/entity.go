package domain

import "time"

type Role struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ID          string
	Key         string
	Name        string
	Description string
	Module      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Policy struct {
	RoleID       string
	PermissionID string
}

type CreateRoleInput struct {
	Name        string
	Description string
}

type CreatePermissionInput struct {
	Key         string
	Name        string
	Description string
	Module      string
}

type EditPermissionInput struct {
	Name        *string
	Description *string
	Module      *string
}
