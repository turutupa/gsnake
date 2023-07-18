package gsnake

import (
	"fmt"
	"sync"
	"time"
	"turutupa/gsnake/log"
)

const WAITING_PLAYERS_IN_S = 15
const MAX_ROOM_SIZE int = 4

type Room struct {
	id                 string
	players            []*Player
	started            bool
	finished           bool
	game               *MultiGame
	lock               *sync.Mutex
	notifyPlayerJoined chan bool
	notifyGameStart    chan bool
	notifyGameEnd      chan bool
}

func NewRoom(game *MultiGame) *Room {
	return &Room{
		id:                 generateUUID(),
		players:            []*Player{},
		started:            false,
		finished:           false,
		game:               game,
		lock:               &sync.Mutex{},
		notifyPlayerJoined: make(chan bool),
		notifyGameStart:    make(chan bool),
		notifyGameEnd:      make(chan bool),
	}
}

func (r *Room) Close() {
	close(r.notifyPlayerJoined)
	close(r.notifyGameStart)
	close(r.notifyGameEnd)
}

func (r *Room) Run() {
	waitingForPlayers := true
	for waitingForPlayers {
		select {
		case <-r.notifyPlayerJoined:
			r.game.DefaultLayout(r.started)
			info := fmt.Sprintf("New player joined. Players in room %s:", r.id)
			for _, player := range r.players {
				info = info + "\n\t* " + player.name
			}
			log.Info(info)
		case <-r.notifyGameStart:
			if len(r.players) > 1 {
				waitingForPlayers = false
			}
		case <-time.After(WAITING_PLAYERS_IN_S * time.Second):
			if len(r.players) > 1 {
				waitingForPlayers = false
			}
		}
	}
	log.Info("Starting game for room %s", r.id)
	r.started = true
	for {
		r.game.DefaultLayout(r.started)
		r.game.Countdown("Game starts in...")
		r.game.DefaultLayout(r.started)
		r.game.Run()
		if len(r.players) == 1 {
			for _, p := range r.players {
				p.screen.RenderWarning(r.game.board, "Sorry. All players left!")
				r.game.Countdown("Exiting in...")
			}
			break
		}
		if r.game.HasWinner() {
			r.game.Leaderboard()
			time.Sleep(5 * time.Second)
			break
		} else {
			time.Sleep(time.Second)
			r.game.Leaderboard()
			r.game.RestartRound()
			time.Sleep(3 * time.Second)
		}
	}
	for i := 0; i < len(r.players); i++ {
		r.notifyGameEnd <- true
	}
	r.finished = true
}

func (r *Room) AddPlayer(player *Player) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.game.players) >= MAX_ROOM_SIZE {
		return false
	}
	r.game.AddPlayer(player)
	r.players = append(r.players, player)
	r.notifyPlayerJoined <- true
	if len(r.game.players) == MAX_ROOM_SIZE { // n should be the maximum number of players allowed in a room
		r.notifyGameStart <- true
	}
	return true
}

func (r *Room) OnExit(player *Player) {
	r.lock.Lock()
	defer r.lock.Unlock()

	players := []*Player{}
	for _, p := range r.players {
		if p != player {
			players = append(players, p)
		}
	}
	r.players = players
	r.game.players = players
}
