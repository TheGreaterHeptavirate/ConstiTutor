package app

import (
	"bytes"
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/internal/assets"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"golang.org/x/image/colornames"
	"image"
	"image/png"
)

const (
	windowTitle              = "Consti Tutor"
	resolutionX, resolutionY = 800, 600
	searchButtonW            = 100
)

const logoLeftRightSpacing = 200

type App struct {
	window *giu.MasterWindow
	data   []*data.LegalAct

	searchPhrase string
	rows         []*giu.TableRowWidget
	logo         struct {
		w, h    int
		texture *giu.Texture
	}

	searchOptions struct {
		actNames   bool
		paragraphs bool
		text       bool
	}
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
	if err := giu.ParseCSSStyleSheet(assets.DefaultTheme); err != nil {
		panic(err)
	}

	a.Research("")

	// load/set tray icon
	icon, err := png.Decode(bytes.NewReader(assets.TrayIcon))
	if err != nil {
		panic(err)
	}

	a.window.SetIcon([]image.Image{icon})

	// load logo image
	logo, err := png.Decode(bytes.NewReader(assets.Logo))
	if err != nil {
		panic(err)
	}

	giu.EnqueueNewTextureFromRgba(logo, func(t *giu.Texture) {
		a.logo.texture = t
	})

	a.logo.w, a.logo.h = logo.Bounds().Dx(), logo.Bounds().Dy()

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
	// calculate logo H
	logoW := availableW - 2*logoLeftRightSpacing - spacingW
	logoH := int(float32(a.logo.h) / float32(a.logo.w) * logoW)

	giu.Layout{
		giu.Row(
			giu.Dummy(logoLeftRightSpacing, 0),
			giu.Image(a.logo.texture).Size(logoW, float32(logoH)),
		),
		giu.Row(
			giu.InputText(&a.searchPhrase).Size(availableW-searchButtonW-spacingW).OnChange(func() {
				a.Research(a.searchPhrase)
			}),
			giu.Button("Szukaj").Size(searchButtonW, 0).OnClick(func() {
				a.Research(a.searchPhrase)
			}),
		),
		giu.Row(
			giu.Checkbox("W nazwach aktów", &a.searchOptions.actNames),
			giu.Checkbox("W paragrafach", &a.searchOptions.paragraphs),
			giu.Checkbox("W treści", &a.searchOptions.text),
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

func (a *App) addRow(actName string, rule *data.Rule) {
	a.rows = append(a.rows, giu.TableRow(
		giu.Label(actName+" "+rule.Identifier),
		giu.Label(rule.Text),
	))
}

func (a *App) getMenubar() *giu.MenuBarWidget {
	return giu.MenuBar().Layout(
		giu.Menu("Plik").Layout(),
	)
}
