package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Task struct {
	ID          uuid.UUID
	UserID      string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
