package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"gitlab.com/tslocum/cview"
)

type ComposeModal struct {
	*cview.Modal
	app *App
}

func NewComposeModal(app *App) *ComposeModal {

	c := &ComposeModal{
		Modal: cview.NewModal(),
		app:   app,
	}

	c.SetBorder(true)
	c.SetTitle("what's up?")
	c.AddButtons([]string{"back", "toot!"})
	c.SetBackgroundColor(tcell.ColorDefault)

	form := c.GetForm()
	form.AddInputField("content", "", 80, nil, nil)

	c.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "back" {
			c.app.ViewTimeline()
			c.app.FocusTimeline()
		}
		if buttonLabel == "toot!" {
			item := form.GetFormItemByLabel("content")
			input, ok := item.(*cview.InputField)
			if !ok {
				panic("couldnt get input")
			}
			content := input.GetText()
			c.app.client.Toot(content)

			c.app.ViewTimeline()
			c.app.FocusTimeline()
		}
	})

	return c
}

type ReplyForm struct {
	*cview.Form
	app    *App
	status *mastodon.Status
}

func (r *ReplyForm) sendReply() {
	r.app.Notify("Replying to @%s", r.status.Account.Acct)

	item := r.GetFormItemByLabel("Reply")
	input, ok := item.(*cview.InputField)
	if !ok {
		panic("failed to get reply content")
	}
	content := input.GetText()
	r.app.client.Reply(r.status, content)
	r.app.ViewTimeline()
	r.app.FocusTimeline()
}

func NewReplyForm(app *App) *ReplyForm {
	r := &ReplyForm{
		Form: cview.NewForm(),
		app:  app,
	}
	r.SetBackgroundColor(tcell.ColorDefault)
	r.SetFieldBackgroundColor(tcell.ColorDefault)
	r.SetFieldBackgroundColorFocused(tcell.ColorDefault)
	r.SetButtonBackgroundColor(tcell.ColorDefault)
	r.SetBorder(true)

	r.AddInputField("Reply", "", 80, nil, nil)
	r.AddButton("Send", r.sendReply)
	return r

}

func (r *ReplyForm) SetStatus(status *mastodon.Status) {
	r.status = status
}
