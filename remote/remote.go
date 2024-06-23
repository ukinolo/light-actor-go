package remote

import (
	"fmt"
	"light-actor-go/actor"
)

type Remote struct {
	remoteReciever RemoteReceiver
	actorSystem    *actor.ActorSystem
}

func NewRemote(remoteConfing RemoteConfig, actorSystem *actor.ActorSystem) *Remote {
	return &Remote{remoteReciever: *NewRemoteReceiver(&remoteConfing, actorSystem),
		actorSystem: actorSystem,
	}
}

func (r *Remote) Listen() {
	go r.remoteReciever.startServer()
}

func (r *Remote) SpawnRemoteActor(address string, name string) (actor.PID, error) {
	newPID, err := actor.NewPID()
	if err != nil {
		return newPID, nil
	}

	remoteSender := NewRemoteSender(address)
	envelopeChan := make(chan actor.Envelope, 10)

	go func() {
		for {
			envelope := <-envelopeChan
			err := remoteSender.SendMessage(envelope.Message, name)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	r.actorSystem.AddRemoteActor(newPID, envelopeChan)
	return newPID, nil
}

func (r *Remote) MakeActorDiscoverable(actorPID actor.PID, name string) error {
	return r.remoteReciever.AddRemoteActor(name, actorPID)
}

// func (r *Remote) findActorName(actorPID actor.PID) string {
// 	return r.remoteActorRegistry.Find(actorPID)
// }
