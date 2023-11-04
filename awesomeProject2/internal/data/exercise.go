package data

import (
	"awesomeProject2/internal/validator"
	"database/sql"
	"errors"
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

type ExerciseModel struct {
	DB *sql.DB
}

func (e ExerciseModel) Get(id int64) (*Exercise, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
               SELECT id, created_at, title, runtime
               FROM exercise
               WHERE id = $1`

	var exercise Exercise

	err := e.DB.QueryRow(query, id).Scan(
		&exercise.ID,
		&exercise.CreatedAt,
		&exercise.Title,
		&exercise.Runtime,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &exercise, nil
}

func (e ExerciseModel) Update(exercise *Exercise) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
          UPDATE exercise
          SET title = $1, runtime = $3
          WHERE id = $5
          RETURNING title`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		exercise.Title,
		exercise.Runtime,
		exercise.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return e.DB.QueryRow(query, args...).Scan(&exercise.Title)
}

func (e ExerciseModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
          DELETE FROM exercise
          WHERE id = $1`
	result, err := e.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (e ExerciseModel) Insert(exercise *Exercise) error {
	query := `
          INSERT INTO exercise (title, runtime)
          VALUES ($1, $2)
          RETURNING id, created_at`

	args := []interface{}{exercise.Title, exercise.Runtime}

	return e.DB.QueryRow(query, args...).Scan(&exercise.ID, &exercise.CreatedAt)
}
