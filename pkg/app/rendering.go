package app

import (
	"image/color"

	"github.com/AllenDang/giu"
	"github.com/pkg/browser"
	"golang.org/x/image/colornames"
)

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
						childH := availableH - spacingH*2 - aboutUsDialogueButtonH

						giu.Child().Layout(
							a.renderLogo(40),
							giu.Markdown(&aboutUs),
						).Size(-1, childH).Build()
					}),
				),
			giu.Separator(),
			giu.Row(
				giu.Button("Zgłoś błąd").OnClick(a.reportBug).Size(0, aboutUsDialogueButtonH),
				a.playClickSound(),
				giu.Button("Zamknij").OnClick(a.aboutAppPopup.Close).Size(0, aboutUsDialogueButtonH),
				a.playClickSound(),
			),
		),
		giu.PrepareMsgbox(),
		giu.Custom(a.renderMainView),
	)
}

func (a *App) getMenubar() *giu.MenuBarWidget {
	return giu.MenuBar().Layout(
		giu.Menu("Plik").Layout(
			giu.MenuItem("Zamknij").Shortcut("Ctrl+Q").OnClick(a.window.Close),
			a.playClickSound(),
		),
		a.playClickSound(),
		giu.Menu("Pomoc").Layout(
			giu.MenuItem("O programie").Shortcut("F1").OnClick(a.aboutAppPopup.Open),
			a.playClickSound(),
			giu.MenuItem("Zobacz na GitHubie").OnClick(func() {
				if err := browser.OpenURL(projectURL); err != nil {
					a.ReportError(err)
				}
			}),
			a.playClickSound(),
			giu.Separator(),
			giu.MenuItem("Zgłoś błąd").OnClick(a.reportBug),
			a.playClickSound(),
		),
		a.playClickSound(),
	)
}

func (a *App) renderLogo(percentageH int) giu.Widget {
	return giu.Custom(func() {
		availableW, availableH := giu.GetAvailableRegion()
		spacingW, _ := giu.GetItemSpacing()
		scale := float32(a.logo.h) / float32(a.logo.w)
		maxLogoH := availableH * float32(percentageH) / 100
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
			giu.Custom(a.logo.loadTextureMutex.Lock),
			giu.Image(a.logo.texture).Size(logoW, logoH),
			giu.Custom(a.logo.loadTextureMutex.Unlock),
		).Build()
	})
}

func (a *App) renderMainView() {
	availableW, _ := giu.GetAvailableRegion()
	spacingW, _ := giu.GetItemSpacing()
	searchFieldW := (availableW)*searchPercentage/100 - spacingW/2
	searchButtonW := availableW - searchFieldW - spacingW/2

	giu.Layout{
		a.renderLogo(logoHPercentage),
		giu.Row(
			giu.InputText(&a.searchPhrase).Size(searchFieldW).OnChange(func() {
				a.Research(a.searchPhrase)
			}).Hint("Szukaj..."),
			giu.Layout{
				giu.Button("Szukaj").Size(searchButtonW, 0).OnClick(func() {
					a.Research(a.searchPhrase)
				}),
				a.playClickSound(),
			},
		),
		// workaround; see https://github.com/AllenDang/giu/issues/572
		giu.Dummy(0, 0),
		giu.Row(
			giu.Checkbox("W nazwach aktów", &a.searchOptions.actNames).OnChange(func() {
				a.Research(a.searchPhrase)
			}),
			a.playClickSound(),
			giu.Checkbox("W paragrafach", &a.searchOptions.paragraphs).OnChange(func() {
				a.Research(a.searchPhrase)
			}),
			a.playClickSound(),
			giu.Checkbox("W treści", &a.searchOptions.text).OnChange(func() {
				a.Research(a.searchPhrase)
			}),
			a.playClickSound(),
		),
		// workaround; see https://github.com/AllenDang/giu/issues/572
		giu.Dummy(0, 0),
		giu.Condition(len(a.rows) > 0,
			giu.Layout{
				giu.Child().Layout(
					giu.Table().Flags(
						giu.TableFlagsScrollY|
							giu.TableFlagsScrollX|
							giu.TableFlagsBordersInner|
							giu.TableFlagsBordersInnerH|
							giu.TableFlagsResizable,
					).
						Columns(
							giu.TableColumn("Paragraf").
								Flags(giu.TableColumnFlagsWidthStretch).
								InnerWidthOrWeight(.3),
							giu.TableColumn("Treść").
								Flags(giu.TableColumnFlagsWidthStretch).
								InnerWidthOrWeight(.7),
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
