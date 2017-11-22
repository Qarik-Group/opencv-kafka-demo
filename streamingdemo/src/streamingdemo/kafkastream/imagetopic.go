package kafkastream

import (
	"fmt"
	"os"
	"os/signal"
	"streamingdemo/mjpeg"
	"syscall"

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
	topic.constructConsumer()
	topic.watchImages(topic.MJPEGStream.UpdateJPEG)
	return
}

func (topic *ImageTopic) constructConsumer() {
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
}

// WatchImages pulls images off the Kafka topic and pushes them to any listening mjpeg.Stream clients
func (topic *ImageTopic) watchImages(latestJPEG func([]byte)) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	err := topic.Consumer.SubscribeTopics([]string{topic.TopicName}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topic '%s': %v", topic.TopicName, err)
		os.Exit(1)
	}

	go func() {
		run := true

		for run == true {
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				run = false
			default:
				ev := topic.Consumer.Poll(100)
				if ev == nil {
					continue
				}

				switch e := ev.(type) {
				case *kafka.Message:
					fmt.Printf("%% Frame on %s\n", e.TopicPartition)
					latestJPEG([]byte(e.Value))
				case kafka.PartitionEOF:
					fmt.Printf("%% Reached %v\n", e)
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
					run = false
					// TODO: perhaps insert an "Error" image
				default:
					fmt.Printf("Ignored %v\n", e)
				}
			}
		}
		os.Exit(0)
	}()
}
