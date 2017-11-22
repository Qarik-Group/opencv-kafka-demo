package kafkastream

import (
	"fmt"
	"os"
	"streamingdemo/mjpeg"
	"strings"
)

// KafkaTopicStreams is the collection of KafkaTopics being watched for new frames
type KafkaTopicStreams struct {
	StreamTypes  []string
	DeviceIDs    []string
	MJPEGStreams map[string]*mjpeg.Stream
}

// NewKafkaTopicStreams constructs a NewKafkaTopicStreams using
// environment variables to discover accessible kafka topics for image streams
func NewKafkaTopicStreams() (streams *KafkaTopicStreams) {
	streams = &KafkaTopicStreams{
		StreamTypes:  []string{"raw", "objectdetection"},
		MJPEGStreams: map[string]*mjpeg.Stream{},
	}
	streams.discoverDeviceIDs()
	return
}

// DiscoverDeviceIDs looks at environment variables to discover which Kafka topics
// the application should be watching, and which mjpeg.Stream to construct for viewing
func (streams *KafkaTopicStreams) discoverDeviceIDs() {
	if os.Getenv("DEVICE_IDS") != "" {
		streams.DeviceIDs = strings.Split(os.Getenv("DEVICE_IDS"), ",")
		streams.setupLocalKafka()
		return
	}
	fmt.Fprintf(os.Stderr, "No support for Cloud Foundry yet. Must set $DEVICE_IDS.")
	os.Exit(1)
}

// setupLocalKafka assumes the topic name for each streamType/deviceID with local kafka
func (streams *KafkaTopicStreams) setupLocalKafka() {
	for _, streamType := range streams.StreamTypes {
		for _, deviceID := range streams.DeviceIDs {
			kafkaTopic := streams.localKafkaTopic(streamType, deviceID)
			mjpegStreamKey := streams.mjpegStreamKey(streamType, deviceID)

			fmt.Println("Watching local Kafka topic", kafkaTopic)
			streams.MJPEGStreams[mjpegStreamKey] = mjpeg.NewStream()
			// TODO Watching local Kafka topic -> mjpeg.Stream
		}
	}
}

func (streams *KafkaTopicStreams) localKafkaTopic(streamType, deviceID string) string {
	return fmt.Sprintf("opencv-kafka-demo-%s-%s", streamType, deviceID)
}

func (streams *KafkaTopicStreams) mjpegStreamKey(streamType, deviceID string) string {
	return fmt.Sprintf("%s/%s", streamType, deviceID)
}
