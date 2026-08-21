package domain

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUserInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	IsActive  bool
	RoleIDs   []string
}

type UpdateUserInput struct {
	Email     *string
	Password  *string
	FirstName *string
	LastName  *string
	IsActive  *bool
}
