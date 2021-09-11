package webui

import (
	"fmt"
	"html/template"
	"net/http"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

var rootTemplate = template.Must(template.New("root").Funcs(funcMap).ParseFS(templatesFS, "html/base.html", "html/root.html"))

// The page template variable
type RootPage struct {
	User *model.User
}

// The root handler of the webui
func rootHandler(e *env, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/" {
		user, err := tryAndResumeSession(e, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}
		w.Header().Set("Cache-Control", "no-store, no-cache")
		p := RootPage{
			User: user,
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
