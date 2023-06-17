//go:build darwin
// +build darwin

package gsnake

import "golang.org/x/sys/unix"

const (
	getTermiosRequest = unix.TIOCGETA
	setTermiosRequest = unix.TIOCSETA
)
