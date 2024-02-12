package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"enfr/shared"

	"github.com/nats-io/nats.go"
)

func set_log() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func send_msg(nc *nats.Conn, tps int, count int) {

	log.Printf("tps = %v, count = %v\n", tps, count)

	// eif 메시지 수신

	for i := 0; i < count; i++ {

		natsMsg := shared.NatsMsg{
			Subject:       shared.IOS_subject,
			ReturnSubject: shared.Eif_return_subject,
			MsgSeqNum:     fmt.Sprintf("%d", i),
			SendTime:      time.Now(),
			Contents:      []byte(shared.IOS_subject)}

		jsonData, _ := json.Marshal(natsMsg)

		log.Printf("Send message to IOS: %s\n", jsonData)

		err := nc.Publish(shared.IOS_subject, jsonData)
		if err != nil {
			log.Printf("nc.Publish(%v) fail: %v\n", shared.IOS_subject, err)
			return
		}
		log.Printf("MsgSeq(%v), time(%v) Sended.\n", natsMsg.MsgSeqNum, natsMsg.SendTime)

		time.Sleep(time.Duration(tps) * time.Millisecond)
	}

	log.Printf("send_msg() Ended.")
}

func handle_recv_msg(msgCh <-chan *nats.Msg) {
	log.Printf("handle_recv_msg() started.")

	for msg := range msgCh {
		log.Printf("Received message: %s\n", msg.Data)
	}
	log.Printf("handle_recv_msg() ended.")
}

func main() {

	set_log()

	proc_name := filepath.Base(os.Args[0])
	log.Printf("(%v) started!!\n", proc_name)

	tps := 1
	count := 100

	if len(os.Args) > 2 {
		tps, _ = strconv.Atoi(os.Args[1])
		count, _ = strconv.Atoi(os.Args[2])
	} else {
		log.Printf("Need #tps and #count!")
		return
	}

	// nats-server 연동
	nc, err := nats.Connect(shared.NATS_URL)
	if err != nil {
		log.Printf("nats.Connect(%v) fail: %v\n", shared.NATS_URL, err)
		return
	}
	defer nc.Close()

	go send_msg(nc, tps, count)

	msgCh := make(chan *nats.Msg)
	sub, err := nc.Subscribe(shared.Eif_return_subject, func(msg *nats.Msg) {
		msgCh <- msg
	})
	if err != nil {
		log.Printf("nc.Subscribe(%v) fail: %v\n", shared.Eif_return_subject, err)
		return
	}

	for i := 0; i < shared.EifSim_thread_cnt; i++ {
		go handle_recv_msg(msgCh)
	}

	waitForSignal()
	sub.Unsubscribe()

	fmt.Printf("(%v) is completed!!\n", proc_name)
}

func waitForSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	fmt.Println("\nReceived termination signal. Exiting...")
	time.Sleep(time.Second) // Give a little time to gracefully shutdown
}
