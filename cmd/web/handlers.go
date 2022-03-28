package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sevskii111/microservices-goods/pkg/forms"
	"github.com/sevskii111/microservices-goods/pkg/models"
)

func (app *application) list(w http.ResponseWriter, r *http.Request) {
	goods, err := app.goods.List()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, goods)
}

func (app *application) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.notFound(w)
		return
	}

	g, err := app.goods.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.json(w, g)
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "description", "price", "inStock")
	form.MaxLength("name", 255)
	form.MinValue("price", 1)
	form.PermittedValues("inStock", "false", "true")

	if !form.Valid() {
		app.formResult(w, form)
		return
	}

	price, err := strconv.Atoi(form.Get("price"))
	if err != nil {
		form.Errors.Add("price", "Must be numeric")
		app.formResult(w, form)
		return
	}
	inStock, err := strconv.ParseBool(form.Get("inStock"))
	if err != nil {
		form.Errors.Add("inStock", "Should be either 0 or 1")
		app.formResult(w, form)
		return
	}

	_, err = app.goods.Insert(form.Get("name"), form.Get("description"), price, inStock)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.formResult(w, form)
}

func (app *application) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.notFound(w)
		return
	}

	form := forms.New(r.PostForm)
	form.MaxLength("name", 255)
	form.MinValue("price", 1)

	if !form.Valid() {
		app.formResult(w, form)
		return
	}

	update := &models.GoodUpdate{}

	if form.Has("name") {
		name := form.Get("name")
		update.Name = &name
	}
	if form.Has("description") {
		description := form.Get("description")
		update.Description = &description
	}
	if form.Has("price") {
		price, err := strconv.Atoi(form.Get("price"))
		if err != nil {
			form.Errors.Add("price", "Must be numeric")
			app.formResult(w, form)
			return
		}
		update.Price = &price
	}
	if form.Has("inStock") {
		inStock, err := strconv.ParseBool(form.Get("inStock"))
		if err != nil {
			form.Errors.Add("inStock", "Should be either 0 or 1")
			app.formResult(w, form)
			return
		}
		update.InStock = &inStock
	}
	if update.IsEmpty() {
		form.Errors.Add("generic", "Update is empty")
		app.formResult(w, form)
		return
	}

	err = app.goods.Update(id, update)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.formResult(w, form)
}
