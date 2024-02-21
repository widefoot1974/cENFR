package main

import (
	"context"
	"sync"

	m "enfr/message"
)

type MsgStore struct {
	mutex   sync.RWMutex
	canFunc map[int]context.CancelFunc
	recvMsg map[int]m.NatsMsg
	sendMsg map[int]m.NatsMsg
}

func NewMsgStore() *MsgStore {
	return &MsgStore{
		recvMsg: make(map[int]m.NatsMsg),
		sendMsg: make(map[int]m.NatsMsg),
		canFunc: make(map[int]context.CancelFunc),
	}
}

func (ms *MsgStore) AddMsgStore(
	messageId int,
	recvMsg m.NatsMsg,
	sendMsg m.NatsMsg,
	cancelFunc context.CancelFunc) {

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.canFunc[messageId] = cancelFunc
	ms.recvMsg[messageId] = recvMsg
	ms.sendMsg[messageId] = sendMsg
}

func (ms *MsgStore) RemoveMsgStore(messageId int) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	delete(ms.canFunc, messageId)
	delete(ms.recvMsg, messageId)
	delete(ms.sendMsg, messageId)
}

func (ms *MsgStore) GetCancelFunc(messageId int) (context.CancelFunc, bool) {
	cancelFunc, found := ms.canFunc[messageId]
	if found {
		return cancelFunc, true
	} else {
		return nil, false
	}
}
