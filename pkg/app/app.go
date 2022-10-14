package app

import (
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
)

const (
	windowTitle              = "Consti Tutor"
	resolutionX, resolutionY = 800, 600
	searchButtonW            = 100
)

type App struct {
	window *giu.MasterWindow
	data   []*data.LegalAct

	searchPhrase string
}

func New() (*App, error) {
	d, err := data.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize json data: %w", err)
	}

	return &App{
		data:   d,
		window: giu.NewMasterWindow(windowTitle, resolutionX, resolutionY, 0),
	}, nil
}

func (a *App) Run() {
	// initialization

	// run render loop
	a.window.Run(a.render)
}

func (a *App) render() {
	giu.SingleWindowWithMenuBar().Layout(
		a.getMenubar(),
		giu.Custom(a.renderMainView),
	)
}

func (a *App) renderMainView() {
	availableW, _ := giu.GetAvailableRegion()
	spacingW, _ := giu.GetItemSpacing()
	giu.Layout{
		giu.Row(
			giu.InputText(&a.searchPhrase).Size(availableW-searchButtonW-spacingW),
			giu.Button("Szukaj").Size(searchButtonW, 0).OnClick(func() {

			}),
		),
	}.Build()
}

func (a *App) getMenubar() *giu.MenuBarWidget {
	return giu.MenuBar().Layout(
		giu.Menu("Plik").Layout(),
	)
}
