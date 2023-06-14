package gsnake

import (
	"fmt"
	"strconv"
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

func (s *Screen) clear() {
	for i := 1; i < s.rows-1; i++ {
		for j := 1; j < s.cols-1; j++ {
			s.matrix[i][j] = ' '
		}
	}
}

/*
* updates the snake on the matrix
 */
func (s *Screen) update(fruit *Fruit, node *Node, score int) {
	// render scoreboard
	scoreboard := []rune{'[', ' ', 'S', 'C', 'O', 'R', 'E', ' ', ']', ' '}
	j := 4
	for _, r := range scoreboard {
		s.matrix[0][j] = r
		j++
	}
	j += 2
	for _, r := range "[ " + strconv.Itoa(score) + " ]" {
		s.matrix[0][j] = r
		j++
	}

	// render fruit
	fruits := []rune{'*', '@', '#', '¶', 'ø'}
	s.matrix[fruit.x][fruit.y] = rune(fruits[randInt(len(fruits)-1)])
	// s.matrix[fruit.x][fruit.y] = rune('@')

	// render snake
	s.matrix[node.x][node.y] = rune(node.pointing)
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
		node = node.next
	}
}

func (s *Screen) render() {
	for i := 0; i < s.rows; i++ {
		var row string
		for j := 0; j < s.cols; j++ {
			row = row + string(s.matrix[i][j])
		}
		fmt.Println(row)
	}
}
