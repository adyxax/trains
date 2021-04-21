package webui

import (
	"log"
	"net/http"

	"git.adyxax.org/adyxax/trains/pkg/config"
	"git.adyxax.org/adyxax/trains/pkg/database"
	"git.adyxax.org/adyxax/trains/pkg/navitia_api_client"
)

func Run(c *config.Config, dbEnv *database.DBEnv) {
	e := env{
		conf:    c,
		dbEnv:   dbEnv,
		navitia: navitia_api_client.NewClient(c.Token),
	}
	http.Handle("/", handler{&e, rootHandler})
	http.Handle("/login", handler{&e, loginHandler})
	http.Handle("/static/", http.FileServer(http.FS(staticFS)))

	listenStr := c.Address + ":" + c.Port
	log.Printf("Starting webui on %s", listenStr)
	log.Fatal(http.ListenAndServe(listenStr, nil))
}
