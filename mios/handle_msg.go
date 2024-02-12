package main

import (
	"context"
	"encoding/json"
	"enfr/shared"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const MSG_CONT int = 1
const MSG_ERR int = -1

func check_eif_msg(msg *nats.Msg) (shared.NatsMsg, bool) {
	var eifRecvMsg shared.NatsMsg
	err := json.Unmarshal(msg.Data, &eifRecvMsg)
	if err != nil {
		log.Printf("json.Unmarshal() fail: %v\n", err)
		return eifRecvMsg, false
	}

	// Do Task

	// 메시지 저장

	return eifRecvMsg, true
}

func check_aaa_msg(msg *nats.Msg) (shared.NatsMsg, bool) {
	var aaaRecvMsg shared.NatsMsg
	err := json.Unmarshal(msg.Data, &aaaRecvMsg)
	if err != nil {
		log.Printf("json.Unmarshal() fail: %v\n", err)
		return aaaRecvMsg, false
	}

	// Do Task

	return aaaRecvMsg, true
}

func handle_eif_msg(nc *nats.Conn, eifCh <-chan *nats.Msg, msgStore *MsgStore) {

	for msg := range eifCh {

		log.Printf("Received message from eif: %s\n", msg.Data)

		// Check EIF Message
		recvMsg, result := check_eif_msg(msg)

		if result {

			messageId := fmt.Sprintf("%d", time.Now().UnixNano())
			log.Printf("messageId = %#v\n", messageId)

			// Send To AAA
			aaaSendmsg := shared.NatsMsg{
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

			go func(ctx context.Context, messageId string) {
				<-ctx.Done()
				log.Printf("messageId(%v) ctx.Done()\n", messageId)
				msgStore.RemoveMsgStore(messageId)
				if ctx.Err() == context.DeadlineExceeded {
					log.Printf("Timeout occurred for message: %s\n", messageId)
				}
			}(ctx, messageId)

		} else {
			// Return To EIF
		}
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
			eifSendmsg := shared.NatsMsg{
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
