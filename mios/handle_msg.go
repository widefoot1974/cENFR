package main

import (
	"encoding/json"
	"enfr/shared"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const MSG_CONT int = 1
const MSG_ERR int = -1

func check_eif_msg(msg *nats.Msg) (int, int) {
	var eifRecvMsg shared.NatsMsg
	err := json.Unmarshal(msg.Data, &eifRecvMsg)
	if err != nil {
		log.Printf("json.Unmarshal() fail: %v\n", err)
		return MSG_ERR, 0
	}

	// Do Task

	return MSG_CONT, eifRecvMsg.MsgSeqNum
}

func check_aaa_msg(msg *nats.Msg) (int, int) {
	var aaaRecvMsg shared.NatsMsg
	err := json.Unmarshal(msg.Data, &aaaRecvMsg)
	if err != nil {
		log.Printf("json.Unmarshal() fail: %v\n", err)
		return MSG_ERR, 0
	}

	// Do Task

	return MSG_CONT, aaaRecvMsg.MsgSeqNum
}

func handle_eif_msg(nc *nats.Conn, eifCh <-chan *nats.Msg) {
	for msg := range eifCh {

		log.Printf("Received message from eif: %s\n", msg.Data)

		// Check EIF Message
		result, msgSeqNum := check_eif_msg(msg)

		if result == MSG_CONT {
			// Send To AAA
			aaaSendmsg := shared.NatsMsg{
				Subject:       shared.AAA_subject,
				ReturnSubject: shared.IOS_return_subject,
				MsgSeqNum:     msgSeqNum,
				SendTime:      time.Now(),
				Contents:      []byte(shared.IOS_return_subject)}

			jsonData, _ := json.Marshal(aaaSendmsg)

			log.Printf("Send message to AAA: %s\n", jsonData)

			err := nc.Publish(shared.AAA_subject, jsonData)
			if err != nil {
				log.Printf("nc.Publish(%v) fail: %v\n", shared.AAA_subject, err)
				return
			}

		} else {
			// Return To EIF
		}
	}
}

func handle_aaa_msg(nc *nats.Conn, aaaCh <-chan *nats.Msg) {
	for msg := range aaaCh {

		log.Printf("Received message from aaa: %s\n", msg.Data)

		// Check AAA Message
		result, msgSeqNum := check_aaa_msg(msg)

		if result == MSG_CONT {
			// Send To EIF
			eifSendmsg := shared.NatsMsg{
				Subject:       shared.Eif_return_subject,
				ReturnSubject: "",
				MsgSeqNum:     msgSeqNum,
				SendTime:      time.Now(),
				Contents:      []byte(shared.Eif_return_subject)}

			jsonData, _ := json.Marshal(eifSendmsg)

			log.Printf("Send message to EIF: %s\n", jsonData)

			err := nc.Publish(shared.Eif_return_subject, jsonData)
			if err != nil {
				log.Printf("nc.Publish(%v) fail: %v\n", shared.Eif_return_subject, err)
				return
			}

		} else {
			// Send To EIF
		}
	}
}
