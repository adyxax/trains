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

	if i, err := dbEnv.CountStops(); err == nil && i == 0 {
		log.Printf("No trains stops data found, updating...")
		if trainStops, err := e.navitia.GetStops(); err == nil {
			log.Printf("Updated trains stops data from navitia api, got %d results", len(trainStops))
			if err = dbEnv.ReplaceAndImportStops(trainStops); err != nil {
				if dberr, ok := err.(*database.QueryError); ok {
					log.Printf("%+v", dberr.Unwrap())
				}
			}
		} else {
			log.Printf("Failed to get trains stops data from navitia api : %+v", err)
		}
	}

	listenStr := c.Address + ":" + c.Port
	log.Printf("Starting webui on %s", listenStr)
	log.Fatal(http.ListenAndServe(listenStr, nil))
}
