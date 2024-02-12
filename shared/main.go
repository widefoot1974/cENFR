package shared

import "time"

var (
	NATS_URL           string = "localhost:4222"
	Eif_subject        string = "ToEIF"
	Eif_return_subject string = "ReturnToEIF"
	IOS_subject        string = "ToIOS"
	IOS_return_subject string = "ReturnToIOS"
	AAA_subject        string = "ToAAA"
	AAA_return_subject string = "ReturnToAAA"

	EifSim_thread_cnt int = 5
	AAASim_thread_cnt int = 5
	EifCh_thread_cnt  int = 5
	AAACh_thread_cnt  int = 5
)

type NatsMsg struct {
	Subject       string    `json:"subject"`
	ReturnSubject string    `json:"return_subject"`
	MsgSeqNum     string    `json:"msg_seq_num"`
	SendTime      time.Time `json:"send_time"`
	Contents      []byte    `json:"contents"`
}
