package main

import (
	"github.com/leveldorado/experiment/api/bootstrap"
	"github.com/leveldorado/experiment/app"
)

func main() {
	app.MustRun(&bootstrap.App{})
}
