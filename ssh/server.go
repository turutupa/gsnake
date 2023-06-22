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
	"time"

	"turutupa/gsnake/events"
	"turutupa/gsnake/log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
)

const DEFAULT_PRIV_KEY_FILENAME = "gsnake_ed25519"
const DEFAULT_PUB_KEY_FILENAME = "gsnake_ed25519.pub"

type SshApp interface {
	Run()
}

type SshServer struct {
	port int
}

func NewSshServer(port int) *SshServer {
	return &SshServer{port}
}

func (s *SshServer) Run(sshAppInjector func(io.Writer, events.EventPoller) SshApp) {
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
			log.Error("Failed to handshake ", err)
			continue
		}
		username := sshConn.User()
		log.Info(username + " connected from " + netConn.RemoteAddr().String() + " " + string(sshConn.ClientVersion()))
		go ssh.DiscardRequests(reqs)
		go s.handleChannels(username, chans, sshAppInjector)
	}
}

func (s *SshServer) handleChannels(
	username string,
	chans <-chan ssh.NewChannel,
	sshAppInjector func(io.Writer, events.EventPoller) SshApp) {
	// Service the incoming Channel channel in go routine
	for newChannel := range chans {
		go s.handleChannel(username, newChannel, sshAppInjector) // propagating channel and sshApp
	}
}

func (s *SshServer) handleChannel(
	username string,
	newChannel ssh.NewChannel,
	sshAppInjector func(io.Writer, events.EventPoller) SshApp) {
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
	defer channel.Close()
	defer log.Info(username + " disconnected")

	// Set up terminal emulation
	term := terminal.NewTerminal(channel, "")
	sshInputReader := NewSshInputReader(channel)
	sshApp := sshAppInjector(term, sshInputReader)

	go s.activityMonitor(username, term, sshInputReader, channel)
	// Run SSH APP
	sshApp.Run()
	s.closeChannel(channel)
}

func (s *SshServer) activityMonitor(username string, term *term.Terminal, inputReader *SshInputReader, channel ssh.Channel) {
	idleTimeout := 5 * time.Minute
	checkTimeout := 1 * time.Minute
	for {
		select {
		case <-time.After(checkTimeout):
			if time.Since(inputReader.lastKeyPressedTime) > idleTimeout {
				term.Write([]byte("Session closed. Idle for too long (5 mins)."))
				s.closeChannel(channel)
				log.Info(username + " disconnected")
				return
			}
		}
	}
}

func (s *SshServer) closeChannel(channel ssh.Channel) {
	_, err := channel.SendRequest("dummy-request", true, nil)
	if err != nil { // Channel is closed
		return
	}
	channel.Close()
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
