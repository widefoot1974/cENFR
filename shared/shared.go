package shared

var (
	NATS_URL string = "localhost:4222"

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
