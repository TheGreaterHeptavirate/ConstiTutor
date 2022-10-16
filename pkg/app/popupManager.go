package app

import (
	"github.com/AllenDang/giu"
)

type PopupModal struct {
	id          string
	popup       *giu.PopupModalWidget
	isOpenInGiu bool
	isOpen      bool
}

func NewPopupModal(id string) *PopupModal {
	return &PopupModal{
		id:     id,
		popup:  giu.PopupModal(id),
		isOpen: false,
	}
}

func (p *PopupModal) Open() {
	p.isOpen = true
}

func (p *PopupModal) Close() {
	p.isOpen = false
}

func (p *PopupModal) Layout(l ...giu.Widget) *PopupModal {
	p.popup.Layout(l...)
	return p
}

func (p *PopupModal) Build() {
	if !p.isOpen {
		return
	}

	p.popup.IsOpen(&p.isOpen).Build()

	if !p.isOpenInGiu {
		giu.OpenPopup(p.id)
	}
}
