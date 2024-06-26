package actor

import (
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

	//Start actor in separate gorutine
	startActor(a, system, prop, mailboxPID, actorChan)

	//Put mailbox chanel in registry
	err = system.registry.Add(mailboxPID, mailboxChan)
	if err != nil {
		return mailboxPID, err
	}

	system.Send(NewEnvelope(startingActor{}, mailboxPID))

	return mailboxPID, nil
}

func startActor(a Actor, system *ActorSystem, prop *ActorProps, mailboxPID PID, actorChan chan Envelope) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Actor recovered, need restarting:", r)
				//TODO handle failure
			}
		}()

		//Setup basic actor context
		actorContext := NewActorContext(system, prop, mailboxPID)

		aProcess := actorProcess{actorContext: actorContext, mailboxPID: mailboxPID, actorChan: actorChan, actor: a}

		aProcess.start()
	}()
}

func startMailbox(mailbox *Mailbox) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Mailbox recovered, need restarting:", r)
				//TODO handle failure
			}
		}()
		mailbox.start()
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

func (system *ActorSystem) ForcefulShutdown(pid PID) {
	envelope := NewEnvelope(forcefulShutdown{}, pid)
	system.Send(envelope)
}

func (system *ActorSystem) GracefulShutdown(pid PID) {
	envelope := NewEnvelope(gracefulShutdown{}, pid)
	system.Send(envelope)
}

func (system *ActorSystem) delete(actorPID PID) {
	system.Send(NewEnvelope(closeMailbox{}, actorPID))
	close(system.registry.Find(actorPID))
	system.registry.Delete(actorPID)
}

func (system *ActorSystem) deleteWithoutSend(actorPID PID) {
	system.registry.Delete(actorPID)
}
