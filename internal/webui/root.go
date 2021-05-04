package webui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

var rootTemplate = template.Must(template.ParseFS(templatesFS, "html/base.html", "html/root.html"))

// The page template variable
type Page struct {
	User       *model.User
	Departures []model.Departure
	Title      string
}

// The root handler of the webui
func rootHandler(e *env, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/" {
		user, err := tryAndResumeSession(e, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}
		var departures []model.Departure
		if departures, err := e.navitia.GetDepartures(e.conf.TrainStop); err != nil {
			log.Printf("%s; data returned: %+v\n", err, departures)
			return newStatusError(http.StatusInternalServerError, fmt.Errorf("Could not get departures"))
		} else {
			w.Header().Set("Cache-Control", "no-store, no-cache")
		}
		p := Page{
			User:       user,
			Departures: departures,
			Title:      "Horaires des prochains trains à Crépieux la Pape",
		}
		err = rootTemplate.ExecuteTemplate(w, "root.html", p)
		if err != nil {
			return newStatusError(http.StatusInternalServerError, err)
		}
		return nil
	} else {
		return newStatusError(http.StatusNotFound, fmt.Errorf("Invalid path in rootHandler"))
	}
}
