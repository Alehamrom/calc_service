package main

import (
	"github.com/Alehamrom/calc_service/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
