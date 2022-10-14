package app

import (
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"golang.org/x/image/colornames"
	"strings"
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
	rows         []*giu.TableRowWidget
}

func New() (*App, error) {
	d, err := data.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize json data: %w", err)
	}

	return &App{
		data:   d,
		window: giu.NewMasterWindow(windowTitle, resolutionX, resolutionY, 0),
		rows:   make([]*giu.TableRowWidget, 0),
	}, nil
}

func (a *App) Run() {
	// initialization
	a.research("")

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
			giu.InputText(&a.searchPhrase).Size(availableW-searchButtonW-spacingW).OnChange(func() {
				a.research(a.searchPhrase)
			}),
			giu.Button("Szukaj").Size(searchButtonW, 0).OnClick(func() {
				a.research(a.searchPhrase)
			}),
		),
		giu.Label(""),
		giu.Condition(len(a.rows) > 0,
			giu.Layout{
				giu.Table().
					Columns(
						giu.TableColumn("Paragraf"),
						giu.TableColumn("Treść"),
					).
					Rows(a.rows...),
			},
			giu.Layout{
				giu.Style().SetColor(giu.StyleColorText, colornames.Gray).To(
					giu.Label("Brak wyników..."),
				),
			},
		),
	}.Build()
}

func (a *App) research(phrase string) {
	a.rows = make([]*giu.TableRowWidget, 0)

	for _, act := range a.data {
		for _, rule := range act.Rules {
			if phrase == "" || strings.Contains(rule.Text, phrase) {
				a.rows = append(a.rows, giu.TableRow(
					giu.Label(act.ActName+" "+rule.Identifier),
					giu.Label(rule.Text),
				))
			}
		}
	}
}

func (a *App) getMenubar() *giu.MenuBarWidget {
	return giu.MenuBar().Layout(
		giu.Menu("Plik").Layout(),
	)
}
