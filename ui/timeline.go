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
	t.SetBackgroundColor(tcell.ColorDefault)
	t.SetMainTextColor(tcell.ColorLightCyan)
	t.SetSecondaryTextColor(tcell.ColorNavajoWhite)
	t.SetSelectedBackgroundColor(tcell.ColorIndianRed)
	t.SetInputCapture(t.HandleInput)

	t.fillToots(toots)
	return t

}

func (t *Timeline) fillToots(toots []*mastodon.Status) {
	t.Clear()
	t.Toots = toots
	for _, toot := range t.Toots {
		tc := NewToot(t.app, toot)
		t.AddItem(tc.ListItem)
	}

}

func (t *Timeline) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	ref := t.GetCurrentItem().GetReference()
	toot, ok := ref.(*Toot)
	if !ok {
		return nil
	}
	status := toot.status

	key := event.Key()
	switch key {
	case tcell.KeyEnter:

		m := NewStatusModal(t.app, toot.status)
		t.app.ui.SetRoot(m, true)

		return nil

	case tcell.KeyRune:
		switch event.Rune() {
		case 't': // Toot
			m := NewComposeModal(t.app)
			t.app.ui.SetRoot(m, true)

		case 'r': // Refresh
			t.app.FocusTimeline()

		case 'l': // Like
			if toot.IsFavorite() {
				t.app.client.Unlike(status)
			} else {
				t.app.client.Like(status)
			}
			t.app.FocusTimeline()

			return nil
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
