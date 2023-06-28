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

type State struct {
	state       AppState
	mode        GameMode
	diff        string
	subscribers []func(AppState)
}

func NewStateBus() *State {
	return &State{
		state:       MAIN_MENU, // defaults to MAIN MENU
		mode:        SINGLE,    // defaults to single player
		diff:        NORMAL,    // defaults to NORMAL
		subscribers: []func(AppState){},
	}
}

func (sb *State) set(appState AppState) *State {
	sb.state = appState
	return sb
}

func (sb *State) gameMode(gameMode GameMode) *State {
	sb.mode = gameMode
	return sb
}

func (sb *State) difficulty(difficulty string) *State {
	sb.diff = difficulty
	return sb
}

func (sb *State) get() AppState {
	return sb.state
}
