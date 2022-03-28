package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, app.requireLogin)
	formMiddleware := alice.New(app.requireAdmin, app.parseForm)

	mux := pat.New()
	mux.Get("/goods", standardMiddleware.ThenFunc(app.list))
	mux.Post("/goods", formMiddleware.ThenFunc(app.create))
	mux.Get("/goods/:id", standardMiddleware.ThenFunc(app.get))
	mux.Patch("/goods/:id", formMiddleware.ThenFunc(app.update))

	return standardMiddleware.Then(mux)
}
