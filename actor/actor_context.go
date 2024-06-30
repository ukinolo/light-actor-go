package actor

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Define the actorState constants
type actorState int32

const (
	actorStart actorState = iota
	actorStop
	actorStopping
)

// ActorContext holds the state and context of an actor
type ActorContext struct {
	actorSystem *ActorSystem
	ctx         context.Context
	Props       *ActorProps
	envelope    Envelope
	state       actorState
	children    map[uuid.UUID]PID
	self        PID
	mu          sync.RWMutex
}

// NewActorContext creates and initializes a new actorContext
func NewActorContext(ctx context.Context, actorSystem *ActorSystem, Props *ActorProps, self PID) *ActorContext {
	context := new(ActorContext)
	context.ctx = ctx
	context.Props = Props
	context.state = actorStart
	context.actorSystem = actorSystem
	context.self = self
	context.children = make(map[uuid.UUID]PID) // Initialize children as a map
	return context
}

// Adds envelope to the current actor context
func (ctx *ActorContext) AddEnvelope(envelope Envelope) {
	ctx.envelope = envelope
}

// Spawns child actor
func (ctx *ActorContext) SpawnActor(actor Actor, Props ...ActorProps) (PID, error) {
	prop := ConfigureActorProps(Props...)
	prop.AddParent(&ctx.self)

	// pid, err := NewPID()
	// if err != nil {
	// 	return PID{}, err
	// }

	pid, err := ctx.actorSystem.SpawnActor(actor, *prop)
	if err != nil {
		return PID{}, err
	}

	ctx.mu.Lock()
	ctx.children[pid.ID] = pid
	ctx.mu.Unlock()

	return pid, nil
}

// Send message
func (ctx *ActorContext) Send(message interface{}, receiver PID) {
	sendEnvelope := NewEnvelope(message, receiver)
	ctx.actorSystem.Send(sendEnvelope)
}

func (ctx *ActorContext) Message() interface{} {
	return ctx.envelope.Message
}

// Context returns the context attached to the actorContext
func (ctx *ActorContext) Context() context.Context {
	return ctx.ctx
}

// Message returns the current message being processed
func (ctx *ActorContext) Envelope() Envelope {
	return ctx.envelope
}

// State returns the current state of the actor
func (ctx *ActorContext) State() actorState {
	return ctx.state
}

func (ctx *ActorContext) Self() PID {
	return ctx.self
}

func (ctx *ActorContext) ActorSystem() *ActorSystem {
	return ctx.actorSystem
}

func (ctx *ActorContext) HandleSystemMessage(msg SystemMessage) {
	switch msg.Type {
	case SystemMessageStart:
		ctx.Start()
	case SystemMessageStop:
		ctx.Stop()
	case SystemMessageGracefulStop:
		ctx.GracefulStop()
	case SystemMessageChildTerminated:
		if child, ok := msg.Extras.(PID); ok {
			ctx.ChildTerminated(child)
		} else {
			fmt.Println("System message extras not a PID")
		}
	default:
		fmt.Println("System message unknown")
	}
}

func (ctx *ActorContext) Start() {
	ctx.state = actorStart
	// fmt.Println("System message start:", ctx.self)
}

func (ctx *ActorContext) Stop() {

	ctx.mu.RLock()
	if len(ctx.children) > 0 {
		for _, child := range ctx.children {
			ctx.actorSystem.SendSystemMessage(child, SystemMessage{Type: SystemMessageStop})
		}
		ctx.mu.RUnlock()
	} else {
		ctx.mu.RUnlock()
	}

	ctx.actorSystem.RemoveActor(ctx.self, SystemMessage{Type: DeleteMailbox})
	ctx.state = actorStop
	// fmt.Println("System message stop", ctx.self)
}
func (ctx *ActorContext) GracefulStop() {
	ctx.state = actorStopping
	// fmt.Println("System message stopping", ctx.self)

	ctx.mu.RLock()

	if len(ctx.children) > 0 {
		for _, child := range ctx.children {
			ctx.actorSystem.SendSystemMessage(child, SystemMessage{Type: SystemMessageGracefulStop})
			// fmt.Println("Sending graceful stop to child:", child)
		}
		ctx.mu.RUnlock()
	} else {
		ctx.mu.RUnlock()
		if ctx.Props.Parent != nil {
			ctx.actorSystem.SendSystemMessage(*ctx.Props.Parent, SystemMessage{Type: SystemMessageChildTerminated, Extras: ctx.self})
		}
		// fmt.Println("No children, actor stopped", ctx.self)
		ctx.actorSystem.RemoveActor(ctx.self, SystemMessage{Type: DeleteMailbox})
		ctx.state = actorStop
	}
}

func (ctx *ActorContext) ChildTerminated(child PID) {
	// fmt.Println("System message child terminated:", child)

	ctx.mu.Lock()
	delete(ctx.children, child.ID)
	ctx.mu.Unlock()

	ctx.mu.RLock()
	if len(ctx.children) == 0 && ctx.state == actorStopping {
		ctx.mu.RUnlock()
		if ctx.Props.Parent != nil {
			ctx.actorSystem.SendSystemMessage(*ctx.Props.Parent, SystemMessage{Type: SystemMessageChildTerminated, Extras: ctx.self})
		}
		ctx.actorSystem.RemoveActor(ctx.self, SystemMessage{Type: DeleteMailbox})
		ctx.state = actorStop
		// fmt.Println("No more children left, actor stopping", ctx.self)
	} else {
		ctx.mu.RUnlock()
	}
}
