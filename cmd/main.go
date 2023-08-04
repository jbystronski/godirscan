package main

import (
	"github.com/jbystronski/godirscan/pkg/terminal"
	u "github.com/jbystronski/godirscan/pkg/utils"
)

func init() {
	u.ParseConfigFile(u.Cfg)
	u.ParseColorSchema(u.Cfg.CurrentSchema)
}

func main() {
	terminal.ClearScreen()

	u.RunMainLoop()
}
