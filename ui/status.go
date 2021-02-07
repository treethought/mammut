package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"gitlab.com/tslocum/cview"
)

type StatusModal struct {
	*cview.Modal
	status *mastodon.Status
	app    *App
}

func NewStatusModal(app *App, status *mastodon.Status) *StatusModal {
	m := cview.NewModal()

	s := &StatusModal{
		Modal:  m,
		status: status,
		app:    app,
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

	return &StatusModal{Modal: m, status: status}
}

type StatusFrame struct {
	*cview.Frame
	status *mastodon.Status
	app    *App
}

func NewStatusFrame(app *App) *StatusFrame {

	frame := cview.NewFrame(cview.NewBox())
	frame.SetBackgroundColor(tcell.ColorDefault)
	frame.SetBorders(2, 2, 2, 2, 4, 4)
	frame.SetBorder(true)
	f := &StatusFrame{
		Frame: frame,
		app:   app,
	}

	return f

}

func (f *StatusFrame) SetStatus(status *mastodon.Status) {
	f.Clear()
	f.SetBackgroundColor(tcell.ColorDefault)

	f.status = status

	content := formatContent(status.Content)

	text := cview.NewTextView()
	text.SetText(content)
	text.SetBackgroundColor(tcell.ColorDefault)

	f.Frame = cview.NewFrame(text)

	if f.status == nil {
		return
	}

	ct := status.CreatedAt

	created := fmt.Sprintf("%d-%02d-%02dT%02d",
		ct.Year(), ct.Month(), ct.Day(),
		ct.Hour(), ct.Minute())

	f.AddText(status.Account.DisplayName, true, cview.AlignLeft, tcell.ColorWhite)
	f.AddText(status.Account.Acct, true, cview.AlignCenter, tcell.ColorWhite)
	f.AddText(status.Account.Username, true, cview.AlignRight, tcell.ColorWhite)
	f.AddText(created, true, cview.AlignCenter, tcell.ColorWhite)
	f.AddText(content, false, cview.AlignLeft, tcell.ColorWhite)
	f.AddText("Footer middle", false, cview.AlignCenter, tcell.ColorGreen)
	f.AddText("Footer second middle", false, cview.AlignCenter, tcell.ColorGreen)

}
