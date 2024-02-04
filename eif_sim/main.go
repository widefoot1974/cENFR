package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

var NATS_URL = "localhost:4222"

type NatsMsg struct {
	MsgSeq   int
	Contents []byte
	Time     time.Time
}

func set_log() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func main() {

	set_log()

	tps := 1
	count := 100

	if len(os.Args) > 1 {
		tps, _ = strconv.Atoi(os.Args[1])
	} else if len(os.Args) > 2 {
		count, _ = strconv.Atoi(os.Args[2])
	} else {
		log.Printf("Need #tps and #count!")
		return
	}

	log.Printf("tps = %v, count = %v\n", tps, count)

	// nats-server 연동
	nc, err := nats.Connect(NATS_URL)
	if err != nil {
		log.Printf("nats.Connect(%v) fail: %v\n", NATS_URL, err)
		return
	}
	defer nc.Close()

	// eif 메시지 수신
	eif_subject := "elf.subject"
	for i := 0; i < count; i++ {

		natsMsg := NatsMsg{MsgSeq: i, Contents: []byte(eif_subject), Time: time.Now()}
		jsonData, _ := json.Marshal(natsMsg)

		err := nc.Publish(eif_subject, []byte(jsonData))
		if err != nil {
			log.Printf("nc.Publish(%v) fail: %v\n", eif_subject, err)
			return
		}
		log.Printf("MsgSeq(%v), time(%v) Sended.\n", natsMsg.MsgSeq, natsMsg.Time)

		time.Sleep(1 * time.Second)
	}

	fmt.Printf("")
}
