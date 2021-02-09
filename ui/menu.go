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

	items := []string{
		"local",
		"public",
		"liked",
		"profile",
		"tags",
	}
	for _, i := range items {
		li := cview.NewListItem(i)
		m.AddItem(li)
	}
	return m
}
