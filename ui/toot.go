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
	converter := md.NewConverter("", true, nil)

	mdContent, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}

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

func NewToot(app *App, status *mastodon.Status) *Toot {
	t := &Toot{
		ListItem: cview.NewListItem(status.Account.DisplayName),
		status:   status,
		app:      app,
	}

	content := formatContent(t.status.Content)

	main := t.status.Account.DisplayName
	if t.IsFavorite() {
		main += emoji.Sprint(":heart:")
	} else {
		main += emoji.Sprint(":white_heart:")
	}

	t.SetMainText(main)
	t.SetSecondaryText(content)
	t.SetReference(t)

	return t
}
