2024/02/04
nats-server와 연동해서 메세지를 처리하는 데몬 프로세스 코드
1) eif 프로세스로부터 다량의 메시지를 수신한다. 단, Request가 아닌 Subcriber 사용함
2) A 프로세스로부터 수신한 메세지의 유효성을 체크한다.
   메세지가 유효하지 않는 경우, 에러를 리턴함
3) 메세지가 유효한 경우, 메세지를 분석해서 aaa or bbb 프로세스로 전달한다. 
   단, Request가 아닌 Subcriber 사용함 
   또한 eif로부터 수신한 메세지와 aaa or bbb 프로세스로 전달한 메세지를 묶어서 저장함
4) aaa or bbb 프로세스로부터 메세지 수신을 기다린다.
   전송한 메시지에 대해서 일정시간 동안 응답이 없는 경우, eif 프로세스로 에러 메세지를 전송한다.
5) aaa or bbb 프로세스로부터 메세지를 수신하면, 
   메시지를 검사해서 A 프로세스로 받은 메세지를 찾아서 응답하거나 다시 aaa or bbb 프로세스로 전달한다.

-----------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

// Message structure to store received and delivered messages
type Message struct {
	Content []byte
	Time    time.Time
}

var receivedMessages []Message
var deliveredMessages []Message
var mu sync.Mutex

func main() {
	// Connect to NATS server
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to messages from the eif process
	sub, err := nc.Subscribe("eif.subject", func(msg *nats.Msg) {
		// Step 2: Check the validity of the message from process A
		if isValid(msg.Data) {
			// Step 3: Analyze and deliver the message to aaa or bbb process
			deliverMessageToProcess(msg.Data)
		} else {
			// Invalid message, return an error
			log.Println("Invalid message received. Returning error.")
			// You may want to implement error handling or notify the sender
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Handle termination signals gracefully
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	signal.Notify(signalCh, os.Kill)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		<-signalCh
		fmt.Println("Daemon process shutting down.")
		wg.Done()
	}()

	// Wait for termination signal
	wg.Wait()
}

// isValid checks the validity of the received message
func isValid(message []byte) bool {
	// Implement your message validation logic here
	// Return true if the message is valid, false otherwise
	return true
}

// deliverMessageToProcess analyzes and delivers the message to aaa or bbb process
func deliverMessageToProcess(message []byte) {
	// Implement your message processing logic here
	// You may choose to forward the message to aaa or bbb process
	// and wait for the response using a subscriber
	// For example, subscribe to "aaa.process.response" or "bbb.process.response"
	sub, err := nc.Subscribe("aaa.process.response", func(msg *nats.Msg) {
		mu.Lock()
		defer mu.Unlock()

		// Step 5: Inspect the response and respond to process A or forward to aaa or bbb process
		inspectAndRespond(msg.Data, message)
	})
	if err != nil {
		log.Println("Error subscribing to response:", err)
		// Handle the error accordingly
		return
	}
	defer sub.Unsubscribe()

	// Step 4: Wait for message reception from the aaa or bbb process with a timeout
	timeout := time.After(10 * time.Second) // Adjust the timeout period as needed
	select {
	case msg := <-sub.Channel():
		mu.Lock()
		defer mu.Unlock()

		// Store delivered message
		deliveredMessages = append(deliveredMessages, Message{Content: msg.Data, Time: time.Now()})

		// Step 5: Inspect the response and respond to process A or forward to aaa or bbb process
		inspectAndRespond(msg.Data, message)
	case <-timeout:
		// If there is no response within the timeout period, send an error message to the eif process
		sendErrorMessageToEIF("Timeout waiting for response from aaa or bbb process", message)
	}
}

// inspectAndRespond inspects the response and responds to process A or forwards to aaa or bbb process
func inspectAndRespond(response []byte, originalMessage []byte) {
	// Implement your logic to inspect the response and respond accordingly
	// You may choose to send the response back to process A or forward it to aaa or bbb process
	// based on your requirements

	// Example: If the response matches the original message, respond to process A
	if string(response) == string(originalMessage) {
		fmt.Println("Responding to process A:", string(response))
		// Implement logic to respond to process A
	} else {
		fmt.Println("Forwarding to aaa or bbb process:", string(response))
		// Implement logic to forward to aaa or bbb process
	}
}

// sendErrorMessageToEIF sends an error message to the eif process
func sendErrorMessageToEIF(errorMessage string, originalMessage []byte) {
	// Implement your logic to send an error message to the eif process
	// Use NATS Publish function for this purpose
	mu.Lock()
	defer mu.Unlock()

	// Store the original message with the error information
	receivedMessages = append(receivedMessages, Message{Content: append([]byte(errorMessage), originalMessage...), Time: time.Now()})

	err := nc.Publish("eif.error.subject", []byte(errorMessage))
	if err != nil {
		log.Println("Error sending error message to eif process:", err)
		// Handle the error accordingly
	}
}