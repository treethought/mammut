package ui

import (
	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

type Menu struct {
	*cview.List
	app *App
}

func NewMenu(app *App) *Menu {
	m := &Menu{
		List: cview.NewList(),
		app:  app,
	}
	m.SetBorder(true)
	m.SetBackgroundColor(tcell.ColorDefault)
	m.SetSelectedTextColor(tcell.ColorTeal)
	m.SetHighlightFullLine(false)

	for _, i := range TimelineTypes {
		li := cview.NewListItem(i.String())
		li.SetReference(i)
		m.AddItem(li)
	}

	m.SetSelectedFunc(func(idx int, li *cview.ListItem) {
		ref := m.GetCurrentItem().GetReference()
		ttype, ok := ref.(TimelineType)
		if !ok {
			return
		}
		m.app.timeline.SetTimeline(ttype)
		m.app.FocusTimeline()

	})
	return m
}
