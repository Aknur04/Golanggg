package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthCheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/exercise", app.requirePermission("exercise:read", app.listexerciseHandler))
	router.HandlerFunc(http.MethodPost, "/v1/exercise", app.requirePermission("exercise:write", app.createexerciseHandler))
	router.HandlerFunc(http.MethodGet, "/v1/exercise/:id", app.requirePermission("exercise:read", app.showexerciseHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/exercise/:id", app.requirePermission("exercise:write", app.updateexerciseHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/exercise/:id", app.requirePermission("exercise:write", app.deleteexerciseHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
