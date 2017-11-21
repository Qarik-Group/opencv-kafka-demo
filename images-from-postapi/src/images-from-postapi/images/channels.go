package images

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
)

type ImageStreamChannels struct {
	Channels map[string]*ImagesChannel
}

func NewImageStreamChannelsWithDiscovery() (channels *ImageStreamChannels) {
	if os.Getenv("DEVICE_IDS") != "" {
		deviceIDs := strings.Split(os.Getenv("DEVICE_IDS"), ",")
		return NewImageStreamChannels(deviceIDs)
	}
	cfApp, err := cfenv.Current()
	if err != nil {
		fmt.Fprint(os.Stderr, "Missing $DEVICE_IDS and not running on Cloud Foundry")
		os.Exit(1)
	}

	channels = &ImageStreamChannels{
		Channels: map[string]*ImagesChannel{},
	}

	cfServiceNamePrefix := "raw-images-topic"
	imageTopicServices, err := cfApp.Services.WithNameUsingPattern(cfServiceNamePrefix + "-.*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot find service names with prefix '%s': %v", cfServiceNamePrefix, err)
		os.Exit(1)
	}
	for _, imageTopicService := range imageTopicServices {
		regex, err := regexp.Compile("(?i)^" + cfServiceNamePrefix + "-(.*)$")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create regexp.Compile: %v", err)
			os.Exit(1)
		}
		deviceID := regex.FindStringSubmatch(imageTopicService.Name)[1]

		hostname, _ := imageTopicService.CredentialString("hostname")
		topicName, _ := imageTopicService.CredentialString("topicName")
		channels.Channels[deviceID] = NewImagesChannel(
			deviceID,
			hostname,
			topicName,
		)
		fmt.Println("Discovered Kafka topic of raw images for deviceID", deviceID)
	}
	return
}

func NewImageStreamChannels(deviceIDs []string) (channels *ImageStreamChannels) {
	channels = &ImageStreamChannels{
		Channels: map[string]*ImagesChannel{},
	}
	for _, deviceID := range deviceIDs {
		topicName := fmt.Sprintf("opencv-kafka-demo-raw-%s", deviceID)
		channels.Channels[deviceID] = NewImagesChannel(
			deviceID,
			"localhost:9092",
			topicName,
		)
		fmt.Printf("Using Kafka topic '%s' for raw images from deviceID '%s'\n", topicName, deviceID)
	}
	return
}

func (channels *ImageStreamChannels) PostImage(image *[]byte, deviceID string) {
	if channels.Channels[deviceID] == nil {
		fmt.Printf("Unknown deviceID '%s'. Please configure Kafka topics.\n", deviceID)
		return
	}

	channels.Channels[deviceID].PostImage(image)
}
