package status

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type StatusChannel struct {
	Hostname  string
	TopicName string
	Producer  *kafka.Producer
}

// NewStatusChannel constructs StatusChannel and sets up connection
func NewStatusChannel() (statusChannel *StatusChannel) {
	statusChannel = &StatusChannel{
		Hostname:  "localhost:9092",
		TopicName: "opencv-kafka-demo-status",
	}

	cfApp, err := cfenv.Current()
	if err == nil {
		statusTopicService, err := cfApp.Services.WithName("status-topic")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot find service name 'status-topic': %v", err)
			os.Exit(1)
		}
		statusChannel.Hostname, _ = statusTopicService.CredentialString("hostname")
		statusChannel.TopicName, _ = statusTopicService.CredentialString("topicName")
	} else {
		fmt.Fprintf(os.Stderr, "Not running inside Cloud Foundry. Assume local Kafka on localhost:9092\n")
	}

	statusChannel.Producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": statusChannel.Hostname})
	if err != nil {
		fmt.Printf("Failed to create 'status-topic' producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created '%s' producer %v\n", statusChannel.TopicName, statusChannel.Producer)
	return
}

// PostStatus is an all-in-one helper to post a status message to statis-topic
func (statusChannel *StatusChannel) PostStatus(status string) {
	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)

	value := fmt.Sprintf("{\"status\":\"%s\", \"client\": \"posted_images_to_kafka\", \"language\": \"golang\"}", status)
	err := statusChannel.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &statusChannel.TopicName,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(value),
	}, deliveryChan)
	if err != nil {
		fmt.Printf("Failed .Produce for'status-topic' message: %s\n", err)
		os.Exit(1)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		// } else {
		// 	fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
		// 		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}
