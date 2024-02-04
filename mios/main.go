package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/nats-io/nats.go"
)

var NATS_URL string = "localhost:4222"

// 로그 설정
func set_log() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func handle_eif_msg(msg *nats.Msg) {
	log.Printf("Received message: %s\n", msg.Data)
	return
}

func main() {

	set_log()

	proc_name := os.Args[0]
	log.Println("############################################################")
	log.Printf(" [%v] Started.\n", proc_name)
	log.Println("############################################################")

	// nats-server 연동
	nc, err := nats.Connect(NATS_URL)
	if err != nil {
		log.Printf("nats.Connect(%v) fail: %v\n", NATS_URL, err)
		return
	}
	defer nc.Close()

	// eif 메시지 수신
	eif_subject := "elf.subject"
	subEIF, err := nc.Subscribe(eif_subject, handle_eif_msg)
	if err != nil {
		log.Printf("nc.Subscribe(%v) fail: %v\n", eif_subject, err)
	}
	defer subEIF.Unsubscribe()

	// aaa 메세지 수신

	// Handle terminate signal gracefully
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	signal.Notify(signalCh, os.Kill)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		signal := <-signalCh
		log.Printf("signal(%v) received.\n", signal)
		wg.Done()
	}()

	wg.Wait()

	log.Println("############################################################")
	log.Printf(" [%v] Ended.\n", proc_name)
	log.Println("############################################################")
}
