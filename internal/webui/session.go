package webui

import (
	"net/http"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

func tryAndResumeSession(e *env, r *http.Request) (*model.User, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}
	user, err := e.dbEnv.ResumeSession(cookie.Value)
	if err != nil {
		return nil, err
	}
	return user, nil
}
