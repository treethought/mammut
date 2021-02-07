package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"gitlab.com/tslocum/cview"
)

type StatusModal struct {
	*cview.Modal
	toot *mastodon.Status
	app  *App
}

func NewStatusModal(app *App, status *mastodon.Status) *StatusModal {
	m := cview.NewModal()

	s := &StatusModal{
		Modal: m,
		toot:  status,
		app:   app,
	}

	s.SetBorder(true)
	s.SetTitle(status.Account.DisplayName)
	s.AddButtons([]string{"back", "boost", "reply"})
	s.SetTitle(status.Account.DisplayName)
	s.SetText(formatContent(status.Content))
	s.SetBackgroundColor(tcell.ColorDefault)

	s.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "back" {
			s.app.FocusTimeline()
		}
	})

	return &StatusModal{Modal: m, toot: status}
}
