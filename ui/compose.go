package ui

import (
	"github.com/gdamore/tcell/v2"
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

			c.app.FocusTimeline()
		}
	})

	return c
}
