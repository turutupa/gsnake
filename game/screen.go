package gsnake

import (
	"fmt"
	"strconv"
	"strings"
)

const VERTICAL rune = '│'
const TOP_LEFT rune = '╭'
const TOP_RIGHT rune = '╮'
const BOTTOM_LEFT rune = '╰'
const BOTTOM_RIGHT rune = '╯'
const HORIZONTAL rune = '─'

type Screen struct {
	rows          int
	cols          int
	matrix        [][]rune
	isFirstRender bool
}

func NewScreen(rows int, cols int) *Screen {
	var matrix [][]rune
	for i := 0; i < rows; i++ {
		row := make([]rune, cols)
		matrix = append(matrix, row)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			var cell rune
			if i == 0 && j == 0 { // top left
				cell = TOP_LEFT
			} else if i == 0 && j == cols-1 { // top right
				cell = TOP_RIGHT
			} else if i == rows-1 && j == 0 { // bottom left
				cell = BOTTOM_LEFT
			} else if i == rows-1 && j == cols-1 { // bottom right
				cell = BOTTOM_RIGHT
			} else if i == 0 || i == rows-1 { // top && bottom rows
				cell = HORIZONTAL
			} else if j == 0 || j == cols-1 { // left && right cols
				cell = VERTICAL
			} else {
				cell = ' '
			}
			matrix[i][j] = cell
		}
	}
	return &Screen{rows, cols, matrix, true}
}

func (s *Screen) restart() {
	newScreen := NewScreen(s.rows, s.cols)
	s.matrix = newScreen.matrix
	s.isFirstRender = newScreen.isFirstRender
}

func (s *Screen) updateScoreboard(score int) {
	padded_score := strconv.Itoa(score)
	padded_score = strings.Repeat("0", 5-len(padded_score)) + padded_score
	scoreboard := "[ SCORE ]──[ " + padded_score + " ]"
	quitMsg := "[ 'q' to Quit ]"
	padding := 4
	scoreboard = scoreboard + strings.Repeat("─", s.cols-len(scoreboard)-len(quitMsg)-padding) + quitMsg
	j := padding
	for _, r := range scoreboard {
		s.matrix[0][j] = r
		j++
	}
}

// print borders
func (s *Screen) init() {
	s.updateScoreboard(0)
	for i := 0; i < s.cols; i++ {
		s.print(0, i, s.matrix[0][i])
		s.print(s.rows-1, i, s.matrix[s.rows-1][i])
	}
	for i := 0; i < s.rows; i++ {
		s.print(i, 0, s.matrix[i][0])
		s.print(i, s.cols-1, s.matrix[i][s.cols-1])
	}
}

func (s *Screen) clear(fruit *Fruit, head *Node, tail *Node, score int) {
	s.print(head.x, head.y, ' ')
	s.print(tail.x, tail.y, ' ')
}

func (s *Screen) renderMainMenu(selected int) {
	s.init()
	title := "SELECT GAME MODE"
	startLine := s.cols/2 - len(title)/2 - 1
	row := s.rows / 5
	gameModes := []string{"EASY", "NORMAL", "HARD", "INSANITY", "EXIT"}
	s.printBold(row, startLine, title)
	row += 2
	optionIndex := 0
	for i, game := range gameModes {
		paddingRight := 8
		if optionIndex == selected {
			selectedIndicatorLeft := "> "
			selectedIndicatorRight := " <"
			gameFmt := strings.Repeat(" ", len(title)-len(selectedIndicatorLeft)-(len(game)/2)-paddingRight)
			gameFmt = gameFmt + selectedIndicatorLeft + game + selectedIndicatorRight
			gameFmt = gameFmt + strings.Repeat(" ", paddingRight-len(selectedIndicatorRight))
			s.printBold(row, startLine, gameFmt)
		} else {
			gameFmt := strings.Repeat(" ", len(title)-(len(game)/2)-paddingRight)
			gameFmt = gameFmt + game
			gameFmt = gameFmt + strings.Repeat(" ", paddingRight)
			for i, r := range gameFmt {
				s.print(row, startLine+i, r)
			}
		}
		if i == len(gameModes)-2 {
			row += 2
		} else {
			row++
		}
		optionIndex++
	}
	s.finishPrint()
}

func (s *Screen) renderSnake(fruit *Fruit, head *Node, tail *Node, score int) {
	// - render fruit -
	if s.matrix[fruit.x][fruit.y] == '@' {
		s.print(fruit.x, fruit.y, '@')
	}

	// - render snake -
	// the first time we render the entire snake
	// the rest of the time only the parts of the snake
	// that requires re-rendering
	if s.isFirstRender {
		node := head
		for node != nil && node.validated {
			s.print(node.x, node.y, node.render)
			node = node.next
		}
		s.isFirstRender = false
	} else {
		// - render snake -
		s.print(head.x, head.y, head.render)
		s.print(head.next.x, head.next.y, head.next.render)
		node := tail
		for !node.validated {
			node = node.prev
		}
		s.print(node.x, node.y, node.render)
	}

	// - top border -
	// we only want to re-render the score so
	// calcualte where the score is positioned,
	// which is between the second set of brackets
	bracket_counter := 0
	for j := 0; j < s.cols; j++ {
		if s.matrix[0][j] == '[' || s.matrix[0][j] == ']' {
			bracket_counter = bracket_counter + 1
			continue
		}
		if bracket_counter == 3 {
			s.print(0, j, s.matrix[0][j])
		}
		if bracket_counter >= 4 {
			break
		}
	}
	s.finishPrint()
}

func (s *Screen) renderScoreboard(scores []int) {
	title := "| " + "TOP SCORES" + " |"
	for len(scores) < 5 {
		scores = append(scores, 0)
	}
	scores = scores[:5]
	marginLeft := len(title)/2 - 1
	startLine := s.cols/2 - marginLeft
	row := s.rows/3 - 1
	for i := 0; i < len(title); i++ {
		position := startLine + i
		if i == 0 {
			s.print(row, position, TOP_LEFT)
		} else if i == len(title)-1 {
			s.print(row, position, TOP_RIGHT)
		} else {
			s.print(row, position, HORIZONTAL)
		}
	}
	row++
	for i, r := range title {
		s.print(row, s.cols/2-marginLeft+i, r)
	}
	row++

	for _, score := range scores {
		rightPadding := 2
		scoreStr := strconv.Itoa(score)
		scoreFmt := strings.Repeat(" ", len(title)-len(scoreStr)-rightPadding)
		scoreFmt = scoreFmt + scoreStr
		scoreFmt = scoreFmt + strings.Repeat(" ", rightPadding) // add padding to the right if needed

		for j, r := range scoreFmt {
			position := startLine + j
			if j == 0 || j == len(title)-1 {
				s.print(row, position, '|')
			} else {
				s.print(row, position, r)
			}
		}
		row++
	}

	for i := range title {
		position := s.cols/2 - marginLeft + i
		if i == 0 {
			s.print(row, position, BOTTOM_LEFT)
		} else if i == len(title)-1 {
			s.print(row, position, BOTTOM_RIGHT)
		} else {
			s.print(row, position, HORIZONTAL)
		}
	}
	s.finishPrint()
}

/*
* updates the snake on the matrix
 */
func (s *Screen) update(fruit *Fruit, node *Node, score int) {
	// render scoreboard
	s.updateScoreboard(score)

	// render fruit
	s.matrix[fruit.x][fruit.y] = rune('@')

	// render snake
	s.matrix[node.x][node.y] = rune(node.pointing)
	node.render = rune(node.pointing)
	node = node.next
	for node != nil && node.validated {
		var cell rune
		if node.next != nil {
			if node.pointing == node.next.pointing {
				if node.pointing == UP || node.pointing == DOWN {
					cell = VERTICAL
				} else {
					cell = HORIZONTAL
				}
			} else if node.pointing == UP {
				if node.next.pointing == LEFT {
					cell = BOTTOM_LEFT
				} else {
					cell = BOTTOM_RIGHT
				}
			} else if node.pointing == DOWN {
				if node.next.pointing == LEFT {
					cell = TOP_LEFT
				} else {
					cell = TOP_RIGHT
				}
			} else if node.pointing == LEFT {
				if node.next.pointing == UP {
					cell = TOP_RIGHT
				} else {
					cell = BOTTOM_RIGHT
				}
			} else if node.pointing == RIGHT {
				if node.next.pointing == UP {
					cell = TOP_LEFT
				} else {
					cell = BOTTOM_LEFT
				}
			}
			s.matrix[node.x][node.y] = cell
		} else {
			if node.pointing == UP || node.pointing == DOWN {
				cell = VERTICAL
			} else {
				cell = HORIZONTAL
			}
		}

		x := node.x
		y := node.y

		if x < 1 || x >= s.rows-1 || y < 1 || y >= s.cols-1 {
			continue
		}

		s.matrix[x][y] = cell
		node.render = cell
		node = node.next
	}
}

func (s *Screen) print(row, col int, r rune) {
	fmt.Printf("\033[%d;%dH%c", row+1, col+1, r)
}

func (s *Screen) printBold(row, col int, r string) {
	fmt.Printf("\033[%d;%dH\033[1m%s\033[0m", row+1, col+1, r)
}

func (s *Screen) finishPrint() {
	fmt.Printf("\033[%d;%dH", s.rows+2, 0)
}

func (s *Screen) GameOver() {
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


 

                PRESS ENTER TO CONTINUE
`
	fmt.Println(gameOver)
}

func (s *Screen) printLogo() {
	logo := `
                                              $$\                 
                                              $$ |                
       $$$$$$\   $$$$$$$\ $$$$$$$\   $$$$$$\  $$ |  $$\  $$$$$$\  
      $$  __$$\ $$  _____|$$  __$$\  \____$$\ $$ | $$  |$$  __$$\ 
      $$ /  $$ |\$$$$$$\  $$ |  $$ | $$$$$$$ |$$$$$$  / $$$$$$$$ |
      $$ |  $$ | \____$$\ $$ |  $$ |$$  __$$ |$$  _$$<  $$   ____|
      \$$$$$$$ |$$$$$$$  |$$ |  $$ |\$$$$$$$ |$$ | \$$\ \$$$$$$$\ 
       \____$$ |\_______/ \__|  \__| \_______|\__|  \__| \_______|
      $$\   $$ |                                                  
      \$$$$$$  |                                                  
       \______/
  `
	fmt.Println(logo)
}
