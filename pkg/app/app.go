package app

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"image"
	"image/png"
	"log"
	"runtime"

	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/internal/assets"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"github.com/hajimehoshi/go-mp3"
)

const (
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
	if err := a.InitializeSound(); err != nil {
		log.Panic(err)
	}

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
