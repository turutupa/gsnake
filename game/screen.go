package gsnake

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const VERTICAL rune = '│'
const TOP_LEFT rune = '╭'
const TOP_RIGHT rune = '╮'
const BOTTOM_LEFT rune = '╰'
const BOTTOM_RIGHT rune = '╯'
const HORIZONTAL rune = '─'

const MAX_HIGH_SCORE_DIGITS = 4

type Screen struct {
	writer        io.Writer
	isFirstRender bool
	termRows      int
	termCols      int
}

func NewScreen(writer io.Writer) *Screen {
	return &Screen{writer, true, 30, 50}
}

func (s *Screen) Restart() {
	s.isFirstRender = true
}

func (s *Screen) Clear() {
	s.writer.Write([]byte("\033[H\033[2J"))
}

func (s *Screen) SetSize(rows int, cols int) {
	s.termRows = rows
	s.termCols = cols
}

func (s *Screen) PromptPlayerName() {}

func (s *Screen) RenderMainMenu(termSize TermSize, selected int) {
	title := "SELECT GAME MODE"
	startLine := termSize.cols/2 - len(title)/2 - 1
	row := termSize.rows / 3
	s.printBold(row, startLine, title)
	row += 2
	optionIndex := 0
	for _, game := range MENU_OPTIONS {
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
		if game == INSANITY {
			row += 2
		} else {
			row++
		}
		optionIndex++
	}
}

func (s *Screen) RenderBoard(board *Board) {
	offsetRow, offsetCol := s.offsets(board)
	for i := 0; i < board.cols; i++ {
		s.print(offsetRow+0, offsetCol+i, board.matrix[0][i])
		s.print(offsetRow+board.rows-1, offsetCol+i, board.matrix[board.rows-1][i])
	}
	for i := 0; i < board.rows; i++ {
		s.print(offsetRow+i, offsetCol+0, board.matrix[i][0])
		s.print(offsetRow+i, offsetCol+board.cols-1, board.matrix[i][board.cols-1])
	}
}

func (s *Screen) RenderCountdown(n int) {
	switch n {
	case 1:
		s.writer.Write([]byte(ONE))
	case 2:
		s.writer.Write([]byte(TWO))
	case 3:
		s.writer.Write([]byte(THREE))
	case 4:
		s.writer.Write([]byte(FOUR))
	case 5:
		s.writer.Write([]byte(FIVE))
	}
}

func (s *Screen) RenderWarning(rows, cols int, txt string) {
	row := rows / 2
	col := (cols / 2) - len(txt)/2
	s.printBold(row, col, txt)
}

func (s *Screen) RenderFruit(board *Board, fruit *Fruit) {
	offsetRow, offsetCol := s.offsets(board)
	s.print(offsetRow+fruit.x, offsetCol+fruit.y, '@')
}

// used in solo game
func (s *Screen) RenderScore(board *Board, score int) {
	// - top border -
	// we only want to re-render the score so
	// calcualte where the score is positioned,
	// which is between the second set of brackets
	offsetRow, offsetCol := s.offsets(board)
	bracket_counter := 0
	for j := 0; j < board.cols; j++ {
		if board.matrix[0][j] == '[' || board.matrix[0][j] == ']' {
			bracket_counter = bracket_counter + 1
			continue
		}
		if bracket_counter == 3 {
			s.print(offsetRow, offsetCol+j, board.matrix[0][j])
		}
		if bracket_counter >= 4 {
			break
		}
	}
}

func (s *Screen) Remove(board *Board, head *Node) {
	s.print(s.offsetRow(board)+head.x, s.offsetCol(board)+head.y, ' ')
	last := head
	for head != nil && head.validated {
		last = head
		head = head.next
	}
	s.print(s.offsetRow(board)+last.x, s.offsetCol(board)+last.y, ' ')
}

func (s *Screen) RenderSnake(board *Board, head *Node) {
	offsetRow, offsetCol := s.offsets(board)
	// room for optimization
	s.print(offsetRow+head.x, offsetCol+head.y, head.render)
	node := head.next
	for node != nil && node.validated {
		s.print(offsetRow+node.x, offsetCol+node.y, node.render)
		node = node.next
	}
}

func (s *Screen) RenderLeaderboard(board *Board, difficulty string, scores []*Player, player *Player) {
	title := "| " + "TOP SCORES" + " |"
	// how many high scores to render, if there is a new one,
	// that one is going to be the fifth top score
	var leaderboarSize int
	if player != nil {
		leaderboarSize = 4
	} else {
		leaderboarSize = 5
	}
	for len(scores) < leaderboarSize {
		// this is just to make sure we always render 5
		// rows of scores, even if player hasn't already
		// played 5 times
		scores = append(scores, NewPlayer(""))
	}
	scores = scores[:leaderboarSize]
	marginLeft := len(title)/2 - 1
	startLine := board.cols/2 - marginLeft
	row := board.rows/3 + 1
	offsetRow, offsetCol := s.offsets(board)

	// print top border
	for i := 0; i < len(title); i++ {
		position := startLine + i
		switch i {
		case 0:
			s.print(offsetRow+row, offsetCol+position, TOP_LEFT)
		case len(title) - 1:
			s.print(offsetRow+row, offsetCol+position, TOP_RIGHT)
		default:
			s.print(offsetRow+row, offsetCol+position, HORIZONTAL)

		}
	}
	row++

	// print title
	for i, r := range title {
		s.print(offsetRow+row, offsetCol+startLine+i, r)
	}
	row++

	// print difficulty
	padding := strings.Repeat(" ", (len(title)-len(difficulty)-2)/2)
	diff := "|" + padding + difficulty + padding + "|"
	for i, r := range diff {
		s.print(offsetRow+row, offsetCol+startLine+i, r)
	}
	row++

	// print empty line
	emptyLine := "|" + strings.Repeat(" ", len(title)-2) + "|"
	for i, r := range emptyLine {
		s.print(offsetRow+row, offsetCol+startLine+i, r)
	}
	row++

	// print scores
	renderedNewHighScore := false
	i := 0
	for i < len(scores) {
		score := scores[i]
		var scoreFmt string
		var sc int
		var pl string
		isHighScore := player != nil && !renderedNewHighScore && player.score > score.score
		if isHighScore && !renderedNewHighScore {
			sc = player.score
			pl = player.name
			renderedNewHighScore = true
		} else {
			sc = score.score
			pl = score.name
			i++
		}

		rightPadding := 2
		leftPadding := 2
		scoreStr := strconv.Itoa(sc)
		scoreFmt = strings.Repeat(" ", leftPadding) + pl
		if len(pl) < MAX_PLAYER_NAME_LEN && isHighScore {
			// format when submitting high score
			cursor := "_"
			scoreFmt = scoreFmt + cursor + strings.Repeat(" ", MAX_PLAYER_NAME_LEN-len(pl)-1)
		} else {
			// add padding spaces for player names
			scoreFmt = scoreFmt + strings.Repeat(" ", MAX_PLAYER_NAME_LEN-len(pl))
		}
		// add space between names and actual scores
		scoreFmt = scoreFmt + strings.Repeat(" ", 2)
		// add leading zeroes to score
		scoreFmt = scoreFmt + strings.Repeat("0", MAX_HIGH_SCORE_DIGITS-len(strconv.Itoa(sc)))
		scoreFmt = scoreFmt + scoreStr
		scoreFmt = scoreFmt + strings.Repeat(" ", rightPadding) // add padding to the right if needed

		for j, r := range scoreFmt {
			position := startLine + j
			if j == 0 || j == len(title)-1 {
				s.print(offsetRow+row, offsetCol+position, '|')
			} else {
				s.print(offsetRow+row, offsetCol+position, r)
			}
		}
		row++
	}

	// render border bottom
	for i := range title {
		position := startLine + i
		if i == 0 {
			s.print(offsetRow+row, offsetCol+position, BOTTOM_LEFT)
		} else if i == len(title)-1 {
			s.print(offsetRow+row, offsetCol+position, BOTTOM_RIGHT)
		} else {
			s.print(offsetRow+row, offsetCol+position, HORIZONTAL)
		}
	}
	row = row + 2
	var msg string
	if player != nil {
		msg = " PRESS ENTER TO SUBMIT NEW HIGH SCORE "
	} else {
		msg = " PRESS ENTER TO CONTINUE "
	}
	s.printBold(offsetRow+row, offsetCol+board.cols/2-len(msg)/2, msg)
}

func (s *Screen) HideCursor() {
	fmt.Fprint(s.writer, "\033[?25l")
}

func (s *Screen) ShowCursor() {
	fmt.Fprint(s.writer, "\033[?25h")
}

func (s *Screen) GameOver() {
	s.writer.Write([]byte(GAME_OVER))
}

func (s *Screen) PrintLogo() {
	s.writer.Write([]byte(LOGO))
}

func (s *Screen) print(row, col int, r rune) {
	fmt.Fprintf(s.writer, "\033[%d;%dH%c", row+1, col+1, r)
}

func (s *Screen) printBold(row, col int, r string) {
	fmt.Fprintf(s.writer, "\033[%d;%dH\033[1m%s\033[0m", row+1, col+1, r)
}

func (s *Screen) offsets(board *Board) (int, int) {
	return s.offsetRow(board), s.offsetCol(board)
}

func (s *Screen) offsetRow(board *Board) int {
	return (s.termRows - board.rows) / 2
}

func (s *Screen) offsetCol(board *Board) int {
	return (s.termCols - board.cols) / 2
}
