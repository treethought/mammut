package ui

import (
	"log"

	ma "github.com/treethought/masto/mastodon"
	"gitlab.com/tslocum/cview"
)

type App struct {
	client   ma.Client
	ui       *cview.Application
	grid     *cview.Grid
	timeline *Timeline
}

func New() *App {
	client := ma.NewClient()
	return &App{
		client: client,
	}

}

func (app *App) FocusTimeline() {
	// Set the grid as the application root and focus the timeline
	app.ui.SetRoot(app.grid, true)
	app.ui.SetFocus(app.timeline)

}

func (app *App) Start() {
	// Initialize application
	app.ui = cview.NewApplication()

	toots := app.client.GetTimeline()
	app.timeline = NewTimeline(app, toots)

	// Create Grid containing the application's widgets
	grid := cview.NewGrid()
	grid.SetColumns(-1, -3, -1)
	grid.SetRows(1, -1, 1)
	grid.AddItem(app.timeline, 1, 1, 1, 1, 0, 0, false) // Left - 3 rows

	app.grid = grid

	app.FocusTimeline()

	err := app.ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
