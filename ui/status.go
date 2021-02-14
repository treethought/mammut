package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/gdamore/tcell/v2"
	"github.com/kyokomi/emoji/v2"
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
	toot *Toot
	app  *App
}

func NewStatusFrame(app *App) *StatusFrame {

	box := cview.NewBox()
	box.SetBackgroundColor(tcell.ColorDefault)

	frame := cview.NewFrame(box)
	// frame.SetBackgroundTransparent(true)
	frame.SetBackgroundColor(tcell.ColorDefault)
	frame.SetBorders(2, 2, 3, 3, 4, 4)
	frame.SetBorder(true)
	f := &StatusFrame{
		Frame: frame,
		app:   app,
	}

	return f

}

func (f *StatusFrame) SetStatus(toot *Toot) {
	f.Clear()

	f.toot = toot
	status := toot.status

	content := formatContent(status.Content)

	text := cview.NewTextView()
	text.SetBackgroundColor(tcell.ColorDefault)
	text.SetDynamicColors(true)

	_, _, w, h := f.GetInnerRect()

	for _, m := range toot.status.MediaAttachments {
		if m.Type == "image" {

			w = w - 5
			h = h - len(strings.Split(content, "\n")) - 5

			img, err := buildImage(m.URL, w, h)
			if err != nil {
				continue
			}

			ans := img.Render()
			trans := cview.TranslateANSI(ans)
			content = fmt.Sprintf("%s\n%s", content, trans)
		}
	}

	text.SetText(content)

	f.Frame = cview.NewFrame(text)
	f.SetBackgroundColor(tcell.ColorDefault)
	f.SetBorder(true)

	if f.toot == nil {
		return
	}

	ct := status.CreatedAt

	created := fmt.Sprintf("%02d:%02d %d-%02d-%02d",
		ct.Hour(), ct.Minute(),
		ct.Year(), ct.Month(), ct.Day())

	replies := emoji.Sprintf(":speech_balloon: %d", status.RepliesCount)
	boosts := emoji.Sprintf(":repeat_button: %d", status.ReblogsCount)

	likes := ""
	if toot.IsFavorite() {
		likes += emoji.Sprintf(":heart: %d", status.FavouritesCount)
	} else {
		likes += emoji.Sprintf(":white_heart: %d", status.FavouritesCount)
	}

	info := strings.Join([]string{replies, boosts, likes}, " | ")

	f.AddText(status.Account.DisplayName, true, cview.AlignLeft, tcell.ColorWhite)
	f.AddText(status.Account.Acct, true, cview.AlignCenter, tcell.ColorWhite)
	f.AddText(status.Account.Username, true, cview.AlignRight, tcell.ColorWhite)
	f.AddText(created, true, cview.AlignCenter, tcell.ColorWhite)

	f.AddText(info, false, cview.AlignCenter, tcell.ColorWhite)
	if status.Reblog != nil {
		boosted := fmt.Sprintf("Boosted from %s", status.Reblog.Account.DisplayName)
		f.AddText(boosted, false, cview.AlignRight, tcell.ColorLightCyan)
	}

}

func buildImage(url string, x, y int) (*ansimage.ANSImage, error) {
	pix, err := ansimage.NewScaledFromURL(url, y, x, color.Transparent, ansimage.ScaleModeResize, ansimage.NoDithering)
	if err != nil {
		return nil, err
	}
	return pix, nil

}
