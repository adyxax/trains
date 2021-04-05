package main

import (
	"log"

	"git.adyxax.org/adyxax/trains/internal/webui"
	"git.adyxax.org/adyxax/trains/pkg/config"
)

func main() {
	c, err := config.LoadFile("/home/julien/.config/adyxax-trains/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	webui.Run(c)
}
