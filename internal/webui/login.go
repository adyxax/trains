package webui

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"git.adyxax.org/adyxax/trains/pkg/database"
	"git.adyxax.org/adyxax/trains/pkg/model"
)

const sessionCookieName = "session-trains-webui"

var validUsername = regexp.MustCompile(`^[a-zA-Z]\w*$`)
var validPassword = regexp.MustCompile(`^.+$`)

var loginTemplate = template.Must(template.ParseFS(templatesFS, "html/base.html", "html/login.html"))

type LoginPage struct{} // no variables to pass for now, but a previous error message would be good

// The login handler of the webui
func loginHandler(e *env, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/login" {
		_, err := tryAndResumeSession(e, r)
		if err == nil {
			// already logged in
			http.Redirect(w, r, "/", http.StatusFound) // TODO fail some other way, at least check if username parameter matches the logged in user
			return nil
		}
		switch r.Method {
		case http.MethodPost:
			r.ParseForm()
			// username
			username, ok := r.Form["username"]
			if !ok {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("No username field in POST"))
			}
			if len(username) != 1 {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("Invalid multiple username fields in POST"))
			}
			if ok := validUsername.MatchString(username[0]); !ok {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("Invalid username field in POST"))
			}
			// password
			password, ok := r.Form["password"]
			if !ok {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("No password field in POST"))
			}
			if len(password) != 1 {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("Invalid multiple password fields in POST"))
			}
			if ok := validPassword.MatchString(password[0]); !ok {
				return newStatusError(http.StatusBadRequest, fmt.Errorf("Invalid password field in POST"))
			}
			// try to login
			user, err := e.dbEnv.Login(&model.UserLogin{Username: username[0], Password: password[0]})
			if err != nil {
				switch e := err.(type) {
				case database.PasswordError:
					// TODO : handle in page
					return e
				case database.QueryError:
					return e
				default:
					return e
				}
			}
			token, err := e.dbEnv.CreateSession(user)
			if err != nil {
				return newStatusError(http.StatusInternalServerError, err)
			}
			cookie := http.Cookie{Name: sessionCookieName, Value: *token, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, MaxAge: 3600000}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		case http.MethodGet:
			p := LoginPage{}
			err := loginTemplate.ExecuteTemplate(w, "login.html", p)
			if err != nil {
				return newStatusError(http.StatusInternalServerError, err)
			}
			return nil
		default:
			return newStatusError(http.StatusMethodNotAllowed, fmt.Errorf(http.StatusText(http.StatusMethodNotAllowed)))
		}
	} else {
		return newStatusError(http.StatusNotFound, fmt.Errorf("Invalid path in loginHandler"))
	}
}
