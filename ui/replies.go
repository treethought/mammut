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

	// apply thread specific bindings on top of timeline bindings
	r.initThreadBindings()

	return r

}

func (t *TootReplies) handleBack(ev *tcell.EventKey) *tcell.EventKey {
	t.app.ViewTimeline()
	return nil

}

// initRepliesBindings adds toot context specific bindings on top of timeline bindings
func (t *TootReplies) initThreadBindings() {

	t.inputHandler.SetRune(tcell.ModNone, 'h', t.handleBack)
	t.inputHandler.SetKey(tcell.ModNone, tcell.KeyEscape, t.handleBack)

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
