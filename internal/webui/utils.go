package webui

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"git.adyxax.org/adyxax/trains/pkg/config"
	"git.adyxax.org/adyxax/trains/pkg/database"
	"git.adyxax.org/adyxax/trains/pkg/navitia_api_client"
)

//go:embed html/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

// Template functions
var funcMap = template.FuncMap{
	"odd": func(i int) bool {
		return i%2 == 1
	},
}

// the environment that will be passed to our handlers
type env struct {
	conf    *config.Config
	dbEnv   *database.DBEnv
	navitia navitia_api_client.Client
}

type handlerError interface {
	error
	Status() int
}

type statusError struct {
	code int
	err  error
}

func (e *statusError) Error() string           { return e.err.Error() }
func (e *statusError) Status() int             { return e.code }
func newStatusError(code int, err error) error { return &statusError{code: code, err: err} }

type handler struct {
	e *env
	h func(e *env, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our handler type to satisfy http.Handler
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	err := h.h(h.e, w, r)
	if err != nil {
		switch e := err.(type) {
		case handlerError:
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default to serving a HTTP 500
			log.Printf("%s : handler returned an unexpected error : %+v", path, e)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
