package main

import (
	"fmt"
	"os"
	"turutupa/gsnake/gsnake"
)

func main() {
	rows := 20
	cols := 55
	term := gsnake.NewTerm()
	screen := gsnake.NewScreen(rows, cols)
	scoreboard, _ := gsnake.NewScoreboard()
	snake := gsnake.NewSnake(screen)
	fruit := gsnake.NewFruit(rows, cols)
	game := gsnake.NewGame(term, screen, scoreboard, fruit, snake, getSnakeSpeed())
	game.Run()
}

func getSnakeSpeed() gsnake.Speed {
	args := os.Args
	var speed gsnake.Speed = gsnake.Normal // Default speed
	if len(args) >= 2 {
		arg := args[1]
		switch arg {
		case "--easy":
			speed = gsnake.Easy
		case "--normal":
			speed = gsnake.Normal
		case "--hard":
			speed = gsnake.Hard
		case "--insanity":
			speed = gsnake.Insanity
		default:
			fmt.Println("Invalid speed option. Using default speed: Normal")
		}
	}
	return speed
}
