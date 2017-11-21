package devicestreamrouting

import (
	"fmt"
	"streamingdemo/kafkastream"

	"github.com/gin-gonic/gin"
)

// KafkaDeviceStreamRouting provides image URLs that will be streaming images from Kafka topics
type KafkaDeviceStreamRouting struct {
	Streams kafkastream.KafkaStreams
}

// NewKafkaDeviceStreamRouting constructs a KafkaDeviceStreamRouting
func NewKafkaDeviceStreamRouting(streams kafkastream.KafkaStreams) (routing *KafkaDeviceStreamRouting) {
	return &KafkaDeviceStreamRouting{Streams: streams}
}

// HTMLLinks returns the information going into HTML template
func (routing *KafkaDeviceStreamRouting) HTMLLinks() (links []gin.H) {
	link := gin.H{
		"Name":                    "drnic-laptop",
		"RawImageURL":             "/kafka/raw/drnic-laptop",
		"ObjectDetectionImageURL": "/kafka/objectdetection/drnic-laptop",
	}
	links = append(links, link)
	return
}

// RegisterGinRouting adds routes to Gin that are required for KafkaDeviceStreamRouting
func (routing *KafkaDeviceStreamRouting) RegisterGinRouting(r *gin.Engine) {
	r.GET("/kafka/:streamType/:deviceID", routing.ImageStream)
}

// ImageStream returns an HTTP image response to a GET /kafka/raw/drnic-laptop.jpg request
func (routing *KafkaDeviceStreamRouting) ImageStream(ctx *gin.Context) {
	streamType := ctx.Param("streamType")
	deviceID := ctx.Param("deviceID")
	fmt.Printf("Request GET jpg stream: streamType=%s deviceID=%s\n", streamType, deviceID)
}
