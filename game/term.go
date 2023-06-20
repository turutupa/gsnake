package gsnake

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type Term struct {
	sig   chan os.Signal
	input chan rune
}

func NewTerm() *Term {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	input := make(chan rune)
	go readInput(input)
	return &Term{
		sig:   sig,
		input: input,
	}
}

func (t *Term) PollEvents() rune {
	return <-t.input
}

func readInput(input chan<- rune) {
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
			close(input)
			return
		}
		input <- rune(buf[0])
		time.Sleep(40 * time.Millisecond)
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

func (t *Term) clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
