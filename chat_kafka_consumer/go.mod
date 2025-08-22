module chat_kafka_broker

go 1.24

require (
	github.com/confluentinc/confluent-kafka-go v1.9.2
	google.golang.org/grpc v1.74.0
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250715232539-7130f93afb79 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace google.golang.org/genproto => google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a
