package main

import (
	"fmt"
	"light-actor-go/actor"
	"light-actor-go/cluster"
	"light-actor-go/remote"
	"time"
)

type SimpleActor struct{}

func (*SimpleActor) Receive(ctx *actor.ActorContext) {
	switch msg := ctx.Message().(type) {
	case actor.ActorStarted:
		fmt.Println("Pokrenut sam")
	case *cluster.StringMessage:
		fmt.Println("Dosla mi je poruka:", msg)
	case actor.ActorStoped:
		fmt.Println("Zaustavljen sam")
	default:
		fmt.Printf("Doslo je nesto tipa %T, a sadrzaj je %v\n", msg, msg)
	}
}

func main() {
	c1 := cluster.NewClusterMember(remote.RemoteConfig{Addr: "127.0.0.1:8081"})
	c1.ExposeActor("Simple Actor", func() actor.Actor { return &SimpleActor{} })
	c1.Start()

	c2 := cluster.NewClusterMember(remote.RemoteConfig{Addr: "127.0.0.1:8082"})
	c2.ExposeActor("Simple Actor", func() actor.Actor { return &SimpleActor{} })
	c2.Start()

	c3 := cluster.NewClusterMember(remote.RemoteConfig{Addr: "127.0.0.1:8083"})
	c3.ExposeActor("Simple Actor", func() actor.Actor { return &SimpleActor{} })
	c3.Start()

	c4 := cluster.NewClusterMember(remote.RemoteConfig{Addr: "127.0.0.1:8084"})
	c4.ExposeActor("Simple Actor", func() actor.Actor { return &SimpleActor{} })
	c4.Start()

	//Wait for all clusters to come alive
	time.Sleep(time.Second)

	c2.JoinCluster("127.0.0.1:8081")
	c3.JoinCluster("127.0.0.1:8081")
	c4.JoinCluster("127.0.0.1:8083")
	time.Sleep(time.Second * 10)

	cc := cluster.NewClusterConnection(*remote.NewRemoteConfig("127.0.0.1:8085"))
	cc.ConnectTo("127.0.0.1:8084")
	pid1, err := cc.SpawnActor("Simple Actor")
	if err != nil {
		fmt.Println(err)
	}
	pid2, err := cc.SpawnActor("Simple Actor")
	if err != nil {
		fmt.Println(err)
	}
	pid3, err := cc.SpawnActor("Simple Actor")
	if err != nil {
		fmt.Println(err)
	}
	pid4, err := cc.SpawnActor("Simple Actor")
	if err != nil {
		fmt.Println(err)
	}
	pid5, err := cc.SpawnActor("Simple Actor")
	if err != nil {
		fmt.Println(err)
	}
	pid6, err := cc.SpawnActor("Simple Actor")
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second * 1)

	cc.Send(&cluster.StringMessage{Value: "Idemo radi 1"}, pid1)
	cc.Send(&cluster.StringMessage{Value: "Idemo radi 2"}, pid2)
	cc.Send(&cluster.StringMessage{Value: "Idemo radi 3"}, pid3)
	cc.Send(&cluster.StringMessage{Value: "Idemo radi 4"}, pid4)
	cc.Send(&cluster.StringMessage{Value: "Idemo radi 5"}, pid5)
	cc.Send(&cluster.StringMessage{Value: "Idemo radi 6"}, pid6)

	cc.Shutdown()

	time.Sleep(time.Second / 5)

}
