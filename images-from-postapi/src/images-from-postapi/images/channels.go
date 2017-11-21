package images

import "fmt"

type ImageStreamChannels struct {
	Channels map[string]*ImagesChannel
}

func NewImageStreamChannels(deviceIDs []string) (channels *ImageStreamChannels) {
	channels = &ImageStreamChannels{
		Channels: map[string]*ImagesChannel{},
	}
	for _, deviceID := range deviceIDs {
		channels.Channels[deviceID] = NewImagesChannel(
			deviceID,
			"localhost:9092",
			fmt.Sprintf("opencv-kafka-demo-raw-%s", deviceID),
		)
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
