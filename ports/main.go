package main

import (
	"github.com/leveldorado/experiment/app"
	"github.com/leveldorado/experiment/ports/bootstrap"
)

func main() {
	app.MustRun(&bootstrap.App{})
}
