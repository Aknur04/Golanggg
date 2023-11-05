package data

import (
	"database/sql"
	"errors"
)

// Example usage of ErrRecordNotFound

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Exercises ExerciseModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Exercises: ExerciseModel{DB: db},
	}
}
