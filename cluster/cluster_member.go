package cluster

import (
	"light-actor-go/actor"
	"light-actor-go/remote"
)

type ClusterMember struct {
	system      *actor.ActorSystem
	remote      remote.Remote
	registry    map[string]func() actor.Actor
	selfAddress string
	gossiperPID actor.PID
}

func NewClusterMember(remoteConfig remote.RemoteConfig) *ClusterMember {
	clusterMember := new(ClusterMember)
	clusterMember.system = actor.NewActorSystem()
	clusterMember.remote = *remote.NewRemote(remoteConfig, clusterMember.system)
	clusterMember.registry = make(map[string]func() actor.Actor)
	clusterMember.selfAddress = remoteConfig.Addr
	return clusterMember
}

func (cMember *ClusterMember) Start() {
	memberList := NewMemberList()
	gossiperPID, err := cMember.system.SpawnActor(NewSwimActor(int32(1), cMember.selfAddress, memberList))
	cMember.gossiperPID = gossiperPID
	if err != nil {
		return
	}

	clusterPID, err := cMember.system.SpawnActor(NewClusterActor(memberList, &cMember.registry, cMember.selfAddress)) //TODO add
	if err != nil {
		return
	}

	cMember.remote.MakeActorDiscoverable(gossiperPID, "gossip")
	cMember.remote.MakeActorDiscoverable(clusterPID, "cluster")

	cMember.remote.Listen()
}

func (cMember *ClusterMember) ExposeActor(name string, actorCreator func() actor.Actor) {
	cMember.registry[name] = actorCreator
}

func (cMember *ClusterMember) JoinCluster(address string) {
	cMember.system.Send(actor.NewEnvelope(JoinCluster{address: address}, cMember.gossiperPID))
}

func (cMember *ClusterMember) CreateActor(name string) actor.Actor {
	return cMember.registry[name]()
}
