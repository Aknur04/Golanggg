package data

import (
	"time"
)

type Exercise struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive
	Title     string    `json:"title"`
	Runtime   Runtime   `json:"runtime,omitempty"`
}
