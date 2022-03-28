package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sevskii111/microservices-goods/pkg/forms"
)

type JSONRes struct {
	Errors  forms.FormErrors
	Success bool
}

func (app *application) json(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(res)
	if err != nil {
		app.serverError(w, err)
	}
	w.Write(jsonResp)
}

func (app *application) formResult(w http.ResponseWriter, form *forms.Form) {
	res := JSONRes{
		Errors:  form.Errors,
		Success: form.Valid(),
	}
	app.json(w, res)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) forbiden(w http.ResponseWriter) {
	app.clientError(w, http.StatusForbidden)
}
