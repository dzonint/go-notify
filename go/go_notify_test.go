package go_notify

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

var gn = NewGoNotify(zap.L())

// TestGoNotify_StartAndStop verifies that channels are added and removed from events map properly.
func TestGoNotify_StartAndStop(t *testing.T) {
	events = make(map[string][]chan interface{})
	outputChan := make(chan interface{})

	gn.Start("test_event", outputChan)
	if channels, _ := events["test_event"]; len(channels) == 0 {
		t.Error("expected channel to be stored for given event, got channel not found")
	}

	gn.Stop("test_event", outputChan)
	if channels, _ := events["test_event"]; len(channels) != 0 {
		t.Error("expected channel to be removed for given event, got channel found")
	}
	events = make(map[string][]chan interface{})
}

// TestGoNotify_Integration makes sure notifier's Post, Start and Stop
// methods are working as intended.
// Scenario: ticker ticks and produces data which notifier posts 3 times to channel before
// configuring it to stop listening for ticket events.
func TestGoNotify_Integration(t *testing.T) {
	var observedEventsCount uint8
	myObserverChan := make(chan interface{})
	ticker := time.NewTicker(100 * time.Millisecond)

	// Producer of "test_event".
	go func() {
		for {
			t := <-ticker.C
			gn.Post("test_event", t)
		}
	}()

	// Observer of "test_event" (normally some independent component that
	// needs to be notified when "test_event" occurs).
	gn.Start("test_event", myObserverChan)

	// Block using observedEventsCount as another goroutine would cause race condition to occur.
	for observedEventsCount != 3 {
		<-myObserverChan
		observedEventsCount++
	}

	// Stop observing "test_event" and check that observedEventsCount did not increase.
	gn.Stop("test_event", myObserverChan)
	if observedEventsCount != 3 {
		t.Errorf("exepcted 3 events ocurring and being observed, got %v", observedEventsCount)
	}

	// Wait for 3 tick events and confirm that events are not being observed anymore.
	time.Sleep(350 * time.Millisecond)
	if observedEventsCount != 3 {
		t.Errorf("expected 3 events ocurring and being observed, got %v", observedEventsCount)
	}

	// Start listening again to confirm we can observe events again.
	gn.Start("test_event", myObserverChan)
	select {
	case <-myObserverChan:
		ticker.Stop()
		return
	case <-time.After(101 * time.Millisecond):
		t.Errorf("events are not being observed after restart")
	}
}
