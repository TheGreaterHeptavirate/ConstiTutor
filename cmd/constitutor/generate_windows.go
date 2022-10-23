//go:build windows
// +build windows

/*
 * Copyright (c) 2022. The Greater Heptavirate (https://github.com/TheGreaterHeptavirate). All Rights Reserved
 *
 * All copies of this software (if not stated otherwise) are dedicated
 * ONLY to personal, non-commercial use.
 */

package main

//go:generate x86_64-w64-mingw32-windres *.rc -O coff -o windows.syso
