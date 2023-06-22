package ssh

import (
	"time"

	"golang.org/x/crypto/ssh"
)

type SshInputReader struct {
	channel            ssh.Channel
	input              chan byte
	lastKeyPressedTime time.Time
}

func NewSshInputReader(channel ssh.Channel) *SshInputReader {
	s := &SshInputReader{channel, make(chan byte), time.Now()}
	go s.readInput()
	return s
}

// events poller
func (s *SshInputReader) Poll() byte {
	return <-s.input
}

func (s *SshInputReader) Close() {
	select {
	case _ = <-s.input:
		close(s.input)
	default:
	}
}

func (s *SshInputReader) readInput() {
	var buf [1]byte
	for {
		_, err := s.channel.Read(buf[:])
		if err != nil {
			close(s.input)
			return
		}
		s.lastKeyPressedTime = time.Now()
		s.input <- buf[0]
	}
}
