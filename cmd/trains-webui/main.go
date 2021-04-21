package main

import (
	"flag"
	"log"
	"os"

	"git.adyxax.org/adyxax/trains/internal/webui"
	"git.adyxax.org/adyxax/trains/pkg/config"
	"git.adyxax.org/adyxax/trains/pkg/database"
)

func main() {
	path := flag.String("c", os.Getenv("HOME")+"/.config/trains/config.yaml", "configuration file path")
	help := flag.Bool("h", false, "display this help message")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	c, err := config.LoadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.InitDB("sqlite3", "file:test.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	webui.Run(c, db)
}
