package events

type EventPoller interface {
	Poll() (byte, error)
	Close()
}
