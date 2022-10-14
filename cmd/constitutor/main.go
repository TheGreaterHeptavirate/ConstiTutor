package main

import (
	"log"

	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Panic(err)
	}

	a.Run()
}