package remote_benchmark_test

import (
	"fmt"
	"light-actor-go/actor"
	"light-actor-go/remote"
	"net"
	"testing"
	"time"
)

type BenchmarkActor struct {
	count int
	done  chan struct{}
}

func (a *BenchmarkActor) Receive(ctx actor.ActorContext) {
	switch ctx.Message().(type) {
	case actor.SystemMessage:
		return
	default:
		a.count++
		if a.count >= totalMessages {
			a.done <- struct{}{}
		}
	}
}

const totalMessages = 100

func getFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port), nil
}

func BenchmarkRemoteActorTest(b *testing.B) {
	port1, err := getFreePort()
	if err != nil {
		b.Fatalf("failed to get free port: %v", err)
	}
	port2, err := getFreePort()
	if err != nil {
		b.Fatalf("failed to get free port: %v", err)
	}

	system := actor.NewActorSystem()
	remote1 := remote.NewRemote(*remote.NewRemoteConfig("127.0.0.1:" + port1), system)
	remote1.Listen()
	//defer remote1.Shutdown() // Ensure the listener is properly closed
	time.Sleep(time.Second)
	// Create and register the receiver actor
	receiverActor := &BenchmarkActor{done: make(chan struct{})}
	receiverPID, err := system.SpawnActor(receiverActor)
	if err != nil {
		b.Fatalf("failed to spawn receiver actor: %v", err)
	}
	remote1.MakeActorDiscoverable(receiverPID, "BenchmarkReceiver")

	// Create a remote actor system
	remoteSystem := actor.NewActorSystem()
	remote2 := remote.NewRemote(*remote.NewRemoteConfig("127.0.0.1:" + port2), remoteSystem)
	remote2.Listen()
	//defer remote2.Shutdown() // Ensure the listener is properly closed
	remotePID, err := remote2.SpawnRemoteActor("127.0.0.1:"+port1, "BenchmarkReceiver")
	if err != nil {
		b.Fatalf("failed to spawn remote actor: %v", err)
	}

	// Wait for the remote actor to be registered
	time.Sleep(time.Second)

	b.ResetTimer()
	start := time.Now()

	fmt.Printf("Sending on PORT %s \n", port1)

	for i := 0; i < totalMessages; i++ {
		remoteSystem.Send(actor.NewEnvelope(&StringMessage{Value: fmt.Sprintf("message-%d", i)}, remotePID))
	}

	// Wait for all messages to be processed or timeout after 10 seconds
	select {
	case <-receiverActor.done:
		// All messages processed
	case <-time.After(10 * time.Second):
		b.Fatalf("timeout: not all messages processed after 10 seconds")
	}

	elapsed := time.Since(start)
	b.StopTimer()

	fmt.Printf("Processed %d messages in %v\n", totalMessages, elapsed)
	fmt.Printf("Messages per second: %f\n", float64(totalMessages)/elapsed.Seconds())
}
