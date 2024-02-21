package main

import (
	"context"
	"encoding/json"
	m "enfr/message"
	"enfr/shared"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func recv_eif_msg(nc *nats.Conn, msgStore *MsgStore) {

	eifCh := make(chan *nats.Msg)
	msgId := NewMsgId()

	subEIF, err := nc.Subscribe(shared.IOS_subject, func(msg *nats.Msg) {
		eifCh <- msg
	})
	if err != nil {
		log.Printf("nc.Subscribe(%v) fail: %v\n", shared.Eif_subject, err)
	}
	defer subEIF.Unsubscribe()

	for i := 0; i < shared.EifCh_thread_cnt; i++ {
		go handle_eif_msg(nc, eifCh, msgStore, msgId)
	}

	fmt.Printf("")
}

func handle_eif_msg(nc *nats.Conn, eifCh <-chan *nats.Msg, msgStore *MsgStore, msgId *MsgId) {

	for msg := range eifCh {

		log.Printf("Received message from eif: %s\n", msg.Data)

		// Check EIF Message
		recvMsg, result := check_eif_msg(msg)

		if result {

			messageId := msgId.NextId()
			log.Printf("messageId = %#v\n", messageId)

			// Send To AAA
			aaaSendmsg := m.NatsMsg{
				Subject:       shared.AAA_subject,
				ReturnSubject: shared.IOS_return_subject,
				MsgSeqNum:     messageId,
				SendTime:      time.Now(),
				Contents:      []byte(shared.IOS_return_subject)}

			jsonData, _ := json.Marshal(aaaSendmsg)

			log.Printf("Send message to AAA: %s\n", jsonData)

			err := nc.Publish(shared.AAA_subject, jsonData)
			if err != nil {
				log.Printf("nc.Publish(%v) fail: %v\n", shared.AAA_subject, err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
			msgStore.AddMsgStore(messageId, recvMsg, aaaSendmsg, cancel)

			go func(ctx context.Context, messageId int) {
				<-ctx.Done()
				log.Printf("messageId(%v) ctx.Done()\n", messageId)
				msgStore.RemoveMsgStore(messageId)
				if ctx.Err() == context.DeadlineExceeded {
					log.Printf("Timeout occurred for message: %v\n", messageId)
				}
			}(ctx, messageId)

		} else {
			// Return To EIF
		}
	}
}

func check_eif_msg(msg *nats.Msg) (m.NatsMsg, bool) {
	var eifRecvMsg m.NatsMsg
	err := json.Unmarshal(msg.Data, &eifRecvMsg)
	if err != nil {
		log.Printf("json.Unmarshal() fail: %v\n", err)
		return eifRecvMsg, false
	}

	// Do Task

	// 메시지 저장

	return eifRecvMsg, true
}
