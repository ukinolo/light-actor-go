package main

import (
	"fmt"
	"time"

	"light-actor-go/actor"
	"light-actor-go/envelope"
)

// SampleActor is a simple actor implementation that processes messages.
type SampleActor struct {
	name     string
	behavior *actor.Behavior
}

func NewSampleActor(name string) *SampleActor {
	act := &SampleActor{
		name:     name,
		behavior: actor.NewBehavior(),
	}
	act.behavior.Become(act.InitialBehavior)
	return act
}

// Initial behavior for the actor
func (a *SampleActor) InitialBehavior(ctx actor.ActorContext) {
	message := ctx.Message()
	fmt.Println("Actor", ctx.Self(), "in InitialBehavior received message:", message)
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
func (a *SampleActor) AlternativeBehavior(ctx actor.ActorContext) {
	message := ctx.Message()
	fmt.Println("Actor", ctx.Self(), "in AlternativeBehavior received message:", message)
	switch message {
	case "reset":
		a.behavior.Become(a.InitialBehavior)
	case "stack":
		a.behavior.BecomeStacked(a.StackedBehavior)
	case "unstack":
		a.behavior.UnbecomeStacked()
	}
}

func (a *SampleActor) StackedBehavior(ctx actor.ActorContext) {
	message := ctx.Message()
	fmt.Println("Actor", ctx.Self(), "in StackedBehavior received message:", message)
	switch message {
	case "unstack":
		a.behavior.UnbecomeStacked()
	}
}

// Receive is the method that processes messages sent to the actor.
func (a *SampleActor) Receive(ctx actor.ActorContext) {
	if a.behavior.IsEmpty() {
		fmt.Println("Actor", ctx.Self(), "using default behavior for message:", ctx.Message())
		a.behavior.Become(a.InitialBehavior)
	} else {
		a.behavior.Receive(ctx)
	}
}

func main() {
	// Initialize the actor system
	actorSystem := actor.NewActorSystem()

	actor := NewSampleActor("Actor2")

	// Spawn the actor
	pid, err := actorSystem.SpawnActor(actor)
	if err != nil {
		fmt.Println("Failed to spawn Actor:", err)
		return
	}

	// Send a message from Actor1 to Actor2
	helloMessage := "Hello to Actor!"
	envelopeHello := envelope.NewEnvelope(helloMessage, pid)
	actorSystem.Send(envelopeHello)

	// Send a "switch" message from Actor1 to Actor2 to change its behavior
	switchMessage := "switch"
	envelopeSwitch := envelope.NewEnvelope(switchMessage, pid)
	actorSystem.Send(envelopeSwitch)

	// Send a "stack" message from Actor1 to Actor2 to stack behavior
	stackMessage := "stack"
	envelopeStack := envelope.NewEnvelope(stackMessage, pid)
	actorSystem.Send(envelopeStack)

	time.Sleep(1 * time.Second)

	// Send an "unstack" message from Actor1 to Actor2 to unstack behavior
	unstackMessage := "unstack"
	envelopeUnstack := envelope.NewEnvelope(unstackMessage, pid)
	actorSystem.Send(envelopeUnstack)

	time.Sleep(1 * time.Second)

	// Send a "reset" message from Actor1 to Actor2 to change its behavior back
	resetMessage := "reset"
	envelopeReset := envelope.NewEnvelope(resetMessage, pid)
	actorSystem.Send(envelopeReset)

	time.Sleep(1 * time.Second)

	actorSystem.Send(envelopeUnstack)

	time.Sleep(1 * time.Second)

	actorSystem.Send(envelopeHello)

	// Allow some time for messages to be processed
	time.Sleep(2 * time.Second)

}
