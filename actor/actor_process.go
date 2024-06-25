package actor

import (
	"context"
	"fmt"
	"time"
)

type actorProcess struct {
	actorChan    chan Envelope
	mailboxPID   PID
	actorContext *ActorContext
	actor        Actor
}

func (aProcess *actorProcess) start() {
	for {
		envelope := <-aProcess.actorChan

		switch envelope.Message.(type) {
		case forcefulShutdown:
			aProcess.forcefullyShutdown()
			return
		case gracefulShutdown:
			aProcess.gracefullyShutdown()
			return
			//TODO more messages
		default:
			aProcess.actorContext.AddEnvelope(envelope)
			aProcess.actor.Receive(aProcess.actorContext)
		}
	}
}

func (aProcess *actorProcess) forceShutdownChildren() {
	for child := range aProcess.actorContext.children {
		aProcess.actorContext.actorSystem.Send(NewEnvelope(forcefulShutdown{}, child))
	}
}

func (aProcess *actorProcess) gracefulShutdownChildren() {
	childrenCount := len(aProcess.actorContext.children)
	if childrenCount == 0 {
		return
	}

	for child := range aProcess.actorContext.children {
		aProcess.actorContext.Send(gracefulShutdown{}, child)
	}

	ct, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for {
		select {
		case <-ct.Done():
			fmt.Println("Timeline exceded, I am no longer waithing for chldren")
			return
		case envelope := <-aProcess.actorChan:
			switch envelope.Message.(type) {
			case childTerminated:
				childrenCount -= 1
				if childrenCount == 0 {
					return
				}
			case forcefulShutdown:
				return
			}
		}
	}
}

func (aProcess *actorProcess) forcefullyShutdown() {
	aProcess.forceShutdownChildren()
	aProcess.actorContext.actorSystem.deleteWithoutSend(aProcess.actorContext.self)
}

func (aProcess *actorProcess) gracefullyShutdown() {
	aProcess.gracefulShutdownChildren()
	if aProcess.actorContext.props.parent != nil {
		aProcess.actorContext.actorSystem.Send(NewEnvelope(childTerminated{}, *aProcess.actorContext.props.parent))
	}
	//Actor reacting to the stopping
	aProcess.actorContext.AddEnvelope(NewEnvelope(ActorStoped{}, aProcess.actorContext.self))
	aProcess.actor.Receive(aProcess.actorContext)

	aProcess.actorContext.actorSystem.delete(aProcess.actorContext.self)
}
