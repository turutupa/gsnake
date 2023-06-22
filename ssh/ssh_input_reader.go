package ssh

import (
	"time"

	"golang.org/x/crypto/ssh"
)

type SshInputReader struct {
	channel            ssh.Channel
	input              chan byte
	lastKeyPressedTime time.Time
	running            bool
}

func NewSshInputReader(channel ssh.Channel) *SshInputReader {
	s := &SshInputReader{channel, make(chan byte), time.Now(), true}
	go s.readInput()
	return s
}

// events poller
func (s *SshInputReader) Poll() byte {
	return <-s.input
}

func (s *SshInputReader) Close() {
	s.running = false
}

func (s *SshInputReader) readInput() {
	var buf [1]byte
	for s.running {
		_, err := s.channel.Read(buf[:])
		if err != nil {
			close(s.input)
			return
		}
		s.lastKeyPressedTime = time.Now()
		s.input <- buf[0]
	}
	close(s.input)
}
