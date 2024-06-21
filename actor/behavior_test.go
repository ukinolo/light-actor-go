package actor

import (
	"fmt"
	"testing"
	"time"

	"light-actor-go/envelope"
)

// LogEntry represents a structured log entry for capturing actor behavior and messages.
type LogEntry struct {
	Behavior string
	Message  string
}

// String converts a LogEntry to a string.
func (le LogEntry) String() string {
	return fmt.Sprintf("Actor in %s received message: %s", le.Behavior, le.Message)
}

// SampleActor is a simple actor implementation that processes messages.
type SampleActor struct {
	name     string
	behavior *Behavior
	logger   *MockLogger // MockLogger for capturing logs during testing
}

// NewSampleActor creates a new instance of SampleActor.
func NewSampleActor(name string) *SampleActor {
	act := &SampleActor{
		name:     name,
		behavior: NewBehavior(),
		logger:   NewMockLogger(),
	}
	act.behavior.Become(act.InitialBehavior)
	return act
}

// Initial behavior for the actor
func (a *SampleActor) InitialBehavior(ctx ActorContext) {
	message := ctx.Message()
	a.logger.Log(LogEntry{"InitialBehavior", message.(string)})
	switch message {
	case "switch":
		a.behavior.Become(a.AlternativeBehavior)
	case "reset":
		a.behavior.Become(a.InitialBehavior)
	case "stack":
		a.behavior.BecomeStacked(a.StackedBehavior)
	case "unstack":
		a.behavior.UnbecomeStacked()
	}
}

// Alternate behavior for the actor
func (a *SampleActor) AlternativeBehavior(ctx ActorContext) {
	message := ctx.Message()
	a.logger.Log(LogEntry{"AlternativeBehavior", message.(string)})
	switch message {
	case "reset":
		a.behavior.Become(a.InitialBehavior)
	case "stack":
		a.behavior.BecomeStacked(a.StackedBehavior)
	case "unstack":
		a.behavior.UnbecomeStacked()
	}
}

// Stacked behavior for the actor
func (a *SampleActor) StackedBehavior(ctx ActorContext) {
	message := ctx.Message()
	a.logger.Log(LogEntry{"StackedBehavior", message.(string)})
	switch message {
	case "unstack":
		a.behavior.UnbecomeStacked()
	}
}

// Receive is the method that processes messages sent to the actor.
func (a *SampleActor) Receive(ctx ActorContext) {
	if a.behavior.IsEmpty() {
		a.logger.Log(LogEntry{"DefaultBehavior", ctx.Message().(string)})
		a.behavior.Become(a.InitialBehavior)
	}
	a.behavior.Receive(ctx)
}

// MockLogger is an implementation of a logger that captures structured logs for testing.
type MockLogger struct {
	Logs []LogEntry
}

// NewMockLogger creates a new instance of MockLogger.
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Log captures structured log entries in MockLogger.
func (l *MockLogger) Log(entry LogEntry) {
	l.Logs = append(l.Logs, entry)
}

// TestSampleActorBehavior tests the behavior of SampleActor.
func TestSampleActorBehavior(t *testing.T) {
	// Create an instance of SampleActor
	actor := NewSampleActor("Actor")

	// Initialize the actor system
	actorSystem := NewActorSystem()

	// Spawn the actor
	pid, err := actorSystem.SpawnActor(actor)
	if err != nil {
		t.Fatalf("Failed to spawn Actor: %v", err)
	}

	// Define a helper function to send messages and verify behavior
	sendMessage := func(message string, expectedLog string) {
		envelopeMsg := envelope.NewEnvelope(message, pid)
		actorSystem.Send(envelopeMsg)
		time.Sleep(500 * time.Millisecond) // Wait for processing

		if len(actor.logger.Logs) == 0 {
			t.Fatalf("No logs captured")
		}
		lastLog := actor.logger.Logs[len(actor.logger.Logs)-1]
		if lastLog.String() != expectedLog {
			t.Errorf("Expected log entry '%s', got '%s'", expectedLog, lastLog.String())
		}
	}

	// Test initial behavior
	sendMessage("Hello to Actor!", "Actor in InitialBehavior received message: Hello to Actor!")

	// Test behavior switch
	sendMessage("switch", "Actor in InitialBehavior received message: switch")

	// Test behavior stack
	sendMessage("stack", "Actor in AlternativeBehavior received message: stack")

	// Test behavior unstack
	sendMessage("unstack", "Actor in StackedBehavior received message: unstack")

	// Test behavior reset
	sendMessage("reset", "Actor in AlternativeBehavior received message: reset")

	sendMessage("Bye to Actor!", "Actor in InitialBehavior received message: Bye to Actor!")
}
