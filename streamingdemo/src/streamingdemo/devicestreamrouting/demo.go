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
			"Name":                    "drnic-laptop",
			"RawImageURL":             "/demo/images/raw-drnic-laptop.png",
			"ObjectDetectionImageURL": "/demo/images/objectdetector-drnic-laptop.png",
		},
	}
}

// RegisterGinRouting adds routes to Gin that are required for DemoDeviceStreams
func (routing *DemoDeviceStreamRouting) RegisterGinRouting(r *gin.Engine) {
	r.Static("/demo/images", "templates/demo/images")
}
