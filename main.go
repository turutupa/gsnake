package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"turutupa/gsnake/events"
	"turutupa/gsnake/game"
	"turutupa/gsnake/log"
	"turutupa/gsnake/ssh"
)

const SSH_MODE = "ssh"
const PORT_FLAG_SHORT = "p"
const PORT_FLAG_LONG = "port"
const MODE_FLAG_SHORT = "m"
const MODE_FLAG_LONG = "mode"
const LOG_FLAG_SHORT = "l"
const LOG_FLAG_LONG = "log"
const HELP_FLAG_SHORT = "h"
const HELP_FLAG_LONG = "help"
const DEFAULT_PORT = 5555

func main() {
	var port int
	var mode string
	var logging bool
	var help bool
	flag.IntVar(&port, PORT_FLAG_SHORT, DEFAULT_PORT, "")
	flag.IntVar(&port, PORT_FLAG_LONG, DEFAULT_PORT, "")
	flag.StringVar(&mode, MODE_FLAG_SHORT, "local", "")
	flag.StringVar(&mode, MODE_FLAG_LONG, "local", "")
	flag.BoolVar(&logging, LOG_FLAG_SHORT, false, "")
	flag.BoolVar(&logging, LOG_FLAG_LONG, false, "")
	flag.BoolVar(&help, HELP_FLAG_SHORT, false, "")
	flag.BoolVar(&help, HELP_FLAG_LONG, false, "")

	flag.Usage = displayHelp
	flag.Parse()

	if logging {
		log.EnableStorage()
	}
	if help {
		displayHelp()
	}
	if mode == SSH_MODE {
		sshServer := ssh.NewSshServer(port)
		sshServer.Run(snakeApp)
	} else {
		snakeApp(os.Stdout, gsnake.NewTerm()).Run()
	}
}

func snakeApp(writer io.Writer, eventsPoller events.EventPoller) ssh.SshApp {
	rows := 20
	cols := 50
	screen := gsnake.NewScreen(writer, rows, cols)
	leaderboard := gsnake.NewLeaderboard()
	snake := gsnake.NewSnake(screen)
	fruit := gsnake.NewFruit(rows, cols)
	return gsnake.NewGame(eventsPoller, screen, leaderboard, fruit, snake)
}

func displayHelp() {
	fmt.Println("Usage: gsnake [-m mode] [-p port]")
	fmt.Println("\nOptions:")
	fmt.Println("  -m, --mode\n\t(Optional) Expected values are 'local' or 'ssh'. If -m flag is set to 'ssh' then it will host a gsnake ssh server. Defaults to 'local'.")
	fmt.Println("  -p, --port\n\t(Optional) Port number. Only used in 'ssh' mode. Defaults to 5555.")
	fmt.Println("  -l, --log\n\t(Optional) Enables logging persistence. Defaults to false.")
	fmt.Println("  -h, --help\n\tDisplay this help text.")
	os.Exit(0)
}
