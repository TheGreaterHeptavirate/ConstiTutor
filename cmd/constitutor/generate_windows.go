//go:build windows
// +build windows

package main

//go:generate x86_64-w64-mingw32-windres *.rc -O coff -o windows.syso
