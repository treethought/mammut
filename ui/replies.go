package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
)

type TootReplies struct {
	*Timeline
	app    *App
	status *mastodon.Status
}

func NewTootReplies(app *App) *TootReplies {
	t := NewTimeline(app, []*mastodon.Status{}, TimelineTootContext)

	r := &TootReplies{
		Timeline: t,
		app:      app,
	}

	r.SetInputCapture(r.HandleInput)

	return r

}

func (r *TootReplies) SetStatus(status *mastodon.Status) {
	r.status = status
	thread := r.app.client.GetStatusContext(r.status)

	r.SetTitle(fmt.Sprintf(" Replies to [cyan]@%s ", status.Account.DisplayName))

	toots := []*mastodon.Status{}

	for _, s := range thread.Ancestors {
		toots = append(toots, s)
	}
	toots = append(toots, status)

	for _, s := range thread.Descendants {
		toots = append(toots, s)
	}
	r.fillToots(toots)
}

func (t *TootReplies) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	toot := t.GetCurrentToot()
	status := toot.status

	key := event.Key()

	switch key {
	case tcell.KeyEnter:
		t.app.ViewThread(toot)
		return nil

	case tcell.KeyEscape:
		t.app.ViewTimeline()
		return nil

	case tcell.KeyRune:
		switch event.Rune() {
		case 'h': // Back
			t.app.ViewTimeline()
			return nil
			m := NewComposeModal(t.app)
			t.app.ui.SetRoot(m, true)

		case 't': // Toot
			m := NewComposeModal(t.app)
			t.app.ui.SetRoot(m, true)

		case 'r': // Refresh
			t.app.FocusTimeline()

		case 'd': // Delete
			if t.app.client.IsOwnStatus(status) {
				t.app.Notify("Deleting toot!")
				t.app.client.Delete(status)
				t.app.FocusTimeline()
				return nil
			}
			return event

		case 'l': // Like
			if toot.IsFavorite() {
				t.app.Notify("Unliking toot!")
				t.app.client.Unlike(status)
				t.app.FocusTimeline()
				return nil
			}
			t.app.Notify("Liking toot!")
			t.app.client.Like(status)
			t.app.FocusTimeline()
			return nil

		case 'o': // Open
			openbrowser(status.URL)
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
