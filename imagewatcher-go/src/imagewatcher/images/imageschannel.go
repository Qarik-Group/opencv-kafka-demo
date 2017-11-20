package images

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ImagesChannel struct {
	Hostname  string
	TopicName string
	Consumer  *kafka.Consumer
}

// NewImagesChannel constructs ImagesChannel and sets up connection
func NewImagesChannel() (imagesChannel *ImagesChannel) {
	imagesChannel = &ImagesChannel{
		Hostname:  "localhost:9092",
		TopicName: "opencv-kafka-demo-raw-images",
	}
	if os.Getenv("IMAGE_TOPIC") != "" {
		imagesChannel.TopicName = os.Getenv("IMAGE_TOPIC")
	}

	cfApp, err := cfenv.Current()
	if err == nil {
		lookupServiceName := "raw-images-topic"
		if os.Getenv("IMAGE_TOPIC_SERVICE_NAME") != "" {
			lookupServiceName = os.Getenv("IMAGE_TOPIC_SERVICE_NAME")
		}
		imagesTopicService, err := cfApp.Services.WithName(lookupServiceName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot find service name 'images-topic': %v", err)
			os.Exit(1)
		}
		imagesChannel.Hostname, _ = imagesTopicService.CredentialString("hostname")
		imagesChannel.TopicName, _ = imagesTopicService.CredentialString("topicName")
	} else {
		fmt.Fprintf(os.Stderr, "Not running inside Cloud Foundry. Assume local Kafka on localhost:9092 on topic '%s'\n", imagesChannel.TopicName)
	}

	imagesChannel.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    imagesChannel.Hostname,
		"group.id":             "imagewatcher",
		"session.timeout.ms":   6000,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created '%s' consumer %v\n", imagesChannel.TopicName, imagesChannel.Consumer)
	return
}

// PostImage is an all-in-one helper to post an image to kafka topic
func (imagesChannel *ImagesChannel) WatchImages(latestJPEG func([]byte)) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	err := imagesChannel.Consumer.SubscribeTopics([]string{imagesChannel.TopicName}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topic '%s': %v", imagesChannel.TopicName, err)
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
				ev := imagesChannel.Consumer.Poll(100)
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
				default:
					fmt.Printf("Ignored %v\n", e)
				}
			}
		}
	}()

}
