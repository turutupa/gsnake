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
	sig         chan os.Signal
	input       chan rune
	OnExit      func()
	OnForceExit func()
}

func NewTerm() *Term {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	input := make(chan rune)
	go readInput(input)
	defaultExitFn := func() {
		fmt.Println("Exiting!")
	}
	return &Term{
		sig:    sig,
		input:  input,
		OnExit: defaultExitFn,
	}
}

func (t *Term) PollEvents() rune {
	select {
	case <-t.sig:
		t.OnExit()
		os.Exit(0)
		return 'q'
	case r := <-t.input:
		if r == 'q' {
			if t.OnExit != nil {
				t.OnExit()
			}
			os.Exit(0)
		}
		return r
	}
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
		time.Sleep(30 * time.Millisecond)
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
