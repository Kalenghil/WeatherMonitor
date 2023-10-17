package main

import (
	"WeatherMonitor/gui"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Setenv("FYNE_THEME", "dark")
	gui.NewPlot()
}
