package app

import (
	"github.com/AllenDang/giu"
	"github.com/pkg/browser"
	"golang.org/x/image/colornames"
	"image/color"
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
						childH := availableH - spacingH*2 - buttonH

						giu.Child().Layout(
							a.renderLogo(40),
							giu.Markdown(&aboutUs),
						).Size(-1, childH).Build()
					}),
				),
			giu.Separator(),
			giu.Row(
				giu.Button("Zgłoś błąd").OnClick(a.reportBug).Size(0, buttonH),
				a.playClickSound(),
				giu.Button("Zamknij").OnClick(a.aboutAppPopup.Close).Size(0, buttonH),
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
				a.Research(a.searchPhrase)
			}),
			a.playClickSound(),
		),
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