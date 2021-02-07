package ui

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	ma "github.com/treethought/mammut/mastodon"
	"gitlab.com/tslocum/cview"
)

type App struct {
	client     ma.Client
	ui         *cview.Application
	root       *cview.Flex
	timeline   *Timeline
	info       *cview.TextView
	statusView *StatusFrame
}

func New() *App {
	client := ma.NewClient()
	return &App{
		client: client,
	}

}

func (app *App) FocusTimeline() {
	// Set the grid as the application root and focus the timeline
	app.ui.SetRoot(app.root, true)
	app.ui.SetFocus(app.timeline)

	app.timeline.SetTitle("...")

	go app.ui.QueueUpdateDraw(func() {
		toots := app.client.GetTimeline()
		app.timeline.fillToots(toots)
		app.timeline.SetTitle("Timeline")
	})
}

func (app *App) SetStatus(toot *Toot) {
	if app.statusView != nil {
		app.statusView.SetStatus(toot)

	}

	// app.Notify(fmt.Sprintf("Viewing status by: %s", status.Account.DisplayName))

	// go app.ui.QueueUpdateDraw(func() {
	// app.statusView.SetStatus(status)
	// frame := NewStatusFrame(app, status)
	// app.statusView = nil
	// app.statusView = frame
	// })
}

func (app *App) Notify(msg string) {
	if app.info == nil {
		return
	}
	app.info.Clear()
	app.info.SetText(msg)
	go app.ui.QueueUpdateDraw(func() {
		time.Sleep(2 * time.Second)
		app.info.Clear()

	})
}

func (app *App) Start() {
	// Initialize application
	app.ui = cview.NewApplication()

	toots := app.client.GetTimeline()
	if len(toots) == 0 {
		log.Fatal("Failed to get toots")
	}
	app.timeline = NewTimeline(app, toots)

	leftpanel := cview.NewBox()
	leftpanel.SetBackgroundColor(tcell.ColorDefault)
	//

	app.statusView = NewStatusFrame(app)

	app.info = cview.NewTextView()

	mid := cview.NewFlex()
	mid.SetBackgroundColor(tcell.ColorDefault)
	mid.SetDirection(cview.FlexRow)
	mid.AddItem(app.timeline, 0, 3, true)
	mid.AddItem(app.statusView, 0, 2, false)
	mid.AddItem(app.info, 2, 1, false)

	flex := cview.NewFlex()
	flex.SetBackgroundTransparent(false)
	flex.SetBackgroundColor(tcell.ColorDefault)
	// flex.AddItem(leftpanel, 0, 2, false)
	flex.AddItem(mid, 0, 3, false)
	// flex.AddItem(leftpanel, 0, 1, false)

	// flex.AddItem(app.info, 0, 1, false)

	// Create Grid containing the application's widgets
	// grid := cview.NewGrid()
	// grid.SetColumns(-1, -4, -1)
	// grid.SetRows(1, -1, 1)
	// grid.AddItem(app.timeline, 1, 1, 1, 1, 0, 0, false) // Left - 3 rows

	app.root = flex

	app.FocusTimeline()

	err := app.ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
