//go:build linux
// +build linux

package gsnake

import "golang.org/x/sys/unix"

const (
	getTermiosRequest = unix.TCGETS
	setTermiosRequest = unix.TCSETS
)
