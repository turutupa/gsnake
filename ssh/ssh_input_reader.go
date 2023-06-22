package ssh

import (
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SshInputReader struct {
	channel            ssh.Channel
	input              chan byte
	errChan            chan error
	lastKeyPressedTime time.Time
	mu                 sync.Mutex
	running            bool
}

func NewSshInputReader(channel ssh.Channel) *SshInputReader {
	s := &SshInputReader{
		channel:            channel,
		input:              make(chan byte),
		errChan:            make(chan error, 1),
		lastKeyPressedTime: time.Now(),
		running:            true,
	}
	go s.readInput()
	return s
}

// events poller
func (s *SshInputReader) Poll() (byte, error) {
	select {
	case input := <-s.input:
		return input, nil
	case err := <-s.errChan:
		return 0, err
	}
}

func (s *SshInputReader) Close() {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
}

func (s *SshInputReader) readInput() {
	var buf [1]byte
	for {
		s.mu.Lock()
		running := s.running
		s.mu.Unlock()
		if !running {
			break
		}
		_, err := s.channel.Read(buf[:])
		if err != nil {
			s.errChan <- err
			break
		}
		s.lastKeyPressedTime = time.Now()
		s.input <- buf[0]
	}
	close(s.input)
	close(s.errChan)
}
