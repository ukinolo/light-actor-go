package main

import (
	"fmt"
	"light-actor-go/actor"
	"time"
)

// ParentActor is a top-level actor that spawns a child actor
type ParentActor struct{}

func (a *ParentActor) Receive(ctx actor.ActorContext) {

	switch msg := ctx.Message().(type) {
	case string:
		fmt.Println("Parent actor received:", msg)
		if msg == "SpawnChild" {
			self := ctx.Self()
			childProps := actor.NewActorProps(&self) // Child actor will not have a parent

			actorSystem := ctx.ActorSystem()

			childPID, _ := ctx.SpawnActor(&ChildActor{}, *childProps)

			time.Sleep(1 * time.Second)

			actorSystem.Send(actor.NewEnvelope("SpawnGrandchild", childPID))
		}
	}
}

// ChildActor is spawned by ParentActor and spawns another actor
type ChildActor struct{}

func (a *ChildActor) Receive(ctx actor.ActorContext) {

	switch msg := ctx.Message().(type) {
	case string:
		fmt.Println("Child actor received:", msg)
		if msg == "SpawnGrandchild" {
			self := ctx.Self()
			grandChildProps := actor.NewActorProps(&self)
			grandChildPID, _ := ctx.SpawnActor(&GrandChildActor{}, *grandChildProps)
			actorSystem := ctx.ActorSystem()

			time.Sleep(1 * time.Second)

			actorSystem.Send(actor.NewEnvelope("Hello from child to grandchild", grandChildPID))
		}
	case actor.SystemMessage:
		switch msg.Type {
		case actor.SystemMessageStop:
			fmt.Println("Child actor received stop message")
		case actor.SystemMessageGracefulStop:
			fmt.Println("Child actor received graceful stop message")
		case actor.SystemMessageChildTerminated:
			fmt.Println("Child actor received child terminated message")
		case actor.SystemMessageStart:
			fmt.Println("Child actor received start message")
		}
	}
}

// GrandChildActor is spawned by ChildActor
type GrandChildActor struct{}

func (a *GrandChildActor) Receive(ctx actor.ActorContext) {
	fmt.Println("Grandchild actor received message:", ctx.Message())

	switch msg := ctx.Message().(type) {
	case string:
		fmt.Println("Grandchild actor received:", msg)
	}
}

func main() {
	actorSystem := actor.NewActorSystem()

	// Spawn the top-level parent actor
	parentPID, err := actorSystem.SpawnActor(&ParentActor{})
	if err != nil {
		fmt.Println("Error spawning parent actor:", err)
		return
	}

	time.Sleep(1 * time.Second)

	// Send a message to parent actor to trigger spawning a child
	actorSystem.Send(actor.NewEnvelope("SpawnChild", parentPID))

	// Wait for a while to simulate some work happening
	time.Sleep(3 * time.Second)

	// Trigger graceful stop of parent actor
	actorSystem.GracefulStop(parentPID)

	// Wait for a while to see the output
	time.Sleep(3 * time.Second)
}
