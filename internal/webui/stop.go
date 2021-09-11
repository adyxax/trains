package webui

import (
	"fmt"
	"html/template"
	"net/http"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

var stopTemplate = template.Must(template.New("stop").Funcs(funcMap).ParseFS(templatesFS, "html/base.html", "html/stop.html"))

// The page template variable
type StopPage struct {
	User  *model.User
	Stops []model.Stop
}

// The stop handler of the webui
func stopHandler(e *env, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/stop" {
		user, err := tryAndResumeSession(e, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}
		switch r.Method {
		case http.MethodGet:
			stops, err := e.dbEnv.GetStops()
			if err != nil {
				return newStatusError(http.StatusInternalServerError, fmt.Errorf("Could not get train stops"))
			} else {
				w.Header().Set("Cache-Control", "no-store, no-cache")
			}
			p := StopPage{
				User:  user,
				Stops: stops,
			}
			err = stopTemplate.ExecuteTemplate(w, "stop.html", p)
			if err != nil {
				return newStatusError(http.StatusInternalServerError, err)
			}
			return nil
		default:
			return newStatusError(http.StatusMethodNotAllowed, fmt.Errorf(http.StatusText(http.StatusMethodNotAllowed)))
		}
	} else {
		return newStatusError(http.StatusNotFound, fmt.Errorf("Invalid path in stopHandler"))
	}
}
