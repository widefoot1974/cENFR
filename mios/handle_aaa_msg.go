package main

import (
	"encoding/json"
	m "enfr/message"
	"enfr/shared"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func recv_aaa_msg(nc *nats.Conn, msgStore *MsgStore) {

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

}

func handle_aaa_msg(nc *nats.Conn, aaaCh <-chan *nats.Msg, msgStore *MsgStore) {
	for msg := range aaaCh {

		log.Printf("Received message from aaa: %s\n", msg.Data)

		// Check AAA Message
		recvMsg, result := check_aaa_msg(msg)

		messageId := recvMsg.MsgSeqNum
		log.Printf("messageId = %v\n", messageId)

		cancelFunc, ok := msgStore.GetCancelFunc(messageId)
		if ok {
			log.Printf("messageId(%v) is exist.\n", messageId)
			cancelFunc()
			// msgStore.RemoveMsgStore(messageId)
		} else {
			log.Printf("messageId(%v) is not exist.\n", messageId)
			return
		}

		if result {
			// Send To EIF
			eifSendmsg := m.NatsMsg{
				Subject:       shared.Eif_return_subject,
				ReturnSubject: "",
				MsgSeqNum:     recvMsg.MsgSeqNum,
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

func check_aaa_msg(msg *nats.Msg) (m.NatsMsg, bool) {
	var aaaRecvMsg m.NatsMsg
	err := json.Unmarshal(msg.Data, &aaaRecvMsg)
	if err != nil {
		log.Printf("json.Unmarshal() fail: %v\n", err)
		return aaaRecvMsg, false
	}

	// Do Task

	return aaaRecvMsg, true
}
