package gsnake

import (
	"fmt"
	"time"
)

type Speed int

const (
	Easy     Speed = 100
	Normal         = 60
	Hard           = 40
	Insanity       = 20
)

type Game struct {
	screen *Screen
	*Term
	snake   *Snake
	fruit   *Fruit
	speed   int
	running bool
	score   int
}

func NewGame(
	term *Term,
	screen *Screen,
	fruit *Fruit,
	snake *Snake,
	speed Speed,
) *Game {
	game := &Game{
		screen:  screen,
		Term:    term,
		snake:   snake,
		fruit:   fruit,
		speed:   int(speed),
		running: true,
	}
	term.OnExit = func() {
		game.running = false
	}
	return game
}

func (g *Game) Run() {
	go g.executeUserInput()
	g.Term.clearTerminal()
	g.screen.init()
	for g.running {
		g.screen.clear(g.fruit, g.snake.head, g.snake.tail, g.score)
		g.snake.move()
		if g.ateFruit() {
			g.score += 10
			g.snake.append()
			g.snake.append()
			g.fruit.new()
		}
		g.screen.update(g.fruit, g.snake.head, g.score)
		g.screen.render(g.fruit, g.snake.head, g.score)
		if g.intersects() {
			time.Sleep(1 * time.Second)
			g.Term.clearTerminal()
			g.GameOver()
			return
		}
		// adding some extra time when going vertical because it feels faster
		// than when going horizontally due to font width/height
		duration := time.Duration(g.speed) * time.Millisecond
		if g.snake.head.pointing == UP || g.snake.head.pointing == DOWN {
			time.Sleep(duration + (duration * 3 / 4))
		} else {
			time.Sleep(duration)
		}
	}
}

func (g *Game) ateFruit() bool {
	x := g.snake.head.x
	y := g.snake.head.y
	return x == g.fruit.x && y == g.fruit.y
}

func (g *Game) intersects() bool {
	head := g.snake.head
	x := head.x
	y := head.y
	if x == 0 || x == g.screen.rows-1 || y == 0 || y == g.screen.cols-1 {
		return true
	}

	node := head.next
	for node != nil && node.validated {
		if x == node.x && y == node.y {
			return true
		}
		node = node.next
	}
	return false
}

func (g *Game) executeUserInput() {
	for {
		event := g.PollEvents()
		pointing := g.snake.head.pointing
		if event == 'w' {
			if pointing != DOWN {
				g.snake.head.pointing = UP
			}
		} else if event == 's' {
			if pointing != UP {
				g.snake.head.pointing = DOWN
			}
		} else if event == 'a' {
			if pointing != RIGHT {
				g.snake.head.pointing = LEFT
			}
		} else if event == 'd' {
			if pointing != LEFT {
				g.snake.head.pointing = RIGHT
			}
		}
	}
}

func (g *Game) GameOver() {
	gameOver := `
       $$$$$$\   $$$$$$\  $$\      $$\ $$$$$$$$\       
      $$  __$$\ $$  __$$\ $$$\    $$$ |$$  _____|      
      $$ /  \__|$$ /  $$ |$$$$\  $$$$ |$$ |            
      $$ |$$$$\ $$$$$$$$ |$$\$$\$$ $$ |$$$$$\          
      $$ |\_$$ |$$  __$$ |$$ \$$$  $$ |$$  __|         
      $$ |  $$ |$$ |  $$ |$$ |\$  /$$ |$$ |            
      \$$$$$$  |$$ |  $$ |$$ | \_/ $$ |$$$$$$$$\        
       \______/ \__|  \__|\__|     \__|\________|       

       $$$$$$\  $$\    $$\ $$$$$$$$\ $$$$$$$\  
      $$  __$$\ $$ |   $$ |$$  _____|$$  __$$\ 
      $$ /  $$ |$$ |   $$ |$$ |      $$ |  $$ |
      $$ |  $$ |\$$\  $$  |$$$$$\    $$$$$$$  |
      $$ |  $$ | \$$\$$  / $$  __|   $$  __$$< 
      $$ |  $$ |  \$$$  /  $$ |      $$ |  $$ |
       $$$$$$  |   \$  /   $$$$$$$$\ $$ |  $$ |
       \______/     \_/    \________|\__|  \__|
`
	fmt.Println(gameOver)
}
