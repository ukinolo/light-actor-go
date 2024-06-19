package remote_test

import (
	"light-actor-go/actor"
	"light-actor-go/remote"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRemoteHandler(t *testing.T) {
	// Define test message
	testMessage := "Hello from Client!"

	// Create channels for envelopes
	receiverEnvelopeChan := make(chan *actor.Envelope, 100)
	senderEnvelopeChan := make(chan *actor.Envelope, 100) // Start server for receiver
	receiverAddr := "127.0.0.1:8091"

	receiverConfig := remote.Configure(receiverAddr)
	receiverHandler := remote.NewRemoteHandler(receiverConfig, receiverEnvelopeChan)

	// Start server for sender
	senderAddr := "127.0.0.1:8092"
	senderConfig := remote.Configure(senderAddr)
	senderHandler := remote.NewRemoteHandler(senderConfig, senderEnvelopeChan)

	// Give servers time to start
	time.Sleep(100 * time.Millisecond)

	// Define client and receiver UUIDs
	receiverUUID := uuid.New()
	//senderUUID := uuid.New()
	stringMessage := &StringMessage{Value: testMessage}

	// Send the message from sender to receiver
	err := senderHandler.SendMessageToAddress(receiverAddr, stringMessage, receiverUUID)
	if err != nil {
		t.Fatalf("Error sending message: %v", err)
	}

	// Wait for the message to be received
	time.Sleep(100 * time.Millisecond)

	// Check if the message was received
	select {
	case receivedMsg := <-receiverHandler.MessageChan:
		if receivedMsg == nil {
			t.Fatal("Received nil message")
		}
		t.Logf("Received message doesn't match sent message. Got: %s, Want: %s", receivedMsg.Message.GetValue(), testMessage)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout while waiting for message")
	}
}
