package app

import "github.com/AllenDang/giu"

type PopupModal struct {
	id         string
	popup      *giu.PopupModalWidget
	isOpen     bool
	shouldOpen bool
}

func NewPopupModal(id string) *PopupModal {
	return &PopupModal{
		id:         id,
		popup:      giu.PopupModal(id),
		isOpen:     false,
		shouldOpen: false,
	}
}

func (p *PopupModal) Open() {
	if !p.isOpen {
		p.shouldOpen = true
	}
}

func (p *PopupModal) Close() {
	if p.isOpen {
		p.isOpen = false
		p.shouldOpen = false
		giu.CloseCurrentPopup()
	}
}

func (p *PopupModal) Layout(l ...giu.Widget) *PopupModal {
	p.popup.Layout(l...)
	return p
}

func (p *PopupModal) Build() {
	if p.shouldOpen {
		p.isOpen = true
		p.shouldOpen = false
		giu.OpenPopup(p.id)
	}

	p.popup.Build()
}
