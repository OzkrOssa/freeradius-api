package domain

import "time"

type User struct {
	ID        uint64
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
