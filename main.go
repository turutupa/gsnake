package main

import (
	"flag"
	"io"
	"os"

	"turutupa/gsnake/events"
	"turutupa/gsnake/game"
	"turutupa/gsnake/ssh"
)

const SSH_MODE = "ssh"
const PORT_FLAG = "p"
const MODE_FLAG = "m"
const DEFAULT_PORT = 5555

func main() {
	port := flag.Int(PORT_FLAG, DEFAULT_PORT, "Port number. Only used in 'ssh' mode")
	mode := flag.String(MODE_FLAG, "local", "Expected values are 'local' or 'ssh'")
	flag.Parse()
	if *mode == SSH_MODE {
		sshServer := ssh.NewSshServer(*port)
		sshServer.Run(snakeApp)
	} else {
		snakeApp(os.Stdout, gsnake.NewTerm()).Run()
	}
}

func snakeApp(writer io.Writer, eventsPoller events.EventPoller) ssh.SshApp {
	rows := 20
	cols := 50
	screen := gsnake.NewScreen(writer, rows, cols)
	scoreboard, _ := gsnake.NewScoreboard()
	snake := gsnake.NewSnake(screen)
	fruit := gsnake.NewFruit(rows, cols)
	return gsnake.NewGame(eventsPoller, screen, scoreboard, fruit, snake)
}
