package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"light-actor-go/actor"
)

type Msg struct {
	Sender actor.PID
}
type Start struct {
	Sender actor.PID
}

type PingActor struct {
	count        int
	wgStop       *sync.WaitGroup
	messageCount int
	batch        int
	batchSize    int
}

type PongActor struct{}

func (p *PongActor) Receive(ctx actor.ActorContext) {
	switch msg := ctx.Message().(type) {
	case *Msg:
		ctx.Send(&Msg{Sender: ctx.Self()}, msg.Sender)
	}
}

func (p *PingActor) sendBatch(ctx actor.ActorContext, sender actor.PID) bool {
	if p.messageCount == 0 {
		return false
	}

	msg := &Msg{
		Sender: ctx.Self(),
	}

	for i := 0; i < p.batchSize; i++ {
		ctx.Send(msg, sender)
	}

	p.messageCount -= p.batchSize
	p.batch = p.batchSize
	return true
}

func (p *PingActor) Receive(ctx actor.ActorContext) {
	switch msg := ctx.Message().(type) {
	case *Start:
		p.sendBatch(ctx, msg.Sender)
	case *Msg:
		p.batch--
		if p.batch > 0 {
			return
		}

		if !p.sendBatch(ctx, msg.Sender) {
			p.wgStop.Done()
		}
	}
}

func NewPingActor(stop *sync.WaitGroup, messageCount int, batchSize int) *PingActor {
	return &PingActor{
		wgStop:       stop,
		messageCount: messageCount,
		batchSize:    batchSize,
	}
}

var (
	cpuprofile   = flag.String("cpuprofile", "", "write cpu profile to file")
	blockProfile = flag.String("blockprof", "", "execute contention profiling and save results here")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *blockProfile != "" {
		prof, err := os.Create(*blockProfile)
		if err != nil {
			log.Fatal(err)
		}
		runtime.SetBlockProfileRate(1)
		defer func() {
			pprof.Lookup("block").WriteTo(prof, 0)
		}()
	}

	system := actor.NewActorSystem()

	var wg sync.WaitGroup

	messageCount := 1000000
	batchSize := 100
	tps := []int{300, 400, 500, 600, 700, 800, 900}
	log.Println("Dispatcher   Throughput\tElapsed Time\tMessages per sec")
	for _, tp := range tps {

		clientCount := runtime.NumCPU() * 2
		clients := make([]actor.PID, clientCount)
		echos := make([]actor.PID, clientCount)
		for i := 0; i < clientCount; i++ {
			client, _ := system.SpawnActor(NewPingActor(&wg, messageCount, batchSize))
			echo, _ := system.SpawnActor(&PongActor{})
			clients[i] = client
			echos[i] = echo
			wg.Add(1)
		}

		start := time.Now()

		for i := 0; i < clientCount; i++ {
			system.Send(actor.NewEnvelope(&Start{Sender: echos[i]}, clients[i]))
		}

		wg.Wait()
		elapsed := time.Since(start)
		x := int(float32(messageCount*2*clientCount) / (float32(elapsed) / float32(time.Second)))
		log.Printf("%v\t%s\t%v", tp, elapsed, x)
		for i := 0; i < clientCount; i++ {
			system.GracefulStop(clients[i])
			system.GracefulStop(echos[i])
		}
		runtime.GC()
		time.Sleep(2 * time.Second)
	}
}
