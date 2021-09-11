package webui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"regexp"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

var validStopId = regexp.MustCompile(`^stop_area:[a-zA-Z]+:\d+$`)

var specificStopTemplate = template.Must(template.New("specificStop").Funcs(funcMap).ParseFS(templatesFS, "html/base.html", "html/specificStop.html"))

// The page template variable
type SpecificStopPage struct {
	User       *model.User
	Stop       string
	Departures []model.Departure
}

// The stop handler of the webui
func specificStopHandler(e *env, w http.ResponseWriter, r *http.Request) error {
	if path.Dir(r.URL.Path) == "/stop" {
		user, err := tryAndResumeSession(e, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}
		switch r.Method {
		case http.MethodGet:
			id := path.Base(r.URL.Path)
			if id == "" {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("No id in query string")) // TODO should we redirect to root page to chose a stop id?
			}
			if ok := validStopId.MatchString(id); !ok {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("Invalid stop id"))
			}
			stop, err := e.dbEnv.GetStop(id)
			if err != nil {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("Stop id not found in database")) // TODO do better
			}
			if departures, err := e.navitia.GetDepartures(stop.Id); err != nil {
				log.Printf("%s; data returned: %+v\n", err, departures)
				return newStatusError(http.StatusInternalServerError, fmt.Errorf("Could not get departures"))
			} else {
				w.Header().Set("Cache-Control", "no-store, no-cache")
				p := SpecificStopPage{
					User:       user,
					Stop:       stop.Name,
					Departures: departures,
				}
				err = specificStopTemplate.ExecuteTemplate(w, "specificStop.html", p)
				if err != nil {
					return newStatusError(http.StatusInternalServerError, err)
				}
				return nil
			}
		default:
			return newStatusError(http.StatusMethodNotAllowed, fmt.Errorf(http.StatusText(http.StatusMethodNotAllowed)))
		}
	} else {
		return newStatusError(http.StatusNotFound, fmt.Errorf("Invalid path in specificStopHandler"))
	}
}
