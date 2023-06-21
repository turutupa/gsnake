package events

type EventPoller interface {
	Poll() byte
	Close()
}
