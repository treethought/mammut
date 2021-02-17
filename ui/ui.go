package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	ma "github.com/treethought/mammut/mastodon"
	"gitlab.com/tslocum/cbind"
	"gitlab.com/tslocum/cview"
)

type App struct {
	client       ma.Client
	ui           *cview.Application
	root         *cview.Flex
	timeline     *Timeline
	info         *cview.TextView
	statusView   *StatusFrame
	menu         *Menu
	thread       *TootReplies
	panels       *cview.Panels
	compose      *ComposeModal
	focusManager *cview.FocusManager
	reply        *ReplyForm
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
		app.timeline.Refresh()
	})
}

func (app *App) SetStatus(toot *Toot) {
	if app.statusView != nil {
		app.statusView.SetStatus(toot)
	}
	if app.reply != nil {
		app.reply.SetStatus(toot.status)
	}

	// app.Notify(fmt.Sprintf("Viewing status by: %s", toot.status.Account.DisplayName))

}

func (app *App) ViewTimeline() {
	app.panels.SetCurrentPanel("timeline")
	app.panels.SendToFront("tiemeline")
	app.ui.SetFocus(app.timeline)
}

func (app *App) ViewThread(toot *Toot) {
	app.thread.SetStatus(toot.status)
	app.panels.SetCurrentPanel("thread")
	app.panels.SendToFront("thread")
	app.ui.SetFocus(app.thread)
}

func (app *App) ViewCompose() {
	app.panels.SetCurrentPanel("compose")
	app.panels.SendToFront("compose")
	app.ui.SetFocus(app.compose)
}

func (app *App) Notify(msg string, a ...interface{}) {
	if app.info == nil {
		return
	}
	app.info.Clear()
	text := fmt.Sprintf(msg, a...)
	app.info.SetText(text)
	go app.ui.QueueUpdateDraw(func() {
		time.Sleep(2 * time.Second)
		app.info.Clear()

	})
}

func (app *App) initViews() {
	toots := app.client.GetTimeline("local")
	if len(toots) == 0 {
		log.Fatal("Failed to get toots")
	}
	app.timeline = NewTimeline(app, toots, TimelineLocal)

	app.menu = NewMenu(app)

	acctInfo := NewInfoView(app)

	app.statusView = NewStatusFrame(app)

	app.thread = NewTootReplies(app)

	app.info = cview.NewTextView()
	app.info.SetBorder(true)
	app.info.SetBackgroundColor(tcell.ColorDefault)

	app.compose = NewComposeModal(app)
	app.reply = NewReplyForm(app)

	panels := cview.NewPanels()
	panels.AddPanel("timeline", app.timeline, true, true)
	panels.AddPanel("thread", app.thread, true, true)
	panels.AddPanel("compose", app.compose, true, true)
	panels.SetCurrentPanel("timeline")
	app.panels = panels

	mid := cview.NewFlex()
	mid.SetBackgroundColor(tcell.ColorDefault)
	mid.SetDirection(cview.FlexRow)
	mid.AddItem(app.panels, 0, 4, true)
	mid.AddItem(app.statusView, 0, 4, false)
	mid.AddItem(app.reply, 0, 1, false)

	flex := cview.NewFlex()
	flex.SetBackgroundTransparent(false)
	flex.SetBackgroundColor(tcell.ColorDefault)

	left := cview.NewFlex()
	left.SetDirection(cview.FlexRow)
	left.AddItem(app.menu, 0, 7, false)
	left.AddItem(app.info, 0, 1, false)
	left.AddItem(acctInfo, 0, 1, false)

	flex.AddItem(left, 0, 1, false)
	flex.AddItem(mid, 0, 4, false)
	app.root = flex

}

func (app *App) handleCompose(ev *tcell.EventKey) *tcell.EventKey {
	if app.reply.HasFocus() {
		return ev
	}
	current, _ := app.panels.GetFrontPanel()
	if current == "compose" {
		return ev
	}
	app.ViewCompose()
	return nil
}

func (app *App) handleToggle(ev *tcell.EventKey) *tcell.EventKey {
	if app.reply.HasFocus() {
		return ev
	}

	current, _ := app.panels.GetFrontPanel()
	if current == "compose" {
		return ev
	}
	app.focusManager.FocusNext()
	return nil

}

func (app *App) handleComment(ev *tcell.EventKey) *tcell.EventKey {
	if app.reply.HasFocus() {
		return ev
	}
	app.ui.SetFocus(app.reply)
	return nil
}

func (app *App) initBindings() {
	c := cbind.NewConfiguration()

	c.SetRune(tcell.ModNone, 't', app.handleCompose)
	c.SetRune(tcell.ModNone, 'c', app.handleComment)
	c.SetRune(tcell.ModNone, 'i', app.handleComment)
	c.SetKey(tcell.ModNone, tcell.KeyTAB, app.handleToggle)
	app.ui.SetInputCapture(c.Capture)

}

func (app *App) initInputHandler() {
	app.focusManager = cview.NewFocusManager(app.ui.SetFocus)
	app.focusManager.SetWrapAround(true)
	app.focusManager.Add(app.menu, app.panels)

}

func (app *App) Start() {
	// Initialize application
	app.ui = cview.NewApplication()

	app.initViews()
	app.initInputHandler()
	app.initBindings()

	app.FocusTimeline()

	err := app.ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func wrap(f func()) func(ev *tcell.EventKey) *tcell.EventKey {
	return func(ev *tcell.EventKey) *tcell.EventKey {
		f()
		return nil
	}
}
