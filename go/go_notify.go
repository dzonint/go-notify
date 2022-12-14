package go_notify

import (
	"go.uber.org/zap"
	"sync"
)

var rwMutex sync.RWMutex
var events = make(map[string][]chan interface{})

type goNotify struct {
	l *zap.Logger
}

func NewGoNotify(l *zap.Logger) *goNotify {
	return &goNotify{l: l}
}

// Post a notification (arbitrary data) to the specified event.
func (g goNotify) Post(event string, data interface{}) error {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	if outChans, ok := events[event]; ok {
		for _, outputChan := range outChans {
			g.l.Debug("Sending data to listeners",
				zap.Any("event", event),
				zap.Any("data", data))
			outputChan <- data
		}
	}
	return nil
}

// Start observing the specified event via provided output channel.
func (g goNotify) Start(event string, outputChan chan interface{}) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	// Avoid panic.
	if events[event] == nil {
		events[event] = make([]chan interface{}, 0)
	}

	events[event] = append(events[event], outputChan)
	g.l.Debug("New listener registered", zap.Any("event", event))
}

// Stop observing the specified event on the provided output channel.
func (g goNotify) Stop(event string, outputChan chan interface{}) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	if outChans, ok := events[event]; ok {
		newChannelArray := make([]chan interface{}, 0)
		for _, ch := range outChans {
			if ch != outputChan {
				newChannelArray = append(newChannelArray, ch)
			}
		}
		events[event] = newChannelArray
	}
	g.l.Debug("Listener unregistered", zap.Any("event", event))
	return nil
}

func (g goNotify) Close() {
	for _, channels := range events {
		for _, channel := range channels {
			close(channel)
		}
	}
	g.l.Info("Notifier shutdown successfully")
}
