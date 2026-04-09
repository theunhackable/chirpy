package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`

	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
