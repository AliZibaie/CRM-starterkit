package domain

import "time"

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Permission struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Module      string    `json:"module"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
