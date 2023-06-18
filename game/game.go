package gsnake

import (
	"math"
	"os"
	"time"
)

type Speed int

const (
	EXIT     Speed = 0
	EASY     Speed = 50
	NORMAL   Speed = 40
	HARD     Speed = 30
	INSANITY Speed = 20
)

var MENU_OPTIONS = []Speed{EASY, NORMAL, HARD, INSANITY, EXIT}

type State int

const (
	MAIN_MENU  State = 1
	PLAYING    State = 2
	SCOREBOARD State = 3
)

type Game struct {
	screen *Screen
	*Term
	scoreboard         *Scoreboard
	fruit              *Fruit
	snake              *Snake
	speed              int
	running            bool
	score              int
	state              State
	selectedMenuOption int
	selectChan         chan bool
}

func NewGame(
	term *Term,
	screen *Screen,
	scoreboard *Scoreboard,
	fruit *Fruit,
	snake *Snake,
	speed Speed,
) *Game {
	game := &Game{
		screen:             screen,
		Term:               term,
		scoreboard:         scoreboard,
		fruit:              fruit,
		snake:              snake,
		speed:              int(speed),
		running:            true,
		score:              0,
		state:              MAIN_MENU,
		selectedMenuOption: 2,
		selectChan:         make(chan bool),
	}
	term.OnExit = func() { game.restart() }
	return game
}

func (g *Game) Run() {
	go g.executeUserInput()
	for {
		for g.state == MAIN_MENU {
			g.mainMenu()
			<-g.selectChan // used for blocking
		}
		g.runGame()
	}
}

func (g *Game) restart() {
	g.state = MAIN_MENU
	g.screen.restart()
	g.snake.restart(g.screen)
}

func (g *Game) mainMenu() {
	g.Term.clearTerminal()
	g.screen.renderMainMenu(g.selectedMenuOption)
}

func (g *Game) runGame() {
	g.Term.clearTerminal()
	g.screen.init()
	for g.state == PLAYING {
		g.screen.clear(g.fruit, g.snake.head, g.snake.tail, g.score)
		g.snake.move()
		if g.ateFruit() {
			g.score += 10
			if g.speed == int(EASY) || g.speed == int(NORMAL) {
				g.snake.append()
				g.snake.append()
			} else {
				g.snake.append()
				g.snake.append()
				g.snake.append()
			}
			g.fruit.new()
		}
		g.screen.update(g.fruit, g.snake.head, g.score)
		g.screen.renderSnake(g.fruit, g.snake.head, g.snake.tail, g.score)
		if g.intersects() {
			scores, ok := g.scoreboard.update(g.score)
			time.Sleep(1 * time.Second)
			g.Term.clearTerminal()
			g.screen.GameOver()
			if ok {
				g.screen.renderScoreboard(scores)
			}
			g.state = SCOREBOARD
			<-g.selectChan
			g.restart()
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
		if g.state == MAIN_MENU {
			g.userActionMainMenu(event)
		} else if g.state == PLAYING {
			g.userActionSnake(event)
		} else if g.state == SCOREBOARD {
			g.userActionScoreboard(event)
		}
	}
}

func (g *Game) userActionMainMenu(event rune) {
	if event == 'w' || event == 'k' {
		g.selectedMenuOption = int(math.Max(float64(0), float64(g.selectedMenuOption-1)))
		g.selectChan <- true
	} else if event == 's' || event == 'j' {
		g.selectedMenuOption = int(math.Min(float64(len(MENU_OPTIONS)-1), float64(g.selectedMenuOption+1)))
		g.selectChan <- true
	} else if event == '\n' {
		g.speed = int(MENU_OPTIONS[g.selectedMenuOption])
		g.state = PLAYING
		if g.selectedMenuOption == len(MENU_OPTIONS)-1 {
			os.Exit(0)
		}
		g.selectChan <- true
	}
}

func (g *Game) userActionSnake(event rune) {
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

func (g *Game) userActionScoreboard(event rune) {
	if event == '\n' {
		g.selectChan <- true
	}
}
