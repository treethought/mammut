package ui

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"gitlab.com/tslocum/cview"
)

type TimelineType int

const (
	TimelineHome TimelineType = iota
	TimelineLocal
	TimelineFederated
	TimelineProfile
	TimelineTag
)

var TimelineTypes = []TimelineType{
	TimelineHome,
	TimelineLocal,
	TimelineFederated,
	TimelineProfile,
	TimelineTag,
}

func (t TimelineType) String() string {
	return [...]string{"home", "local", "federated", "profile", "tags"}[t]
}

type Timeline struct {
	*cview.List
	Toots []*mastodon.Status
	app   *App
	ttype TimelineType
}

func NewTimeline(app *App, toots []*mastodon.Status, ttype TimelineType) *Timeline {
	t := &Timeline{
		List:  cview.NewList(),
		Toots: toots,
		app:   app,
	}
	t.SetTitle("Timeline")
	t.SetBorder(true)
	t.SetPadding(1, 1, 1, 1)
	t.SetBackgroundColor(tcell.ColorDefault)
	t.SetMainTextColor(tcell.ColorLightCyan)
	t.SetSecondaryTextColor(tcell.ColorNavajoWhite)
	t.SetSelectedBackgroundColor(tcell.ColorIndianRed)
	t.SetInputCapture(t.HandleInput)

	t.SetChangedFunc(func(index int, item *cview.ListItem) {
		toot := t.GetCurrentToot()
		app.SetStatus(toot)
	})

	t.fillToots(toots)
	return t
}

func (t *Timeline) SetTimeline(ttype TimelineType) {
	t.ttype = ttype
}

func (t *Timeline) GetCurrentToot() *Toot {
	ref := t.GetCurrentItem().GetReference()
	toot, ok := ref.(*Toot)
	if !ok {
		return nil
	}
	return toot

}

func (t *Timeline) Refresh() {
	toots := t.app.client.GetTimeline(t.ttype.String())
	t.fillToots(toots)
	title := fmt.Sprintf(" Timeline - %s ", strings.Title(t.ttype.String()))
	t.SetTitle(title)
	t.SetTitleColor(tcell.ColorLightCyan)
}

func (t *Timeline) fillToots(toots []*mastodon.Status) {
	t.Clear()
	t.Toots = toots
	for _, toot := range t.Toots {
		tc := NewToot(t.app, toot)
		t.AddItem(tc.ListItem)
	}

}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func (t *Timeline) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	toot := t.GetCurrentToot()
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
