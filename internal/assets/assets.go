package assets

import (
	_ "embed"
)

//go:embed css/stylesheet.css
var DefaultTheme []byte

var (
	//go:embed icons/logo.png
	Logo []byte
	//go:embed icons/icon.png
	TrayIcon []byte
)
