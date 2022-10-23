/*
 * Copyright (c) 2022. The Greater Heptavirate (https://github.com/TheGreaterHeptavirate). All Rights Reservet
 *
 * All copies of this software (if not stated otherway) are dedicated
 * ONLY to personal, non-commercial use.
 */

package app

import (
	"github.com/AllenDang/giu"
)

// PopupModal is a wrapper of giu.PopupModalWidget.
// It makes it possible to use Open / Close methods.
type PopupModal struct {
	id          string
	popup       *giu.PopupModalWidget
	isOpenInGiu bool
	isOpen      bool
}

// NewPopupModal creates a new instance of PopupModal.
func NewPopupModal(id string) *PopupModal {
	return &PopupModal{
		id:     id,
		popup:  giu.PopupModal(id),
		isOpen: false,
	}
}

// Open opens the popup (if not opened).
func (p *PopupModal) Open() {
	p.isOpen = true
}

// Close closes the popup (if opened).
func (p *PopupModal) Close() {
	p.isOpen = false
}

// Layout sets popup's layout (wraps (*giu.PopupModalWidget).Layout.
func (p *PopupModal) Layout(l ...giu.Widget) *PopupModal {
	p.popup.Layout(l...)

	return p
}

// Build implements giu.widget.
func (p *PopupModal) Build() {
	if !p.isOpen {
		return
	}

	p.popup.IsOpen(&p.isOpen).Build()

	if !p.isOpenInGiu {
		giu.OpenPopup(p.id)
	}
}
