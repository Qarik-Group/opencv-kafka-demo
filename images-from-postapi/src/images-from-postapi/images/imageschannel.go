package images

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ImagesChannel struct {
	Hostname  string
	TopicName string
	Producer  *kafka.Producer
}

// NewImagesChannel constructs ImagesChannel and sets up connection
func NewImagesChannel() (imagesChannel *ImagesChannel) {
	imagesChannel = &ImagesChannel{
		Hostname:  "localhost:9092",
		TopicName: "opencv-kafka-demo-raw-images",
	}

	cfApp, err := cfenv.Current()
	if err == nil {
		imagesTopicService, err := cfApp.Services.WithName("raw-images-topic")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot find service name 'images-topic': %v", err)
			os.Exit(1)
		}
		imagesChannel.Hostname, _ = imagesTopicService.CredentialString("hostname")
		imagesChannel.TopicName, _ = imagesTopicService.CredentialString("topicName")
	} else {
		fmt.Fprintf(os.Stderr, "Not running inside Cloud Foundry. Assume local Kafka on localhost:9092\n")
	}

	imagesChannel.Producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": imagesChannel.Hostname})
	if err != nil {
		fmt.Printf("Failed to create 'images-topic' producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created '%s' producer %v\n", imagesChannel.TopicName, imagesChannel.Producer)
	return
}

// PostImage is an all-in-one helper to post an image to kafka topic
func (imagesChannel *ImagesChannel) PostImage(image *[]byte) {
	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)

	err := imagesChannel.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &imagesChannel.TopicName,
			Partition: kafka.PartitionAny,
		},
		Value: *image,
	}, deliveryChan)
	if err != nil {
		fmt.Printf("Failed .Produce for'images-topic' message: %s\n", err)
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
