package webui

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"git.adyxax.org/adyxax/trains/pkg/config"
	"git.adyxax.org/adyxax/trains/pkg/navitia_api_client"
)

// the api client object
var client *navitia_api_client.Client

// the webui configuration
var conf *config.Config

//go:embed html/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

// The page template variable
type Page struct {
	Departures []Departure
	Title      string
}
type Departure struct {
	DisplayName string
	Arrival     string
	Odd         bool
}

// The root handler of the webui
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		var p Page
		if d, err := client.GetDepartures(conf.TrainStop); err != nil {
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
		renderTemplate(w, "index", p)
	} else {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

var templates = template.Must(template.ParseFS(templatesFS, "html/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Run(c *config.Config) {
	conf = c
	client = navitia_api_client.NewClient(c.Token)
	http.Handle("/static/", http.FileServer(http.FS(staticFS)))
	http.HandleFunc("/", rootHandler)

	listenStr := c.Address + ":" + c.Port
	log.Printf("Starting webui on %s", listenStr)
	log.Fatal(http.ListenAndServe(listenStr, nil))
}
