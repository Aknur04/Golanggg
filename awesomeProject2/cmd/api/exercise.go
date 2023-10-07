package main

import (
	"awesomeProject2/internal/data" // New import
	"fmt"
	"net/http"
	"time" // New import
)

func (app *application) createexerciseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new exercise")
}
func (app *application) showexerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}
	exercise := data.Exercise{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "abs body",
		Runtime:   102,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"exercise": exercise}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
