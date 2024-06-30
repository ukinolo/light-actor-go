package cluster

import (
	"fmt"
	"light-actor-go/actor"
	"light-actor-go/remote"

	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type ClusterConnection struct {
	registry           ConnectionActorsRegistry
	clusterSender      *ClusterSender
	system             *actor.ActorSystem
	remote             *remote.Remote
	entryPointAddress  string
	selfAddress        string
	connectionActorPID actor.PID
	clusterActorPID    actor.PID
}

func NewClusterConnection(remoteConfig remote.RemoteConfig) *ClusterConnection {
	ClusterConnection := new(ClusterConnection)
	ClusterConnection.registry = *NewConnectionActorsRegistry()
	ClusterConnection.clusterSender = NewClusterSender()
	ClusterConnection.system = actor.NewActorSystem()
	ClusterConnection.remote = remote.NewRemote(remoteConfig, ClusterConnection.system)
	ClusterConnection.selfAddress = remoteConfig.Addr
	return ClusterConnection
}

func (ClusterConnection *ClusterConnection) ConnectTo(address string) {
	ClusterConnection.entryPointAddress = address
	//Create connection actor
	var err error
	ClusterConnection.connectionActorPID, err = ClusterConnection.system.SpawnActor(NewConnectionActor(&ClusterConnection.registry))
	if err != nil {
		fmt.Println("Error creating connection actor")
	}
	ClusterConnection.remote.MakeActorDiscoverable(ClusterConnection.connectionActorPID, "connection")
	ClusterConnection.clusterActorPID, err = ClusterConnection.remote.SpawnRemoteActor("cluster", address)
	if err != nil {
		fmt.Println("Error connecting to node")
	}
	ClusterConnection.remote.Listen()
}

func (ClusterConnection *ClusterConnection) Shutdown() {
	ClusterConnection.system.GracefulShutdown(ClusterConnection.connectionActorPID)
	for _, v := range ClusterConnection.registry.mapping {
		if v.address != "" {
			ClusterConnection.clusterSender.SendMessage(&DeleteActor{ActorName: v.name}, v.address, "cluster")
		}
	}
}

func (ClusterConnection *ClusterConnection) SpawnActor(name string) (actor.PID, error) {
	newPID, err := actor.NewPID()
	ClusterConnection.registry.Add(newPID, actorInfo{
		address: "",
		name:    name,
		created: false,
	})
	ClusterConnection.clusterSender.SendMessage(&CreateActor{ActorName: name, ConnectionAddress: ClusterConnection.selfAddress, Forwarded: false}, ClusterConnection.entryPointAddress, "cluster")
	return newPID, err
}

func (ClusterConnection *ClusterConnection) Send(message interface{}, actorPID actor.PID) {
	actorInfo := ClusterConnection.registry.Find(actorPID)
	if actorInfo.created {
		anyMsg, err := anypb.New(message.(proto.Message))
		if err != nil {
			return
		}
		ClusterConnection.clusterSender.SendMessage(&MessageActor{
			ActorName: actorInfo.name,
			Message:   anyMsg,
		}, actorInfo.address, "cluster")
	} else {
		fmt.Println("Actor is not created")
	}
}

func (ClusterConnection *ClusterConnection) ShutdownActor(actorPID actor.PID) {
	actorInfo := ClusterConnection.registry.Find(actorPID)
	if actorInfo.created {
		ClusterConnection.clusterSender.SendMessage(&DeleteActor{
			ActorName: actorInfo.name,
		}, actorInfo.address, "cluster")
		ClusterConnection.registry.Delete(actorPID)
	} else {
		fmt.Println("Actor is not created")
	}
}
