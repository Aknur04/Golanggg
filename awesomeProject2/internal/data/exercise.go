package data

import (
	"awesomeProject2/internal/validator"
	"time"
)

type Exercise struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive
	Title     string    `json:"title"`
	Runtime   Runtime   `json:"runtime,omitempty"`
}

func Validateexercise(v *validator.Validator, exercise *Exercise) {
	v.Check(exercise.Title != "", "title", "must be provided")
	v.Check(len(exercise.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(exercise.Runtime != 0, "runtime", "must be provided")
	v.Check(exercise.Runtime > 0, "runtime", "must be a positive integer")
}
