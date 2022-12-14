/*
 * Copyright (c) 2022. The Greater Heptavirate (https://github.com/TheGreaterHeptavirate). All Rights Reserved
 *
 * All copies of this software (if not stated otherwise) are dedicated
 * ONLY to personal, non-commercial use.
 */

// Package app contains the main UI logic for ConstiTutor.
// For stylesheet see internal/assets/css.
package app

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"runtime"
	"sync"

	"github.com/hajimehoshi/oto/v2"

	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/internal/assets"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"github.com/hajimehoshi/go-mp3"
)

const (
	// sound player config (for oto).
	channelCount    = 2
	bitDepthInBytes = 2

	// App config.
	windowTitle              = "ConstiTutor"
	resolutionX, resolutionY = 800, 600
	logoHPercentage          = 15
	searchPercentage         = 80

	// about us dialog.
	aboutUsText = `
ConstiTutor to program służący do wyszukiwania interesujących Cię
aktów prawnych w Konstytucji Rzeczypospolitej Polskiej i innych ustawach.

Wersja: v1.1
Autorzy: [The Greater Heptavirate: programming lodge](https://github.com/TheGreaterHeptavirate)
[Oficialna strona projektu](https://github.com/TheGreaterHeptavirate/ConstiTutor)
`
	projectURL             = "https://github.com/TheGreaterHeptavirate/ConstiTutor"
	bugURL                 = "https://github.com/TheGreaterHeptavirate/ConstiTutor/issues/new"
	aboutUsDialogueButtonH = 30
)

// App represents a UI application
// Create a new instance with New and run it with Run
// NOTE! only one instance of an app could be active at once!
type App struct {
	window *giu.MasterWindow
	data   []*data.LegalAct

	searchPhrase string
	rows         []*giu.TableRowWidget
	logo         struct {
		w, h             int
		texture          *giu.Texture
		loadTextureMutex *sync.Mutex
	}

	searchOptions struct {
		actNames   bool
		paragraphs bool
		text       bool
	}

	aboutAppPopup *PopupModal
	clickSound    oto.Player
}

// New creates a new instance of an app.
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
	result.logo.loadTextureMutex = &sync.Mutex{}

	return result, nil
}

// Run starts the app.
// It'll hold program execution until the app is closed.
// You can call it in goroutine and use channels to communicate with it.
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

	a.initializeFonts()
	// load logo image
	logo, err := png.Decode(bytes.NewReader(assets.Logo))
	if err != nil {
		panic(err)
	}

	giu.EnqueueNewTextureFromRgba(logo, func(t *giu.Texture) {
		a.logo.loadTextureMutex.Lock()
		a.logo.texture = t
		a.logo.loadTextureMutex.Unlock()
	})

	a.logo.w, a.logo.h = logo.Bounds().Dx(), logo.Bounds().Dy()

	// run render loop
	a.window.Run(a.render)
}

// InitializeSound initializes sound player (oto).
func (a *App) InitializeSound() error {
	mp3Data, err := mp3.NewDecoder(bytes.NewReader(assets.ClickSound))
	if err != nil {
		return fmt.Errorf("decoding MP3 file data: %w", err)
	}

	c, ready, err := oto.NewContext(mp3Data.SampleRate(), channelCount, bitDepthInBytes)
	if err != nil {
		return fmt.Errorf("creating oto context: %w", err)
	}

	<-ready

	a.clickSound = c.NewPlayer(mp3Data)
	runtime.KeepAlive(a.clickSound)

	return nil
}

func (a *App) initializeFonts() {
	giu.Context.FontAtlas.SetDefaultFontFromBytes(assets.TimesNewRoman, 14)
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
