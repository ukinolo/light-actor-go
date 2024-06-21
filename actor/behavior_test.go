package actor

import (
	"fmt"
	"testing"
	"time"

	"light-actor-go/envelope"
)

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
	a.logger.Logf("Actor %s in InitialBehavior received message: %s", ctx.Self(), message)
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
	a.logger.Logf("Actor %s in AlternativeBehavior received message: %s", ctx.Self(), message)
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
	a.logger.Logf("Actor %s in StackedBehavior received message: %s", ctx.Self(), message)
	switch message {
	case "unstack":
		a.behavior.UnbecomeStacked()
	}
}

// Receive is the method that processes messages sent to the actor.
func (a *SampleActor) Receive(ctx ActorContext) {
	if a.behavior.IsEmpty() {
		a.logger.Logf("Actor %s using default behavior for message: %s", ctx.Self(), ctx.Message())
		a.behavior.Become(a.InitialBehavior)
	}
	a.behavior.Receive(ctx)
}

// MockLogger is an implementation of a logger that captures logs for testing.
type MockLogger struct {
	Logs []string
}

// NewMockLogger creates a new instance of MockLogger.
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Logf captures logs in MockLogger.
func (l *MockLogger) Logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Logs = append(l.Logs, msg)
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
	sendMessage := func(message string, expectedOutput string) {
		envelopeMsg := envelope.NewEnvelope(message, pid)
		actorSystem.Send(envelopeMsg)
		time.Sleep(500 * time.Millisecond) // Wait for processing
		lastLog := actor.logger.Logs[len(actor.logger.Logs)-1]
		if !containsSubstring(lastLog, expectedOutput) {
			t.Errorf("Expected output containing '%s', got '%s'", expectedOutput, lastLog)
		}
	}

	// Test initial behavior
	sendMessage("Hello to Actor!", "in InitialBehavior received message")

	// Test behavior switch
	sendMessage("switch", "in InitialBehavior received message")

	// Test behavior stack
	sendMessage("stack", "in AlternativeBehavior received message")

	// Test behavior unstack
	sendMessage("unstack", "in StackedBehavior received message")

	// Test behavior reset
	sendMessage("reset", "in AlternativeBehavior received message")

	sendMessage("Bye to Actor!", "in InitialBehavior received message")
}

// containsSubstring checks if a substring exists within a string.
func containsSubstring(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
