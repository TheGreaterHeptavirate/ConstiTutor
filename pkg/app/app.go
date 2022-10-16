package app

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"runtime"

	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/internal/assets"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"github.com/hajimehoshi/go-mp3"
	"github.com/pkg/browser"
)

const (
	sampleRate      = 4800
	channelCount    = 2
	bitDepthInBytes = 2

	windowTitle              = "Consti Tutor"
	resolutionX, resolutionY = 800, 600
	logoHProcentage          = 15
	searchProcentage         = 80
	aboutUsText              = `
ConstiTutor to program służący do wyszukiwania interesujących Cię
aktów prawnych w Konstytucji Rzeczypospolitej Polskiej i innych ustawach.

Wersja: v1.0
Autor: [The Greater Heptavirate: programming lodge](https://github.com/TheGreaterHeptavirate)
[Oficialna strona projektu](https://github.com/TheGreaterHeptavirate/ConstiTutor)
`
	projectURL = "https://github.com/TheGreaterHeptavirate/ConstiTutor"
	bugURL     = "https://github.com/TheGreaterHeptavirate/ConstiTutor/issues/new"
	buttonH    = 30
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

	aboutAppPopup *PopupModal
	clickSound    oto.Player
}

func New() (*App, error) {
	d, err := data.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize json data: %w", err)
	}

	result := &App{
		data:          d,
		window:        giu.NewMasterWindow(windowTitle, resolutionX, resolutionY, 0),
		rows:          make([]*giu.TableRowWidget, 0),
		aboutAppPopup: NewPopupModal("O Aplikacji"),
	}

	result.searchOptions.text = true

	return result, nil
}

func (a *App) Run() {
	// initialization
	a.InitializeSound()
	a.registerShortcuts()

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

func (a *App) InitializeSound() error {
	mp3Data, err := mp3.NewDecoder(bytes.NewReader(assets.ClickSound))
	if err != nil {
		return err
	}

	c, ready, err := oto.NewContext(mp3Data.SampleRate(), channelCount, bitDepthInBytes)
	if err != nil {
		return err
	}

	<-ready

	a.clickSound = c.NewPlayer(mp3Data)
	runtime.KeepAlive(a.clickSound)

	return nil
}

func (a *App) playClickSound() {
	newPos, err := a.clickSound.(io.Seeker).Seek(0, io.SeekStart)
	if err != nil {
		a.ReportError(err)
	}

	if newPos != 0 {
		a.ReportError(fmt.Errorf("failed to seek to the beginning of the click sound"))
	}

	a.clickSound.Play()
}

func (a *App) registerShortcuts() {
	a.window.RegisterKeyboardShortcuts(
		// quit - Ctrl+Q
		giu.WindowShortcut{
			Key:      giu.KeyQ,
			Modifier: giu.ModControl,
			Callback: a.window.Close,
		},
		// close popup - Esc
		giu.WindowShortcut{
			Key:      giu.KeyEscape,
			Modifier: giu.ModNone,
			Callback: func() {
				giu.CloseCurrentPopup()
				a.aboutAppPopup.Close()
			},
		},
		// about app - F1
		giu.WindowShortcut{
			Key:      giu.KeyF1,
			Modifier: giu.ModNone,
			Callback: a.aboutAppPopup.Open,
		},
	)
}

func (a *App) render() {
	aboutUs := aboutUsText
	giu.SingleWindowWithMenuBar().Layout(
		a.getMenubar(),
		giu.Custom(func() {
			giu.SetNextWindowSize(700, 550)
		}),
		a.aboutAppPopup.Layout(
			giu.Style().
				SetColor(giu.StyleColorChildBg, color.RGBA{}).
				SetColor(giu.StyleColorBorder, color.RGBA{}).
				To(
					giu.Custom(func() {
						_, availableH := giu.GetAvailableRegion()
						_, spacingH := giu.GetItemSpacing()
						childH := availableH - spacingH*2 - buttonH

						giu.Child().Layout(
							a.renderLogo(40),
							giu.Markdown(&aboutUs),
						).Size(-1, childH).Build()
					}),
				),
			giu.Separator(),
			giu.Row(
				giu.Button("Zgłoś błąd").OnClick(func() {
					a.playClickSound()
					a.reportBug()
				}).Size(0, buttonH),
				giu.Button("Zamknij").OnClick(func() {
					a.playClickSound()
					a.aboutAppPopup.Close()
				}).Size(0, buttonH),
			),
		),
		giu.PrepareMsgbox(),
		giu.Custom(a.renderMainView),
	)
}

func (a *App) renderMainView() {
	availableW, _ := giu.GetAvailableRegion()
	spacingW, _ := giu.GetItemSpacing()
	searchFieldW := (availableW)*searchProcentage/100 - spacingW/2
	searchButtonW := availableW - searchFieldW - spacingW/2

	giu.Layout{
		a.renderLogo(logoHProcentage),
		giu.Row(
			giu.InputText(&a.searchPhrase).Size(searchFieldW).OnChange(func() {
				a.Research(a.searchPhrase)
			}),
			giu.Button("Szukaj").Size(searchButtonW, 0).OnClick(func() {
				a.playClickSound()
				a.Research(a.searchPhrase)
			}),
		),
		giu.Row(
			giu.Checkbox("W nazwach aktów", &a.searchOptions.actNames).OnChange(a.playClickSound),
			giu.Checkbox("W paragrafach", &a.searchOptions.paragraphs).OnChange(a.playClickSound),
			giu.Checkbox("W treści", &a.searchOptions.text).OnChange(a.playClickSound),
		),
		giu.Label(""),
		giu.Condition(len(a.rows) > 0,
			giu.Layout{
				giu.Child().Layout(
					giu.Table().Flags(
						giu.TableFlagsScrollY|
							giu.TableFlagsBordersInner|
							giu.TableFlagsBordersInnerH,
					).
						Columns(
							giu.TableColumn("Paragraf").
								Flags(giu.TableColumnFlagsWidthStretch).
								InnerWidthOrWeight(.3),
							giu.TableColumn("Treść"),
						).
						Rows(a.rows...),
				),
			},
			giu.Layout{
				giu.Style().SetColor(giu.StyleColorText, colornames.Gray).To(
					giu.Label("Brak wyników..."),
				),
			},
		),
	}.Build()
}

func (a *App) renderLogo(procentage int) giu.Widget {
	return giu.Custom(func() {
		availableW, availableH := giu.GetAvailableRegion()
		spacingW, _ := giu.GetItemSpacing()
		scale := float32(a.logo.h) / float32(a.logo.w)
		maxLogoH := availableH * float32(procentage) / 100
		maxLogoW := availableW
		var logoW, logoH, dummyW float32
		if maxLogoW >= maxLogoH/scale {
			logoH = maxLogoH
			logoW = logoH / scale
			dummyW = (availableW-logoW)/2 - spacingW
			if dummyW < 0 {
				dummyW = 0
			}
		} else {
			logoW = maxLogoW
			logoH = logoW * scale
		}

		giu.Row(
			giu.Dummy(dummyW, 0),
			giu.Image(a.logo.texture).Size(logoW, float32(logoH)),
		).Build()
	})
}

func (a *App) reportBug() {
	if err := browser.OpenURL(bugURL); err != nil {
		a.ReportError(err)
	}
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
			giu.MenuItem("Zamknij").Shortcut("Ctrl+Q").OnClick(func() {
				a.playClickSound()
				a.window.Close()
			}),
		),
		giu.Event().OnClick(giu.MouseButtonLeft, func() {
			a.playClickSound()
		}),
		giu.Menu("Pomoc").Layout(
			giu.MenuItem("O programie").Shortcut("F1").OnClick(func() {
				a.playClickSound()
				a.aboutAppPopup.Open()
			}),
			giu.MenuItem("Zobacz na GitHubie").OnClick(func() {
				a.playClickSound()
				if err := browser.OpenURL(projectURL); err != nil {
					a.ReportError(err)
				}
			}),
			giu.Separator(),
			giu.MenuItem("Zgłoś błąd").OnClick(func() {
				a.playClickSound()
				a.reportBug()
			}),
		),
		giu.Event().OnClick(giu.MouseButtonLeft, func() {
			a.playClickSound()
		}),
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
