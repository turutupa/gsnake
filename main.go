package main

import (
	"io"
	"turutupa/gsnake/game"
	"turutupa/gsnake/ssh"
)

func main() {
	port := 2222
	sshServer := ssh.NewSshServer(port, snakeHandler)
	sshServer.Run()
}

func snakeHandler(writer io.Writer) ssh.Runnable {
	// inject ssh/local cli writer to screen
	rows := 20
	cols := 55
	term := gsnake.NewTerm()
	screen := gsnake.NewScreen(writer, rows, cols)
	scoreboard, _ := gsnake.NewScoreboard()
	snake := gsnake.NewSnake(screen)
	fruit := gsnake.NewFruit(rows, cols)
	game := gsnake.NewGame(term, screen, scoreboard, fruit, snake)
	return game
}
