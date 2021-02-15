package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

type InfoView struct {
	*cview.TextView
	app *App
}

func NewInfoView(app *App) *InfoView {
	i := &InfoView{
		TextView: cview.NewTextView(),
		app:      app,
	}

	i.SetBackgroundColor(tcell.ColorDefault)
	i.SetBorder(true)
	info := fmt.Sprintf("%s\n%s",
		app.client.Account().DisplayName,
		app.client.Server(),
	)
	i.SetText(info)

	return i
}
