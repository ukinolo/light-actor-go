package cluster

import (
	"context"
	"fmt"
	"light-actor-go/actor"
	"time"
)

type SwimActor struct {
	pingInterval  int32
	selfAddress   string
	timerCancel   context.CancelFunc
	clusterSender *ClusterSender
	swimGossiper  *SwimGossiper
}

func NewSwimActor(pingInterval int32, selfAddress string, memberList *MemberList) *SwimActor {
	return &SwimActor{
		pingInterval:  pingInterval,
		selfAddress:   selfAddress,
		swimGossiper:  NewSwimGossiper(memberList, selfAddress),
		clusterSender: NewClusterSender(),
	}
}

func (SwimActor *SwimActor) Receive(ctx *actor.ActorContext) {
	switch msg := ctx.Message().(type) {
	case actor.ActorStarted:
		timerContext, cancel := context.WithCancel(context.Background())
		SwimActor.timerCancel = cancel
		//Start gossip ping timer
		go func() {
			for {
				time.Sleep(time.Duration(SwimActor.pingInterval * int32(time.Duration(time.Second))))
				select {
				case <-timerContext.Done():
					return
				default:
				}
				ctx.Send(SendPing{}, ctx.Self())
			}
		}()
		SwimActor.clusterSender = NewClusterSender()

	case *SwimAck:
		SwimActor.swimGossiper.HandleAck(msg)

	case *SwimPing:
		respond := SwimActor.swimGossiper.HandlePing(msg)
		SwimActor.clusterSender.SendMessage(respond, msg.Sender, "gossip")

	case *SwimIndirectPing:
		if msg.Receiver == SwimActor.selfAddress {
			SwimActor.swimGossiper.HandleIndirectPing(msg)
		} else {
			SwimActor.clusterSender.SendMessage(msg, msg.Receiver, "gossip")
		}

	case *SwimIndirectAck:
		if msg.Receiver == SwimActor.selfAddress {
			SwimActor.swimGossiper.HandleIndirectAck(msg)
		} else {
			SwimActor.clusterSender.SendMessage(msg, msg.Sender, "gossip")
		}

	case SendPing:
		ping, indirectPing, address := SwimActor.swimGossiper.CreateNewPing()
		if address == "" {
			break
		}
		if indirectPing.Receiver != "" {
			SwimActor.clusterSender.SendMessage(indirectPing, address, "gossip")
		} else {
			SwimActor.clusterSender.SendMessage(ping, address, "gossip")
		}

	case JoinCluster:
		SwimActor.clusterSender.SendMessage(&SwimNewGossiper{
			MemberAddress: SwimActor.selfAddress,
			MemberState:   SwimState_Alive,
			Version:       int32(SwimActor.swimGossiper.selfVersion),
		}, msg.address, "gossip")
		SwimActor.swimGossiper.memberList.Add(msg.address, ClusterInfo{state: SwimState_Alive, version: 0})
		SwimActor.swimGossiper.healthyAddresses = append(SwimActor.swimGossiper.healthyAddresses, msg.address)

	case *SwimNewGossiper:
		SwimActor.swimGossiper.HandleNewGossiper(msg)

	case actor.ActorStoped:
		SwimActor.timerCancel()
		//TODO stop everything

	default:
		// fmt.Printf("Doslo je %T**********************************************\n", msg)
		//TODO handle random messages
	}
	fmt.Printf("%v:Doslo je %T**********************************************\n", SwimActor.selfAddress, ctx.Message())
	// fmt.Println("Trenutno stanje je ovakvo////////////////////////////////////////////")
	// fmt.Println("MemberList je ovo:")
	// for i, v := range SwimActor.swimGossiper.memberList.mapping {
	// 	fmt.Printf("Adresa %v ima vrednost %v\n", i, v)
	// }
	// fmt.Println("No responses su:")
	// for i, v := range SwimActor.swimGossiper.noResponseAddresses {
	// 	fmt.Printf("Adresa %v ima vrednost %v\n", i, v)
	// }
	// fmt.Println("HealtyAddress su:")
	// for i, v := range SwimActor.swimGossiper.healthyAddresses {
	// 	fmt.Printf("Adresa %v ima vrednost %v\n", i, v)
	// }
	// fmt.Println("Self version je:", SwimActor.swimGossiper.selfVersion)
	// fmt.Println("Extra info je: ", SwimActor.swimGossiper.extraInfo)
}
