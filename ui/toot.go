package ui

import (
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/mattn/go-mastodon"
	"gitlab.com/tslocum/cview"
)

type Toot struct {
	*cview.ListItem
	status *mastodon.Status
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

func NewToot(status *mastodon.Status) *Toot {
	t := &Toot{
		ListItem: cview.NewListItem(status.Account.DisplayName),
		status:   status,
	}
	content := formatContent(t.status.Content)

	t.SetMainText(t.status.Account.DisplayName)
	t.SetSecondaryText(content)
	t.SetReference(t.status)
	return t
}

