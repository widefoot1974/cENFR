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

// 로그 설정
func set_log() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func main() {

	// 로그 초기화
	set_log()

	proc_name := os.Args[0]
	log.Println("############################################################")
	log.Printf(" [%v] Started.\n", proc_name)
	log.Println("############################################################")

	// nats-server 연동
	nc, err := nats.Connect(shared.NATS_URL)
	if err != nil {
		log.Printf("nats.Connect(%v) fail: %v\n", shared.NATS_URL, err)
		return
	}
	defer nc.Close()

	// 메시지를 저장하기 위한 MAP 구조체
	msgStore := NewMsgStore()

	// eif 메시지 수신/처리
	recv_eif_msg(nc, msgStore)

	// aaa 메세지 수신
	recv_aaa_msg(nc, msgStore)

	// Handle terminate signal gracefully
	waitForSignal()

	log.Printf("len(msgStore.canFunc) = %v\n", len(msgStore.canFunc))
	log.Println("############################################################")
	log.Printf(" [%v] Ended.\n", proc_name)
	log.Println("############################################################")
}

func waitForSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
	log.Println("\nReceived termination signal. Exiting...")
	time.Sleep(time.Second) // Give a little time to gracefully shutdown
}
