package main

import (
	"io"

	"turutupa/gsnake/events"
	"turutupa/gsnake/game"
	"turutupa/gsnake/ssh"
)

func main() {
	port := 2222
	sshServer := ssh.NewSshServer(port)
	sshServer.Run(snakeApp)
}

func snakeApp(writer io.Writer, eventsPoller events.EventPoller) ssh.SshApp {
	rows := 20
	cols := 55
	screen := gsnake.NewScreen(writer, rows, cols)
	scoreboard, _ := gsnake.NewScoreboard()
	snake := gsnake.NewSnake(screen)
	fruit := gsnake.NewFruit(rows, cols)
	game := gsnake.NewGame(eventsPoller, screen, scoreboard, fruit, snake)
	return game
}
