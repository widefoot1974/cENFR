package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"enfr/shared"

	"github.com/nats-io/nats.go"
)

func init() {
	SetLog()
}

func SetLog() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func main() {

	logStartup(os.Args[0])

	// nats-server 연동
	nc := connectToNATS(shared.NATS_URL)
	defer nc.Close()

	// 메시지를 저장하기 위한 MAP 구조체
	msgStore := NewMsgStore()

	// eif 메시지 수신/처리
	handleEIFMessage(nc, msgStore)

	// aaa 메세지 수신
	handleAAAMessge(nc, msgStore)

	// Handle terminate signal gracefully
	waitForSignal()

	logShutdown(os.Args[0], len(msgStore.canFunc))
}

func logStartup(procName string) {
	log.Println("############################################################")
	log.Printf(" [%v] Started.\n", procName)
	log.Println("############################################################")
}

func logShutdown(procName string, storeSize int) {
	log.Printf("len(msgStore.canFunc) = %v\n", storeSize)
	log.Println("############################################################")
	log.Printf(" [%v] Ended.\n", procName)
	log.Println("############################################################")
}

func connectToNATS(url string) *nats.Conn {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Failed to connect to NATS at %s: %v", url, err)
	}
	return nc
}

func waitForSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
	log.Println("\nReceived termination signal. Exiting...")
	time.Sleep(time.Second) // Give a little time to gracefully shutdown
}
