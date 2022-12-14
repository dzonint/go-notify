package kafka_notify

import (
	"github.com/dzonint/go-notify/kafka/mocks"
	"github.com/pkg/errors"
	"testing"
)

var mockKafkaNotify = mocks.NewKafkaNotify()

func TestKafkaNotify_PostError(t *testing.T) {
	mocks.Post = func(event string, data interface{}) error {
		return errors.New("an error occurred")
	}
	err := mockKafkaNotify.Post("test", []string{"test-data"})
	if err == nil {
		t.Error("expected:", nil, "got:", err)
	}
}

func TestKafkaNotify_Post(t *testing.T) {
	mocks.Post = func(event string, data interface{}) error {
		return nil
	}
	err := mockKafkaNotify.Post("test", []string{"test-data"})
	if err != nil {
		t.Error("expected:", err, "got:", nil)
	}
}

func TestKafkaNotify_StartAndStop(t *testing.T) {
	events := make(map[chan interface{}]map[string]bool)
	outputChan := make(chan interface{})

	// Test start - channel is present and set to listen.
	mocks.Start = func(event string, outputChan chan interface{}) {
		// Avoid panic.
		if events[outputChan] == nil {
			events[outputChan] = make(map[string]bool, 0)
		}
		events[outputChan][event] = true
	}
	mocks.Start("test", outputChan)
	if val, ok := events[outputChan]["test"]; !val || !ok {
		t.Error("expected:", true, true, "- got:", val, ok)
	}

	// Test stop - channel is present but set to not listen.
	mocks.Stop = func(event string, outputChan chan interface{}) error {
		if _, ok := events[outputChan][event]; ok {
			events[outputChan][event] = false
		}
		return nil
	}

	mocks.Stop("test", outputChan)
	if val, ok := events[outputChan]["test"]; val || !ok {
		t.Error("expected:", false, false, "- got:", val, ok)
	}
}

func TestKafkaNotify_ListenError(t *testing.T) {
	outputChan := make(chan interface{})
	mocks.Listen = func(outputChan chan interface{}) error {
		return errors.New("an error occurred")
	}
	err := mockKafkaNotify.Listen(outputChan)
	if err == nil {
		t.Error("expected:", err, "got:", nil)
	}
}

func TestKafkaNotify_Listen(t *testing.T) {
	outputChan := make(chan interface{})
	mocks.Listen = func(outputChan chan interface{}) error {
		return nil
	}
	err := mockKafkaNotify.Listen(outputChan)
	if err != nil {
		t.Error("expected:", err, "got:", nil)
	}
}
