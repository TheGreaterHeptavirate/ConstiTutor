package app

import (
	"bytes"
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/internal/assets"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"github.com/pkg/browser"
	"golang.org/x/image/colornames"
	"image"
	"image/png"
	"log"
)

const (
	windowTitle              = "Consti Tutor"
	resolutionX, resolutionY = 800, 600
	searchButtonW            = 100
)

const (
	logoLeftRightSpacing = 200
	aboutUsText          = `
Consti Tutor to program służący do wyszukiwania interesujących Cię
aktów prawnych w Konstytucji Rzeczypospolitej Polskiej i innych ustawach.

Wersja: v1.0
Autor: The Greater Heptavirate: programming lodge
Oficialna strona projektu: https://github.com/TheGreaterHeptavirate/ConstiTutor
`
	projectURL = "https://github.com/TheGreaterHeptavirate/ConstiTutor"
	bugURL     = "https://github.com/TheGreaterHeptavirate/ConstiTutor/issues/new"
)

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
		giu.PrepareMsgbox(),
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
		giu.Menu("Plik").Layout(
			giu.MenuItem("Zamknij").OnClick(func() {
				a.window.Close()
			}),
		),
		giu.Menu("Pomoc").Layout(
			giu.MenuItem("O programie").OnClick(func() {
				giu.Msgbox("O programie", aboutUsText)
			}),
			giu.MenuItem("Zobacz na GitHubie").OnClick(func() {
				if err := browser.OpenURL(projectURL); err != nil {
					a.ReportError(err)
				}
			}),
			giu.Separator(),
			giu.MenuItem("Zgłoś błąd").OnClick(func() {
				if err := browser.OpenURL(bugURL); err != nil {
					a.ReportError(err)
				}
			}),
		),
	)
}

// ReportError prints an error to the log and shows a message box in App.
func (a *App) ReportError(err error) {
	text := "Wystąpił nieznany błąd"
	if err != nil {
		text = fmt.Sprintf("Wystąpił błąd: %s\nProsimy skontaktować się z nami poprzez menu Pomoc->Zgłoś Błąd", err)
	}

	giu.Msgbox("Wystąpił błąd!", text)
	log.Print(err)
}
