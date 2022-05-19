package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app, err := NewTodos()

	if err != nil {
		log.Fatalln("init todo:", err)
	}
	log.Println("inited todos")

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "wails2-todo",
		Width:     1024,
		Height:    768,
		Assets:    assets,
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err)
	}
}
