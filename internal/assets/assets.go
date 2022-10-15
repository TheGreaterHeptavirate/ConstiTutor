package assets

import (
	_ "embed"
)

//go:embed css/stylesheet.css
var DefaultTheme []byte

//go:embed icons/logo.png
var Logo []byte
