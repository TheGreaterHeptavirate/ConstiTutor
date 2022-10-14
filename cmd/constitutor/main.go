package main

import (
	"fmt"

	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
)

func main() {
	d, err := data.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println(d[0])
}
