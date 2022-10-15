package assets

import (
	_ "embed"
)

//go:embed css/stylesheet.css
var DefaultTheme []byte
