package main

import (
	"github.com/jbystronski/godirscan/pkg/app"

	"github.com/jbystronski/godirscan/pkg/app/config"
)

func init() {
	config.Running().Parse()

	schema := config.Running().ColorSchemas[config.Running().CurrentSchema]

	config.SetTheme(schema, config.CurrentTheme)
}

func main() {
	app.New()
}
