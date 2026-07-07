package eventbus

// EventBus is a minimal publish/subscribe surface.
type EventBus interface {
	Publish(topic string, msg any) error
	// Subscribe would return a channel or subscription handle in a full impl.
}

// NoopEventBus is a tiny implementation that does nothing.
type NoopEventBus struct{}

func (NoopEventBus) Publish(topic string, msg any) error { return nil }
