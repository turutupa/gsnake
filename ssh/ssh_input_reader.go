package ssh

import (
	"time"
	"turutupa/gsnake/log"

	"golang.org/x/crypto/ssh"
)

type SshInputReader struct {
	channel ssh.Channel
	input   chan byte
}

func NewSshInputReader(channel ssh.Channel) *SshInputReader {
	s := &SshInputReader{channel, make(chan byte)}
	go s.readInput()
	return s
}

// events poller
func (s *SshInputReader) PollEvents() byte {
	return <-s.input
}

func (s *SshInputReader) readInput() {
	var buf [1]byte
	for {
		_, err := s.channel.Read(buf[:])
		if err != nil {
			close(s.input)
			log.Error("Could not read from channel", err)
			return
		}
		s.input <- buf[0]
		time.Sleep(40 * time.Millisecond)
	}
}
