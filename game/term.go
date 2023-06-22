package gsnake

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

type Term struct {
	sig   chan os.Signal
	input chan byte
}

func NewTerm() *Term {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	input := make(chan byte)
	term := &Term{sig: sig, input: input}
	go term.readInput()
	return term
}

func (t *Term) Poll() (byte, error) {
	return <-t.input, nil
}

func (t *Term) Close() {
	select {
	case _, ok := <-t.input:
		if ok {
			close(t.input)
		}
	default:
	}
}

func (t *Term) readInput() {
	oldState, err := getTermios(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal attributes: ", err)
		return
	}
	defer setTermios(int(os.Stdin.Fd()), oldState)
	newState := *oldState
	newState.Lflag &^= syscall.ECHO | syscall.ICANON
	if err := setTermios(int(os.Stdin.Fd()), &newState); err != nil {
		fmt.Println("Error setting terminal attributes:", err)
		return
	}

	var buf [1]byte
	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil {
			close(t.input)
			return
		}
		t.input <- buf[0]
	}
}

func getTermios(fd int) (*unix.Termios, error) {
	termios, err := unix.IoctlGetTermios(fd, getTermiosRequest)
	if err != nil {
		return nil, fmt.Errorf("ioctl get termios: %v", err)
	}
	return termios, nil
}

func setTermios(fd int, termios *unix.Termios) error {
	if err := unix.IoctlSetTermios(fd, setTermiosRequest, termios); err != nil {
		return fmt.Errorf("ioctl set termios: %v", err)
	}
	return nil
}
