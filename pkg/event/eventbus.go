package event

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type EventCallack func(payload any)

type Event struct {
	Type    string
	Payload any
}

type EventBus struct {
	logger   *zap.Logger
	mu       sync.RWMutex
	handlers map[string][]EventCallack
}

func NewEventBus(logger *zap.Logger) *EventBus {
	return &EventBus{
		logger:   logger,
		handlers: make(map[string][]EventCallack),
	}
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	start := time.Now()

	eb.logger.Info("Send Event to EventBus",
		zap.String("type", event.Type),
		zap.Any("payload", event.Payload),
	)

	if handlers, ok := eb.handlers[event.Type]; ok {
		for _, handler := range handlers {
			eb.logger.Info("Recieve Event from EventBus",
				zap.String("type", event.Type),
				zap.Any("payload", event.Payload),
				zap.Duration("duration", time.Since(start)),
			)
			handler(event.Payload)
		}
	}
}

func (eb *EventBus) On(eventType string, callback EventCallack) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.logger.Info("Subscribe to EventBus",
		zap.String("eventType", eventType),
	)

	eb.handlers[eventType] = append(eb.handlers[eventType], callback)
}
