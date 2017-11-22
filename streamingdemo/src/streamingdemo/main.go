package main

import (
	"os"
	"streamingdemo/devicestreamrouting"
	"streamingdemo/kafkastream"
	"streamingdemo/status"

	"net/http"

	"github.com/gin-gonic/gin"
)

var statusChannel *status.StatusChannel

func main() {
	statusChannel = status.NewStatusChannel()
	statusChannel.PostStatus("starting")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	var routing devicestreamrouting.DeviceStreamRouting
	if os.Getenv("DEMO_IMAGES") != "" {
		routing = devicestreamrouting.NewDemoDeviceStreamRouting()
	} else {
		streams := kafkastream.KafkaTopicStreams{}
		routing = devicestreamrouting.NewKafkaDeviceStreamRouting(streams)
	}
	routing.RegisterGinRouting(r)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"DeviceStreams": routing.HTMLLinks(),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:PORT
}
