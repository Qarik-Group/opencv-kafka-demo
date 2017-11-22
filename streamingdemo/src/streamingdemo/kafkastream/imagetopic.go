package kafkastream

import (
	"fmt"
	"os"
	"streamingdemo/mjpeg"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ImageTopic holds a kafka.Consumer to a Kafka topic within which are images
type ImageTopic struct {
	TopicName   string
	Hostname    string
	MJPEGStream *mjpeg.Stream
	Consumer    *kafka.Consumer
}

// NewImageTopic constructs an ImageTopic and starts the kafka.Consumer
func NewImageTopic(hostname, topicName string, mjpegStream *mjpeg.Stream) (topic *ImageTopic) {
	topic = &ImageTopic{
		Hostname:    hostname,
		TopicName:   topicName,
		MJPEGStream: mjpegStream,
	}
	var err error
	topic.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    topic.Hostname,
		"group.id":             "streamingdemo",
		"session.timeout.ms":   6000,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created '%s' consumer %v\n", topic.TopicName, topic.Consumer)
	return
}
