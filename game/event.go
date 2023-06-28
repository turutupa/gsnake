package gsnake

import (
	"errors"
	"turutupa/gsnake/events"
	"turutupa/gsnake/log"
)

// Int value of arrow keys
const (
	ARROW_UP    int = 65
	ARROW_DOWN  int = 66
	ARROW_RIGHT int = 67
	ARROW_LEFT  int = 68
)

type EventBus struct {
	state       *StateBus
	eventPoller events.EventPoller
	strategies  map[AppState]func(rune)
}

func NewEventBus(state *StateBus, eventPoller events.EventPoller) *EventBus {
	return &EventBus{state, eventPoller, make(map[AppState]func(rune))}
}

func (e *EventBus) subscribe(state AppState, strategy func(rune)) error {
	if _, exists := e.strategies[state]; exists {
		return errors.New("Strategy already exists for state " + string(state))
	}
	e.strategies[state] = strategy
	return nil
}

func (eb *EventBus) Run() {
	for {
		event, err := eb.eventPoller.Poll()
		if err != nil {
			return
		}
		e := rune(event)
		if strategy, exists := eb.strategies[eb.state.get()]; !exists {
			log.Error("Event Bus", errors.New("No active strategy for state "+string(eb.state.get())))
		} else {
			strategy(e)
		}
	}
}

func (eb *EventBus) Stop() {
	eb.eventPoller.Close()
}

func isBackspaceOrDelete(r rune) bool {
	return r == '\b' || r == '\u007F'
}

func isUserAcceptedChar(r rune) bool {
	return byte(r) >= 33 && byte(r) <= 126
}

// Accepted keys for up/down/left and right are
// - wasd
// - hjkl
// - arrow keys
func isUp(event rune) bool {
	return event == 'w' || int(event) == ARROW_UP || event == 'k'
}

func isDown(event rune) bool {
	return event == 's' || int(event) == ARROW_DOWN || event == 'j'
}

func isLeft(event rune) bool {
	return event == 'a' || int(event) == ARROW_LEFT || event == 'h'
}

func isRight(event rune) bool {
	return event == 'd' || int(event) == ARROW_RIGHT || event == 'l'
}

// enter keys are
// - enter
// - spacebar
// - \r which I'm not sure which key is that tbh
func isEnterKey(input rune) bool {
	in := byte(input)
	enterKeys := [2]byte{'\n', '\r'} // Byte representations of "enter" keys
	for _, key := range enterKeys {
		if in == key {
			return true
		}
	}
	if int(in) == 32 { // space
		return true
	}
	return false
}
