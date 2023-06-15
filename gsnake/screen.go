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
	rows   int
	cols   int
	matrix [][]rune
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
	return &Screen{rows, cols, matrix}
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
	for node != nil {
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
	s.print(fruit.x, fruit.y, ' ')
	s.print(head.x, head.y, ' ')
	s.print(tail.x, tail.y, ' ')
}

func (s *Screen) render(fruit *Fruit, node *Node, score int) {
	// - render fruit -
	s.print(fruit.x, fruit.y, '@')

	// - render snake -
	for node != nil {
		if node.validated && node.x < s.rows && node.y < s.cols {
			s.print(node.x, node.y, node.render)
		}
		node = node.next
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

func (s *Screen) print(row, col int, r rune) {
	fmt.Printf("\033[%d;%dH%c", row+1, col+1, r)
}

func (s *Screen) printBold(row, col int, r rune) {
	fmt.Printf("\033[%d;%dH\033[1m%c\033[0m", row, col, r)
}

func (s *Screen) finishPrint() {
	fmt.Printf("\033[%d;%dH", s.rows+2, s.cols+2)
}
