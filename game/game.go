package gsnake

// Snake available speeds
const (
	ONLINE_SPEED   = 80
	EASY_SPEED     = 50
	NORMAL_SPEED   = 40
	HARD_SPEED     = 30
	INSANITY_SPEED = 20
)

const ROWS_SOLO = 20
const COLS_SOLO = 50

const ROWS_MULTI = 30
const COLS_MULTI = 80

type Game interface {
	Run()
	Leaderboard()
	Restart()
	SetDifficulty(string)
	Stop()
	Strategy(event rune)
}
