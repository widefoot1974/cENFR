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

	msgStore := NewMsgStore()

	// eif 메시지 수신
	eifCh := make(chan *nats.Msg)
	subEIF, err := nc.Subscribe(shared.IOS_subject, func(msg *nats.Msg) {
		eifCh <- msg
	})
	if err != nil {
		log.Printf("nc.Subscribe(%v) fail: %v\n", shared.Eif_subject, err)
	}
	defer subEIF.Unsubscribe()

	for i := 0; i < shared.EifCh_thread_cnt; i++ {
		go handle_eif_msg(nc, eifCh, msgStore)
	}

	// aaa 메세지 수신
	aaaCh := make(chan *nats.Msg)
	subAAA, err := nc.Subscribe(shared.IOS_return_subject, func(msg *nats.Msg) {
		aaaCh <- msg
	})
	if err != nil {
		log.Printf("nc.Subscribe(%v) fail: %v\n", shared.IOS_return_subject, err)
	}
	defer subAAA.Unsubscribe()

	for i := 0; i < shared.AAACh_thread_cnt; i++ {
		go handle_aaa_msg(nc, aaaCh, msgStore)
	}

	// Handle terminate signal gracefully
	waitForSignal()

	log.Printf("len(msgStore.canFunc) = %v\n", len(msgStore.canFunc))
	log.Println("############################################################")
	log.Printf(" [%v] Ended.\n", proc_name)
	log.Println("############################################################")
}

func waitForSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	log.Println("\nReceived termination signal. Exiting...")
	time.Sleep(time.Second) // Give a little time to gracefully shutdown
}
