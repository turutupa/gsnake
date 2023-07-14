package gsnake

// Snake available speeds
const (
	ONLINE_SPEED   = 80
	EASY_SPEED     = 50
	NORMAL_SPEED   = 40
	HARD_SPEED     = 30
	INSANITY_SPEED = 20
)

type Game interface {
	Run()
	Leaderboard()
	Restart()
	SetDifficulty(string)
	Stop()
	Strategy(event rune)
}
