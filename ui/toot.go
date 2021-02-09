package ui

import (
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/kyokomi/emoji/v2"
	"github.com/mattn/go-mastodon"
	"gitlab.com/tslocum/cview"
)

type Toot struct {
	*cview.ListItem
	status *mastodon.Status
	app    *App
}

func formatContent(html string) string {
	// opts := &md.Options{}
	converter := md.NewConverter("", true, nil)

	mdContent, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}
	return mdContent

	result := markdown.Render(mdContent, 80, 6)
	return string(result)

}

func (t *Toot) IsFavorite() bool {
	favorited, ok := t.status.Favourited.(bool)
	if !ok {
		return false
	}
	return favorited

}

func (t *Toot) header() string {
	header := t.status.Account.DisplayName
	if t.IsFavorite() {
		header += emoji.Sprint(" :heart:")
	} else {
		header += emoji.Sprint(" :white_heart:")
	}
	return header

}

func NewToot(app *App, status *mastodon.Status) *Toot {

	t := &Toot{
		ListItem: cview.NewListItem(status.Account.DisplayName),
		status:   status,
		app:      app,
	}

	content := formatContent(t.status.Content)
	main := t.header()

	if status.Reblog != nil {
		main = emoji.Sprintf("%s  || :repeat_button:@%s", main, status.Reblog.Account.DisplayName)
	}

	t.SetMainText(main)
	t.SetSecondaryText(content)
	t.SetReference(t)

	return t
}
