package app

import (
	"fmt"
	"io"
	"log"

	"github.com/AllenDang/giu"
	"github.com/pkg/browser"
)

// ReportError prints an error to the log and shows a message box in App.
func (a *App) ReportError(err error) {
	text := "Wystąpił nieznany błąd"
	if err != nil {
		text = fmt.Sprintf("Wystąpił błąd: %s\nProsimy skontaktować się z nami poprzez menu Pomoc->Zgłoś Błąd", err)
	}

	giu.Msgbox("Wystąpił błąd!", text)
	log.Print(err)
}

func (a *App) reportBug() {
	if err := browser.OpenURL(bugURL); err != nil {
		a.ReportError(err)
	}
}

func (a *App) playClickSound() giu.Widget {
	return giu.Event().OnMouseDown(giu.MouseButtonLeft, func() {
		newPos, err := a.clickSound.(io.Seeker).Seek(0, io.SeekStart)
		if err != nil {
			a.ReportError(err)
		}

		if newPos != 0 {
			//nolint:goerr113 // it's just for ReportError call
			a.ReportError(fmt.Errorf("failed to seek to the beginning of the click sound"))
		}

		a.clickSound.Play()
	})
}
