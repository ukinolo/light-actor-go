package actor_spawn_benchmark_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"light-actor-go/actor"
)

type SumActor struct {
	sum       int
	count     int
	level     int
	mu        sync.Mutex
	wg        *sync.WaitGroup
	startTime time.Time // New field to store start time
}

func (a *SumActor) Receive(ctx actor.ActorContext) {
	switch msg := ctx.Message().(type) {
	case ResultMessage:
		a.mu.Lock()
		a.sum += msg.Result
		a.count++
		if a.count == 10 {
			if a.level == 1 { // Leaf level, return sum to parent
				if ctx.Self().ID == rootPID.ID {
					fmt.Println("Final Sum:", a.sum) // Print final sum at root actor
					elapsed := time.Since(a.startTime)
					fmt.Println("Time taken:", elapsed) // Print elapsed time
					a.wg.Done()                         // Signal completion
				} else {
					envelope := actor.NewEnvelope(ResultMessage{Result: a.sum}, *ctx.Props.Parent)
					ctx.ActorSystem().Send(envelope)
				}
			} else {
				envelope := actor.NewEnvelope(ResultMessage{Result: a.sum}, *ctx.Props.Parent)
				ctx.ActorSystem().Send(envelope)
			}
		}
		a.mu.Unlock()
	case StartMessage:
		if a.level <= 6 {
			for i := 0; i < 10; i++ {
				childOrdinal := msg.Ordinal*10 + i
				self := ctx.Self()
				props := actor.NewActorProps(&self)
				childPID, err := ctx.SpawnActor(&SumActor{level: a.level + 1, wg: a.wg}, *props)
				if err != nil {
					fmt.Println("Failed to spawn child:", err)
					return
				}
				envelope := actor.NewEnvelope(StartMessage{Ordinal: childOrdinal}, childPID)
				ctx.ActorSystem().Send(envelope)
			}
		} else {
			envelope := actor.NewEnvelope(ResultMessage{Result: msg.Ordinal}, *ctx.Props.Parent)
			ctx.ActorSystem().Send(envelope)
		}
	default:
		// Handle other cases or system messages
	}
}

type StartMessage struct {
	Ordinal int
}
type ResultMessage struct {
	Result int
}

var rootPID actor.PID

func runBenchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		system := actor.NewActorSystem()
		rootActor := &SumActor{level: 1, startTime: time.Now()} // Initialize startTime
		rootProps := actor.NewActorProps(nil)
		var err error

		// Add a WaitGroup to wait for completion
		var wg sync.WaitGroup
		wg.Add(1)
		rootActor.wg = &wg

		rootPID, err = system.SpawnActor(rootActor, *rootProps)
		if err != nil {
			fmt.Println("Failed to spawn root actor:", err)
			return
		}

		envelope := actor.NewEnvelope(StartMessage{Ordinal: 0}, rootPID)
		system.Send(envelope)

		// Wait for all actors to complete their work
		wg.Wait()
	}
}

func BenchmarkActorSpawn(b *testing.B) {
	runBenchmark(b)
}
