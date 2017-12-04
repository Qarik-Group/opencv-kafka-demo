package devicestreamrouting

import (
	"github.com/gin-gonic/gin"
)

// DemoDeviceStreamRouting represents a set of image URLs to static demo images
type DemoDeviceStreamRouting struct {
}

// NewDemoDeviceStreamRouting constructs a DemoDeviceStreamRouting
func NewDemoDeviceStreamRouting() (routing *DemoDeviceStreamRouting) {
	return &DemoDeviceStreamRouting{}
}

// HTMLLinks returns the information going into HTML template
func (routing *DemoDeviceStreamRouting) HTMLLinks() (links []gin.H) {
	return []gin.H{
		gin.H{
			"Name":                   "laptop",
			"FontAwesomeClass":       "fa-laptop",
			"RawImageURL":            "/demo/images/raw-drnic-laptop.png",
			"ObjectDetectorImageURL": "/demo/images/objectdetector-drnic-laptop.png",
			"EdgeDetectorImageURL":   "/demo/images/edgedetector-drnic-laptop.png",
		},
		gin.H{
			"Name":                   "raspberrypi",
			"FontAwesomeClass":       "fa-camera-retro",
			"RawImageURL":            "/demo/images/raw-drnic-pi.png",
			"ObjectDetectorImageURL": "/demo/images/objectdetector-drnic-pi.png",
			"EdgeDetectorImageURL":   "/demo/images/edgedetector-drnic-pi.png",
		},
	}
}

// RegisterGinRouting adds routes to Gin that are required for DemoDeviceStreams
func (routing *DemoDeviceStreamRouting) RegisterGinRouting(r *gin.Engine) {
	r.Static("/demo/images", "templates/demo/images")
}
