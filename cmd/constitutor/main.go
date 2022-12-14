/*
 * Copyright (c) 2022. The Greater Heptavirate (https://github.com/TheGreaterHeptavirate). All Rights Reserved
 *
 * All copies of this software (if not stated otherwise) are dedicated
 * ONLY to personal, non-commercial use.
 */

// goversioninfo is https://github.com/josephspurrier/goversioninfo/
//
// //go:generate goversioninfo -icon=icon.ico -platform-specific=true
package main

import (
	"log"

	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Panic(err)
	}

	a.Run()
}
