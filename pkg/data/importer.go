package data

import "embed"

//go:generate go run ../../cmd/txt2json/main.go -i DATA/konstytucjaRP.txt -o DATA/konstytucjaRP.json

//go:embed DATA
var data embed.FS
