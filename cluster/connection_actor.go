package cluster

import (
	"light-actor-go/actor"
	"strings"
)

type ConnectionActor struct {
	registry *ConnectionActorsRegistry
}

func NewConnectionActor(registry *ConnectionActorsRegistry) *ConnectionActor {
	return &ConnectionActor{
		registry: registry,
	}
}

func (ConnectionActor *ConnectionActor) Receive(context *actor.ActorContext) {
	switch msg := context.Message().(type) {
	case actor.ActorStarted:

	case *ActorCreated:
		pid, ok := ConnectionActor.registry.GetPid(strings.Split(msg.Name, "_")[0])
		if !ok {
			return
		}
		ConnectionActor.registry.Add(pid, actorInfo{
			address: msg.Address,
			name:    msg.Name,
			created: true,
		})

	case actor.ActorStoped:

	default:
	}
}
