package main

import (
	"fmt"
	"net/http"
)

func (app *application) createexerciseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new exercise")
}
func (app *application) showexerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show notes %d\n", id)
}
