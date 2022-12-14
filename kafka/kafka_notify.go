package kafka_notify

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"sync"
)

type kafkaNotify struct {
	w *kafka.Writer
	r *kafka.Reader
	l *zap.Logger
}

var rwMutex sync.RWMutex
var events = make(map[chan interface{}]map[string]bool)

func NewKafkaNotify(address, topic, groupID string, l *zap.Logger) *kafkaNotify {
	return &kafkaNotify{
		w: &kafka.Writer{
			Addr:     kafka.TCP(address),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{address},
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		l: l,
	}
}

// Post a notification (arbitrary data) to the specified event.
func (k kafkaNotify) Post(event string, data interface{}) error {
	return k.w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(event),
		Value: []byte(fmt.Sprintf("%v", data)),
	})
}

// Start starts listening for an event on a given channel.
func (k kafkaNotify) Start(event string, outputChan chan interface{}) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	// Avoid panic.
	if events[outputChan] == nil {
		events[outputChan] = make(map[string]bool, 0)
	}

	// Do not trigger another goroutine if channel is already being used.
	if len(events[outputChan]) == 0 {
		go k.Listen(outputChan)
	}
	events[outputChan][event] = true
	k.l.Debug("New listener registered", zap.Any("event", event))
}

// Stop stops listening for an event on a given channel.
func (k kafkaNotify) Stop(event string, outputChan chan interface{}) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	if _, ok := events[outputChan][event]; ok {
		events[outputChan][event] = false
	}
	k.l.Debug("Listener unregistered", zap.Any("event", event))
	return nil
}

func (k kafkaNotify) Listen(outputChan chan interface{}) error {
	for {
		m, err := k.r.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		if shouldListen, ok := events[outputChan][string(m.Key)]; ok && shouldListen {
			k.l.Info("Received message",
				zap.Any("key", string(m.Key)),
				zap.Any("val", string(m.Value)),
			)
			outputChan <- m.Value
		}
	}
}

func (k kafkaNotify) Close() {
	err := k.w.Close()
	if err != nil {
		k.l.Error("failed to close kafkaNotify writer", zap.Error(err))
	}

	err = k.r.Close()
	if err != nil {
		k.l.Error("failed to close kafkaNotify reader", zap.Error(err))
	}
	k.l.Info("Notifier shutdown successfully")
}
