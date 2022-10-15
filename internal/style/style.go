package style

import (
	_ "embed"
)

//go:embed stylesheet.css
var DefaultTheme []byte
