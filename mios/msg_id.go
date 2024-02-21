package main

import (
	"sync"
)

const MAX_MSGID_NUM = 10000

type MsgId struct {
	mu     sync.Mutex
	iMsgId int
}

func NewMsgId() *MsgId {
	return &MsgId{iMsgId: 1}
}

func (m *MsgId) NextId() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.iMsgId >= MAX_MSGID_NUM {
		m.iMsgId = 1
	} else {
		m.iMsgId++
	}

	return m.iMsgId
}
