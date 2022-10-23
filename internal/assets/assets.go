/*
 * Copyright (c) 2022. The Greater Heptavirate (https://github.com/TheGreaterHeptavirate). All Rights Reserved
 *
 * All copies of this software (if not stated otherwise) are dedicated
 * ONLY to personal, non-commercial use.
 */

// Package assets uses go:embed to embed asset into the binary
// and shares it with other packages.
package assets

import (
	_ "embed"
)

// DefaultTheme sores a CSS stylesheet for giu application.
//
//go:embed css/stylesheet.css
var DefaultTheme []byte

// Icons.
var (
	//go:embed icons/1.png
	Logo []byte
	//go:embed icons/icon.png
	TrayIcon []byte
)

// ClickSound stores button click sound
//
//go:embed sounds/click.mp3
var ClickSound []byte

// fonts:
var (
	//go:embed fonts/times_new_roman/times_new_roman.ttf
	TimesNewRoman []byte
	//go:embed fonts/times_new_roman/times_new_roman_bold.ttf
	TimesNewRomanBold []byte
)
