package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"turutupa/gsnake/log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const DEFAULT_PRIV_KEY_FILENAME = "gsnake_ed25519"
const DEFAULT_PUB_KEY_FILENAME = "gsnake_ed25519.pub"

type Runnable interface {
	Run()
}

type SshServer struct {
	port int
}

func NewSshServer(port int, requestHandler func(io.Writer) Runnable) *SshServer {
	return &SshServer{port}
}

func (s *SshServer) Run() {
	privateBytes, _, ok := s.getKeyPairOrDefault(DEFAULT_PRIV_KEY_FILENAME, DEFAULT_PUB_KEY_FILENAME)
	if !ok {
		log.Error("No crypto keys found", nil)
		return
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Error("Failed to parse private key: %v", err)
		return
	}

	// Create SSH server configuration
	config := &ssh.ServerConfig{NoClientAuth: true}
	config.AddHostKey(privateKey)

	// Start SSH server
	strPort := strconv.Itoa(s.port)
	addr := "0.0.0.0:" + strPort
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed to listen on "+strPort, err)
	} else {
		log.Info("Listening on port " + strPort)
	}

	for {
		netConn, err := listener.Accept()
		if err != nil {
			log.Error("Failed to accept incoming connection (%s)", err)
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn
		sshConn, chans, reqs, err := ssh.NewServerConn(netConn, config)
		if err != nil {
			log.Error("Failed to handshake (%s)", err)
			continue
		}
		log.Info("User connected from " + netConn.RemoteAddr().String() + " " + string(sshConn.ClientVersion()))
		go ssh.DiscardRequests(reqs)
		go s.handleChannels(chans)
	}
}

func (s *SshServer) handleChannels(chans <-chan ssh.NewChannel) {
	// Service the incoming Channel channel in go routine
	for newChannel := range chans {
		go s.handleChannel(newChannel)
	}
}

func (s *SshServer) handleChannel(newChannel ssh.NewChannel) {
	// Channels have a type, depending on the application level protocol intended.
	// In the case of a shell, the type is "session" and ServerShell may be used to present a simple terminal interface.
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}

	channel, _, err := newChannel.Accept()
	if err != nil {
		log.Error("Could not accept channel", err)
		return
	}

	// Set up terminal emulation
	term := terminal.NewTerminal(channel, "")

	// Start your snake game in a goroutine
	go func() {
		// Run your snake game
		// You can write game output to the terminal using the term.Write() function
		// and read user input from the terminal using the term.ReadLine() function
		term.Write([]byte("hello everyone!"))
	}()

	go func() {
		var buf [1]byte
		for {
			_, err := channel.Read(buf[:])
			if err != nil {
				log.Error("Could not read from channel", err)
				return
			}
			// Use the input character, stored in buf[0]
		}
	}()
}

func (s *SshServer) getKeyPairOrDefault(privateKeyPath, publicKeyPath string) ([]byte, []byte, bool) {
	// Check if files already exist
	sshDir, _ := os.Getwd()
	sshDir = filepath.Join(sshDir, ".ssh")
	privateKeyPath = filepath.Join(sshDir, privateKeyPath)
	publicKeyPath = filepath.Join(sshDir, publicKeyPath)
	if _, err := os.Stat(privateKeyPath); err == nil {
		if _, err := os.Stat(publicKeyPath); err == nil {
			// Both keys exist, no need to create new ones
			privateBytes, _ := os.ReadFile(privateKeyPath)
			publicBytes, _ := os.ReadFile(publicKeyPath)
			return privateBytes, publicBytes, true
		}
	}

	// generate .ssh dir if required
	os.MkdirAll(sshDir, 0700)

	// Generate a new ED25519 key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, false
	}
	// Prepare the private key for writing to a file
	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return []byte{}, []byte{}, false
	}
	// Prepare the private key for writing to a file
	privateBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}
	privateBytes := pem.EncodeToMemory(privateBlock)
	if err := os.WriteFile(privateKeyPath, privateBytes, 0600); err != nil {
		return nil, nil, false
	}
	// Prepare the public key for writing to a file
	publicSSHKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return nil, nil, false
	}
	publicBytes := ssh.MarshalAuthorizedKey(publicSSHKey)
	if err := ioutil.WriteFile(publicKeyPath, publicBytes, 0644); err != nil {
		return nil, nil, false
	}
	return privateBytes, publicBytes, true
}
