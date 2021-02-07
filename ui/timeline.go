package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"gitlab.com/tslocum/cview"
)

type Timeline struct {
	*cview.List
	Toots []*mastodon.Status
	app   *App
}

func NewTimeline(app *App, toots []*mastodon.Status) *Timeline {
	t := &Timeline{
		List:  cview.NewList(),
		Toots: toots,
		app:   app,
	}
	t.SetTitle("Timeline")
	t.SetBorder(true)
	t.SetMainTextColor(tcell.ColorLightCyan)
	t.SetSecondaryTextColor(tcell.ColorNavajoWhite)
	t.SetSelectedBackgroundColor(tcell.ColorIndianRed)
	t.SetInputCapture(t.HandleInput)

	for _, toot := range t.Toots {
		tc := NewToot(toot)
		t.AddItem(tc.ListItem)
		tc.View()

	}
	return t

}

func (t *Timeline) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	key := event.Key()
	switch key {
	case tcell.KeyEnter:

		ref := t.GetCurrentItem().GetReference()
		toot, ok := ref.(*mastodon.Status)
		if !ok {
			return nil
		}

		m := NewStatusModal(t.app, toot)
		t.app.ui.SetRoot(m, true)

		return nil

	case tcell.KeyRune:
		switch event.Rune() {
		case 'g': // Home.
			t.SetCurrentItem(0)
		case 'G': // End.
			t.SetCurrentItem(-1)
		case 'j': // Down.
			cur := t.GetCurrentItemIndex()
			t.SetCurrentItem(cur + 1)
		case 'k': // Up.
			cur := t.GetCurrentItemIndex()
			t.SetCurrentItem(cur - 1)
		}

		return nil
	}

	return event
}
