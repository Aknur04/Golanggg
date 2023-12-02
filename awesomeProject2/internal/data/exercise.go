package data

import (
	"awesomeProject2/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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

func (e ExerciseModel) Insert(exercise *Exercise) error {
	query := `
             INSERT INTO exercise (title, runtime)
             VALUES ($1, $2)
             RETURNING id, created_at`
	args := []interface{}{exercise.Title, exercise.Runtime}

	return e.DB.QueryRow(query, args...).Scan(&exercise.ID, &exercise.CreatedAt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return e.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.ID, &exercise.CreatedAt)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := e.DB.QueryRowContext(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.CreatedAt,
		&exercise.Title,
		&exercise.Runtime,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &exercise, nil
}

func (e ExerciseModel) Update(exercise *Exercise) error {
	query := `
UPDATE exercise
SET title = $1, runtime = runtime + 1
WHERE id = $2 AND runtime = $3
RETURNING runtime`
	args := []interface{}{
		exercise.Title,
		exercise.Runtime,
		exercise.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.Runtime)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (e ExerciseModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM exercise
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := e.DB.ExecContext(ctx, query, id)
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
func (e ExerciseModel) GetAll(title string, filters Filters) ([]*Exercise, Metadata, error) {
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, title, runtime
FROM exercise
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (title @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{title, pq.Array(title), filters.limit(), filters.offset()}
	rows, err := e.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	exercises := []*Exercise{}
	for rows.Next() {
		var exercise Exercise
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&exercise.ID,
			&exercise.CreatedAt,
			&exercise.Title,
			&exercise.Runtime,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		exercises = append(exercises, &exercise)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return exercises, metadata, nil
}
