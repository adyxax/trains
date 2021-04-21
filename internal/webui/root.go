package webui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

var rootTemplate = template.Must(template.ParseFS(templatesFS, "html/base.html", "html/root.html"))

// The page template variable
type Page struct {
	User       *model.User
	Departures []Departure
	Title      string
}
type Departure struct {
	DisplayName string
	Arrival     string
	Odd         bool
}

// The root handler of the webui
func rootHandler(e *env, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/" {
		var p Page
		user, err := tryAndResumeSession(e, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}
		p.User = user
		if d, err := e.navitia.GetDepartures(e.conf.TrainStop); err != nil {
			log.Printf("%+v\n%s\n", d, err)
		} else {
			for i := 0; i < len(d.Departures); i++ {
				t, err := time.Parse("20060102T150405", d.Departures[i].StopDateTime.ArrivalDateTime)
				if err != nil {
					panic(err)
				}
				p.Departures = append(p.Departures, Departure{d.Departures[i].DisplayInformations.Direction, t.Format("Mon, 02 Jan 2006 15:04:05"), i%2 == 1})
			}
			w.Header().Set("Cache-Control", "no-store, no-cache")
		}
		p.Title = "Horaires des prochains trains à Crépieux la Pape"
		err = rootTemplate.ExecuteTemplate(w, "root.html", p)
		if err != nil {
			return newStatusError(http.StatusInternalServerError, err)
		}
		return nil
	} else {
		return newStatusError(http.StatusNotFound, fmt.Errorf("Invalid path in rootHandler"))
	}
}
