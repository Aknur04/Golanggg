package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/yoga", app.yogaHandler)

	router.HandlerFunc(http.MethodGet, "/v1/exercise", app.listexerciseHandler)
	router.HandlerFunc(http.MethodPost, "/v1/exercise", app.createexerciseHandler)
	router.HandlerFunc(http.MethodGet, "/v1/exercise/:id", app.showexerciseHandler)
	// Require a PATCH request, rather than PUT.
	router.HandlerFunc(http.MethodPatch, "/v1/exercise/:id", app.updateexerciseHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/exercise/:id", app.deleteexerciseHandler)
	return app.recoverPanic(router)
}
