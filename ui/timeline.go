package ui

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"gitlab.com/tslocum/cbind"
	"gitlab.com/tslocum/cview"
)

type TimelineType int

const (
	TimelineHome TimelineType = iota
	TimelineLocal
	TimelineFederated
	TimelineProfile
	TimelineLikes
	TimelineTag
	TimelineMedia

	TimelineTootContext
)

var TimelineTypes = []TimelineType{
	TimelineHome,
	TimelineLocal,
	TimelineFederated,
	TimelineProfile,
	TimelineLikes,
	TimelineTag,
	TimelineMedia,
}

func (t TimelineType) String() string {
	return [...]string{"home", "local", "federated", "profile", "likes", "tags", "media"}[t]
}

type Timeline struct {
	*cview.List
	Toots        []*mastodon.Status
	app          *App
	ttype        TimelineType
	inputHandler *cbind.Configuration
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

	t.inputHandler = cbind.NewConfiguration()
	t.initBindings()

	t.SetChangedFunc(func(index int, item *cview.ListItem) {
		toot := t.GetCurrentToot()
		app.SetStatus(toot)
	})

	t.fillToots(toots)
	return t
}

func (t *Timeline) handleSelect(ev *tcell.EventKey) *tcell.EventKey {
	toot := t.GetCurrentToot()
	t.app.ViewThread(toot)

	return nil
}

func (t *Timeline) handleDelete(ev *tcell.EventKey) *tcell.EventKey {

	toot := t.GetCurrentToot()
	status := toot.status

	if t.app.client.IsOwnStatus(status) {
		t.app.Notify("Deleting toot!")
		t.app.client.Delete(status)
		t.app.FocusTimeline()
		return nil
	}
	return ev
}

func (t *Timeline) handleRefresh(ev *tcell.EventKey) *tcell.EventKey {
	t.app.FocusTimeline()
	return nil

}

func (t *Timeline) handleFollow(ev *tcell.EventKey) *tcell.EventKey {
	toot := t.GetCurrentToot()
	status := toot.status

	t.app.Notify("Following %s", status.Account.Acct)
	t.app.client.Follow(status.Account.ID)
	t.app.FocusTimeline()
	return nil

}

func (t *Timeline) handleUnfollow(ev *tcell.EventKey) *tcell.EventKey {
	toot := t.GetCurrentToot()
	status := toot.status

	t.app.Notify("Unfollowing %s", status.Account.Acct)
	t.app.client.Unfollow(status.Account.ID)
	t.app.FocusTimeline()
	return nil

}

func (t *Timeline) handleLike(ev *tcell.EventKey) *tcell.EventKey {
	toot := t.GetCurrentToot()
	status := toot.status

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
}
func (t *Timeline) handleOpen(ev *tcell.EventKey) *tcell.EventKey {
	t.app.Notify("Opening in browser")
	toot := t.GetCurrentToot()
	status := toot.status
	openbrowser(status.URL)
	return nil

}

func (t *Timeline) initBindings() {

	t.inputHandler.SetKey(tcell.ModNone, tcell.KeyEnter, t.handleSelect)
	t.inputHandler.SetRune(tcell.ModNone, 'd', t.handleDelete)
	t.inputHandler.SetRune(tcell.ModNone, 'r', t.handleRefresh)
	t.inputHandler.SetRune(tcell.ModNone, 'l', t.handleLike)
	t.inputHandler.SetRune(tcell.ModNone, 'o', t.handleOpen)
	t.inputHandler.SetRune(tcell.ModNone, 'f', t.handleFollow)
	t.inputHandler.SetRune(tcell.ModNone, 'u', t.handleUnfollow)

	t.SetInputCapture(t.inputHandler.Capture)

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
func (t *Timeline) SetCurrentToot(toot *Toot) {
	for i, item := range t.GetItems() {
		ref := item.GetReference()
		tootc, ok := ref.(*Toot)
		if !ok {
			continue
		}
		if tootc.status.ID == toot.status.ID {
			t.SetCurrentItem(i)
		}
	}
}

func (t *Timeline) Refresh() {
	selected := t.GetCurrentToot()
	toots := t.app.client.GetTimeline(t.ttype.String())
	t.fillToots(toots)
	title := fmt.Sprintf(" Timeline - %s ", strings.Title(t.ttype.String()))
	t.SetTitle(title)
	t.SetTitleColor(tcell.ColorLightCyan)
	if selected != nil {
		t.SetCurrentToot(selected)
	}
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
