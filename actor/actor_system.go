package actor

import (
	"context"
	"fmt"
)

type ActorSystem struct {
	registry *Registry
}

// Creates new actor system that can only be used localy
func NewActorSystem() *ActorSystem {
	return &ActorSystem{registry: NewRegistry()}
}

func (system *ActorSystem) SpawnActor(a Actor, props ...ActorProps) (PID, error) {
	prop := ConfigureActorProps(props...)

	actorChan := make(chan Envelope)
	mailbox := NewMailbox(actorChan)

	mailboxChan := mailbox.GetChan()
	mailboxPID, err := NewPID()
	if err != nil {
		return mailboxPID, err
	}

	//Start mailbox in separate gorutine
	startMailbox(mailbox)
	// StartWorker(mailbox.Start, nil)

	//Start actor in separate gorutine

	startActor(a, system, prop, mailboxPID, actorChan)
	// StartWorker(func() {
	// 	//Setup basic actor context
	// 	actorContext := NewActorContext(context.Background(), system, prop, mailboxPID)
	// 	for {
	// 		envelope := <-actorChan
	// 		//Set only message and send
	// 		actorContext.AddEnvelope(envelope)
	// 		a.Receive(*actorContext)
	// 	}
	// }, nil)

	//Put mailbox chanel in registry
	err = system.registry.Add(mailboxPID, mailboxChan)
	if err != nil {
		return mailboxPID, err
	}

	return mailboxPID, nil
}

func startActor(a Actor, system *ActorSystem, prop *ActorProps, mailboxPID PID, actorChan chan Envelope) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Actor recovered, need restarting:", r)
			}
		}()

		//Setup basic actor context
		actorContext := NewActorContext(context.Background(), system, prop, mailboxPID)

		system.SendSystemMessage(actorContext.self, SystemMessage{Type: SystemMessageStart})
		defer system.SendSystemMessage(actorContext.self, SystemMessage{Type: SystemMessageStop})

		for {
			envelope := <-actorChan
			//Set only message and send
			actorContext.AddEnvelope(envelope)
			switch envelope.Message.(type) {
			case SystemMessage:
				actorContext.HandleSystemMessage(envelope.Message.(SystemMessage))
			default:
				a.Receive(*actorContext)
			}
		}
	}()
}

func startMailbox(mailbox *Mailbox) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Mailbox recovered, need restarting:", r)
			}
		}()
		mailbox.Start()
	}()
}

func (system *ActorSystem) Send(envelope Envelope) {
	ch := system.registry.Find(*envelope.Receiver())
	if ch == nil {
		return
	}
	ch <- envelope
}

func (system *ActorSystem) AddRemoteActor(remoteActorPID PID, senderChan chan Envelope) {
	system.registry.Add(remoteActorPID, senderChan)
}

func (system *ActorSystem) SendSystemMessage(receiver PID, msg SystemMessage) {
	envelope := NewEnvelope(msg, receiver)
	system.Send(envelope)
}
