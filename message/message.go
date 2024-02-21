package message

import (
	"time"
)

type NatsMsg struct {
	Subject       string    `json:"subject"`
	ReturnSubject string    `json:"return_subject"`
	MsgSeqNum     int       `json:"msg_seq_num"`
	SendTime      time.Time `json:"send_time"`
	Contents      []byte    `json:"contents"`
}
