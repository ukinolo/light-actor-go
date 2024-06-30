package cluster

import (
	"fmt"
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
		fmt.Println("Dosao je actorCreated")
		fmt.Printf("Poruka je %v\n", msg)
		for i, v := range ConnectionActor.registry.mapping {
			fmt.Printf("mapping[%v]=%v\n", i, v)
		}
		pid, ok := ConnectionActor.registry.GetPid(strings.Split(msg.Name, "_")[0])
		if !ok {
			return
		}
		ConnectionActor.registry.Add(pid, actorInfo{
			address: msg.Address,
			name:    msg.Name,
			created: true,
		})
		fmt.Println("Ovo je posle")
		for i, v := range ConnectionActor.registry.mapping {
			fmt.Printf("mapping[%v]=%v\n", i, v)
		}

	case actor.ActorStoped:

	default:
	}
}
