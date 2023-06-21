package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"turutupa/gsnake/events"
	"turutupa/gsnake/game"
	"turutupa/gsnake/ssh"
)

const SSH_MODE = "ssh"
const PORT_FLAG = "p"
const MODE_FLAG = "m"
const HELP_FLAG = "h"
const DEFAULT_PORT = 5555

func main() {
	port := flag.Int(PORT_FLAG, DEFAULT_PORT, "Port number. Only used in 'ssh' mode")
	mode := flag.String(MODE_FLAG, "local", "Expected values are 'local' or 'ssh'")
	help := flag.Bool(HELP_FLAG, false, "Display help information.")
	flag.Usage = displayHelp
	flag.Parse()

	if *help {
		displayHelp()
	}
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

func displayHelp() {
	fmt.Println("Usage: gsnake [-m mode] [-p port]")
	fmt.Println("\nOptions:")
	fmt.Println("  -m, --mode\n\t(Optional) Expected values are 'local' or 'ssh'. If -m flag is set to 'ssh' then it will host a gsnake ssh server. Defaults to 'local'.")
	fmt.Println("  -p, --port\n\t(Optional) Port number. Only used in 'ssh' mode. Defaults to 5555.")
	fmt.Println("  -h, --help\n\tDisplay this help text.")
	os.Exit(0)
}
