package cluster

import (
	"fmt"
	"light-actor-go/actor"
	"math/rand"
	"strconv"
)

// type createActor struct {
// 	actorName         string
// 	connectionAddress string
// 	forwarded         bool
// }

// type deleteActor struct {
// 	actorName string
// }

// type messageActor struct {
// 	actorName string
// 	message   interface{}
// }

// type actorCreated struct {
// 	address string
// 	name    string
// }

type ClusterActor struct {
	memberList    *MemberList
	clusterSender ClusterSender
	registry      *map[string]func() actor.Actor
	actors        map[string]actor.PID
	newActorId    int
	selfAddress   string
}

func NewClusterActor(memberlist *MemberList, registry *map[string]func() actor.Actor, selfAddress string) *ClusterActor {
	return &ClusterActor{
		memberList:    memberlist,
		clusterSender: *NewClusterSender(),
		registry:      registry,
		selfAddress:   selfAddress,
		actors:        make(map[string]actor.PID),
		newActorId:    0,
	}
}

func (clusterActor *ClusterActor) Receive(context *actor.ActorContext) {
	switch msg := context.Message().(type) {
	case actor.ActorStarted:

	case *CreateActor:
		if msg.Forwarded {
			pid, err := context.SpawnActor((*clusterActor.registry)[msg.ActorName]())
			if err != nil {
				fmt.Println("Error with creating actor")
				break
			}
			clusterActor.createActor(msg, pid)
			break
		}
		possibleAddressConections := clusterActor.memberList.GetAllAlive()
		id := rand.Intn(len(possibleAddressConections) + 1)

		if id == len(possibleAddressConections) {
			pid, err := context.SpawnActor((*clusterActor.registry)[msg.ActorName]())
			if err != nil {
				fmt.Println("Error with creating actor")
				break
			}
			clusterActor.createActor(msg, pid)
			break
		}

		clusterActor.clusterSender.SendMessage(&CreateActor{
			ActorName:         msg.ActorName,
			ConnectionAddress: msg.ConnectionAddress,
			Forwarded:         true,
		}, possibleAddressConections[id], "cluster")

	case *DeleteActor:
		clusterActor.deleteActor(msg, context)

	case *MessageActor:
		clusterActor.messageActor(msg, context)

	case actor.ActorStoped:
		//TODO stop everything

	default:
	}
}

func (clusterActor *ClusterActor) createActor(message *CreateActor, pid actor.PID) {
	name := message.ActorName + "_" + strconv.Itoa(clusterActor.newActorId)
	clusterActor.newActorId += 1
	clusterActor.actors[name] = pid

	clusterActor.clusterSender.SendMessage(&ActorCreated{
		Address: clusterActor.selfAddress,
		Name:    name,
	}, message.ConnectionAddress, "connection")
	clusterActor.actors[name] = pid
}

func (clusterActor *ClusterActor) messageActor(message *MessageActor, context *actor.ActorContext) {
	pid, ok := clusterActor.actors[message.ActorName]
	if !ok {
		return
	}
	context.Send(message.Message, pid)
}

func (clusterActor *ClusterActor) deleteActor(message *DeleteActor, context *actor.ActorContext) {
	pid, ok := clusterActor.actors[message.ActorName]
	if !ok {
		return
	}
	context.GracefulShutdown(pid)
}
