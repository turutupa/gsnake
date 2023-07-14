package gsnake

const MAX_PLAYER_NAME_LEN = 6

type Player struct {
	id            string
	name          string
	score         int
	screen        *Screen
	snake         *Snake
	isAlive       bool
	keypressCh    chan bool
	nameSubmitted chan bool
}

func NewPlayer(name string) *Player {
	return &Player{"", name, 0, nil, nil, true, make(chan bool), make(chan bool)}
}

func (p *Player) WithSnake(row int, col int) *Player {
	p.snake = NewSnake(row, col)
	return p
}

func (p *Player) WithScore(score int) *Player {
	p.score = score
	return p
}

// only used for online players
func (p *Player) WithUUID() *Player {
	p.id = generateUUID()
	return p
}

func (p *Player) WithScreen(screen *Screen) *Player {
	p.screen = screen
	return p
}

func (p *Player) SetName() {
	label := "Set player name"
	p.screen.InputBox(label, p.name)
	for {
		select {
		case <-p.keypressCh:
			p.screen.InputBox(label, p.name)
		case <-p.nameSubmitted:
			close(p.keypressCh)
			close(p.nameSubmitted)
			return
		}
	}
}

// so far only used to prompt for user input
func (p *Player) SubmitNameStrategy(event rune) {
	name, done := HandleUserInputForm(p.name, event)
	if done && len(name) > 0 {
		p.nameSubmitted <- true
		return
	}
	p.name = name
	p.keypressCh <- true
}

// so far only used to prompt for user input
func (p *Player) MultiGameStrategy(event rune) {
	if event == 'q' {

		return
	}
	snakeStrategy(p.snake, event)
}
