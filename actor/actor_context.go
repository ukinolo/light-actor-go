package actor

import (
	"context"
)

// Define the actorState constants
type actorState int32

const (
	actorStart actorState = iota
	actorStop
)

// ActorContext holds the state and context of an actor
type actorContext struct {
	actorSystem *ActorSystem
	actor       Receiver
	ctx         context.Context
	props       *Props
	message     Envelope
	state       actorState
	mailboxChan chan interface{}
	behavior    *Behavior // Add behavior as an attribute
}

// NewActorContext creates and initializes a new actorContext
func NewActorContext(ctx context.Context, actorSystem *ActorSystem, props *Props) *actorContext {
	context := new(actorContext)
	context.ctx = ctx
	context.props = props
	context.state = actorStart
	context.actorSystem = actorSystem
	context.mailboxChan = make(chan interface{}, 1)
	context.behavior = NewBehavior() // Initialize behavior
	return context
}

// Context returns the context attached to the actorContext
func (ctx *actorContext) Context() context.Context {
	return ctx.ctx
}

// Actor returns the actor associated with the actorContext
func (ctx *actorContext) Actor() Receiver {
	return ctx.actor
}

// Message returns the current message being processed
func (ctx *actorContext) Message() Envelope {
	return ctx.message
}

// State returns the current state of the actor
func (ctx *actorContext) State() actorState {
	return ctx.state
}

// Become sets the actor's behavior to a new one
func (ctx *actorContext) Become(receive ReceiveFunc) {
	ctx.behavior.Become(receive)
}

// BecomeStacked adds a new behavior on top of the current behavior stack
func (ctx *actorContext) BecomeStacked(receive ReceiveFunc) {
	ctx.behavior.BecomeStacked(receive)
}

// UnbecomeStacked removes the top behavior from the behavior stack
func (ctx *actorContext) UnbecomeStacked() {
	ctx.behavior.UnbecomeStacked()
}

// Receive processes a message using the current behavior
func (ctx *actorContext) Receive() {
	ctx.behavior.Receive(*ctx)
}
