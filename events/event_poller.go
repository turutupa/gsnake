package events

type EventPoller interface {
	PollEvents() byte
}
