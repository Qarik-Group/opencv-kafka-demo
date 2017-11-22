package kafkastream

import (
	"fmt"
	"os"
	"regexp"
	"streamingdemo/mjpeg"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

// KafkaTopicStreams is the collection of KafkaTopics being watched for new frames
type KafkaTopicStreams struct {
	StreamTypes  []string
	DeviceIDs    []string
	ImageTopics  map[string]*ImageTopic
	MJPEGStreams map[string]*mjpeg.Stream
}

// NewKafkaTopicStreams constructs a NewKafkaTopicStreams using
// environment variables to discover accessible kafka topics for image streams
func NewKafkaTopicStreams() (streams *KafkaTopicStreams) {
	streams = &KafkaTopicStreams{
		StreamTypes:  []string{"raw", "objectdetector"},
		ImageTopics:  map[string]*ImageTopic{},
		MJPEGStreams: map[string]*mjpeg.Stream{},
	}
	streams.discoverDeviceIDs()
	return
}

// DiscoverDeviceIDs looks at environment variables to discover which Kafka topics
// the application should be watching, and which mjpeg.Stream to construct for viewing
func (streams *KafkaTopicStreams) discoverDeviceIDs() {
	if cfenv.IsRunningOnCF() {
		cfApp, err := cfenv.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR looking up CF services: %s", err)
			os.Exit(1)
		}

		rawServices, err := cfApp.Services.WithNameUsingPattern("raw-images-topic-.*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "At least one topic named 'raw-images-topic-<device-id>' must be bound to app.")
			os.Exit(1)
		}
		streams.setupCFKafka(rawServices)
		return
	}
	if os.Getenv("DEVICE_IDS") != "" {
		streams.DeviceIDs = strings.Split(os.Getenv("DEVICE_IDS"), ",")
		streams.setupLocalKafka()
		return
	}
	fmt.Fprintf(os.Stderr, "Must set $DEVICE_IDS if not running within Cloud Foundry.")
	os.Exit(1)
}

// setupCFKafka uses Cloud Foundry services to access Kafka topics for each image stream
func (streams *KafkaTopicStreams) setupCFKafka(rawServices []cfenv.Service) {
	regex := regexp.MustCompile("(?i)^raw-images-topic-(.*)$")

	deviceIDs := make([]string, len(rawServices))
	for i, rawService := range rawServices {
		deviceID := regex.FindStringSubmatch(rawService.Name)[1]
		deviceIDs[i] = deviceID
	}
	streams.DeviceIDs = deviceIDs

	cfApp, _ := cfenv.Current()

	for _, streamType := range streams.StreamTypes {
		for _, deviceID := range streams.DeviceIDs {
			serviceName := fmt.Sprintf("%s-images-topic-%s", streamType, deviceID)
			service, err := cfApp.Services.WithName(serviceName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to find service '%s': %v", serviceName, err)
				os.Exit(1)
			}
			hostname, _ := service.CredentialString("hostname")
			kafkaTopic, _ := service.CredentialString("topicName")

			mjpegStreamKey := streams.mjpegStreamKey(streamType, deviceID)
			mjpegStream := mjpeg.NewStream()
			streams.MJPEGStreams[mjpegStreamKey] = mjpegStream

			streams.ImageTopics[mjpegStreamKey] = NewImageTopic(hostname, kafkaTopic, mjpegStream)

			fmt.Printf("Watching Cloud Foundry Kafka service %s for stream %s\n", serviceName, mjpegStreamKey)
		}
	}
}

// setupLocalKafka assumes the topic name for each streamType/deviceID with local kafka
func (streams *KafkaTopicStreams) setupLocalKafka() {
	for _, streamType := range streams.StreamTypes {
		for _, deviceID := range streams.DeviceIDs {
			kafkaTopic := streams.localKafkaTopic(streamType, deviceID)
			mjpegStreamKey := streams.mjpegStreamKey(streamType, deviceID)

			fmt.Println("Watching local Kafka topic", kafkaTopic)
			mjpegStream := mjpeg.NewStream()
			streams.MJPEGStreams[mjpegStreamKey] = mjpegStream
			streams.ImageTopics[mjpegStreamKey] = NewImageTopic("localhost:9092", kafkaTopic, mjpegStream)
		}
	}
}

func (streams *KafkaTopicStreams) localKafkaTopic(streamType, deviceID string) string {
	return fmt.Sprintf("opencv-kafka-demo-%s-%s", streamType, deviceID)
}

func (streams *KafkaTopicStreams) mjpegStreamKey(streamType, deviceID string) string {
	return fmt.Sprintf("%s/%s", streamType, deviceID)
}
