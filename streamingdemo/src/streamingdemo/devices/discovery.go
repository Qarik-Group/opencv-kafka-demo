package devices

import (
	"fmt"
	"os"
	"strings"
)

type Devices struct {
	Devices map[string]*Device
}

type Device struct {
}

func NewDiscoveredDevices() (devices *Devices) {
	var deviceIDs []string
	if os.Getenv("DEVICE_IDS") != "" {
		deviceIDs = strings.Split(os.Getenv("DEVICE_IDS"), ",")
		return NewLocalDevices(deviceIDs)
	}
	fmt.Fprintln(os.Stderr, "TODO: implement support for discovery in Cloud Foundry")
	os.Exit(1)
	return
}

func NewLocalDevices(deviceIDs []string) (devices *Devices) {
	devices = &Devices{}
	localTopicNamePrefix := "opencv-kafka-demo"
	streamTypes := []string{"raw", "objectdetection"}

	for _, deviceID := range deviceIDs {
		for _, streamType := range streamTypes {
			topicName := fmt.Sprintf("%s-%s-%s", localTopicNamePrefix, streamType, deviceID)
			// TODO add (streamType, deviceID) => NewDevice to structure
			fmt.Printf("Using Kafka topic '%s' for raw images from deviceID '%s'\n", topicName, deviceID)
		}
	}
	return
}

func (devices *Devices) PostImage(image *[]byte, streamType, deviceID string) {
	// TODO lookup (streamType, deviceID) => NewDevice
	var device *Device
	if device == nil {
		fmt.Printf("Unknown deviceID '%s' for streamType '%s'. Please configure Kafka topics.\n", deviceID, streamType)
		return
	}

	device.PostImage(image)
}
