package localkafka

type KafkaEvent struct {
	Records map[string][]KafkaMessage `json:"records"`
	// Records map[string][]string
}

type KafkaMessage struct {
	Topic     string `json:"topic"`
	Offset    int64  `json:"offset"`
	Timestamp int    `json:"timestamp"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Partition int    `json:"partition"`
}
