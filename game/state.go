package gsnake

// App AppState
type AppState string
type GameMode string

const (
	SINGLE GameMode = "SINGLE"
	MULTI  GameMode = "MULTI"
)

const (
	MAIN_MENU AppState = "MAIN_MENU" // starts here by default
	IN_GAME   AppState = "IN_GAME"
	QUIT      AppState = "QUIT"
)

type StateBus struct {
	state       AppState
	mode        GameMode
	diff        string
	subscribers []func(AppState)
}

func NewStateBus() *StateBus {
	return &StateBus{
		state:       MAIN_MENU, // defaults to MAIN MENU
		mode:        SINGLE,    // defaults to single player
		diff:        NORMAL,    // defaults to NORMAL
		subscribers: []func(AppState){},
	}
}

func (sb *StateBus) set(appState AppState) *StateBus {
	sb.state = appState
	return sb
}

func (sb *StateBus) gameMode(gameMode GameMode) *StateBus {
	sb.mode = gameMode
	return sb
}

func (sb *StateBus) difficulty(difficulty string) *StateBus {
	sb.diff = difficulty
	return sb
}

func (sb *StateBus) get() AppState {
	return sb.state
}
