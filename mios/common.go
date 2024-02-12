package main

import (
	"context"
	"enfr/shared"
	"sync"
)

type MsgStore struct {
	mutex   sync.RWMutex
	canFunc map[string]context.CancelFunc
	recvMsg map[string]shared.NatsMsg
	sendMsg map[string]shared.NatsMsg
}

func NewMsgStore() *MsgStore {
	return &MsgStore{
		recvMsg: make(map[string]shared.NatsMsg),
		sendMsg: make(map[string]shared.NatsMsg),
		canFunc: make(map[string]context.CancelFunc),
	}
}

func (ms *MsgStore) AddMsgStore(messageId string,
	recvMsg shared.NatsMsg,
	sendMsg shared.NatsMsg,
	cancelFunc context.CancelFunc) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.canFunc[messageId] = cancelFunc
	ms.recvMsg[messageId] = recvMsg
	ms.sendMsg[messageId] = sendMsg
}

func (ms *MsgStore) RemoveMsgStore(messageId string) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	delete(ms.canFunc, messageId)
	delete(ms.recvMsg, messageId)
	delete(ms.sendMsg, messageId)
}

func (ms *MsgStore) GetCancelFunc(messageId string) (context.CancelFunc, bool) {
	cancelFunc, found := ms.canFunc[messageId]
	if found {
		return cancelFunc, true
	} else {
		return nil, false
	}
}
