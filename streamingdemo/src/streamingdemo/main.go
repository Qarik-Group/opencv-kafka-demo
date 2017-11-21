package main

import (
	"streamingdemo/images"
	"streamingdemo/mjpeg"
	"streamingdemo/status"

	"net/http"

	"github.com/gin-gonic/gin"
)

var statusChannel *status.StatusChannel
var imagesChannel *images.ImagesChannel
var imageStream *mjpeg.Stream

func main() {
	statusChannel = status.NewStatusChannel()
	statusChannel.PostStatus("starting")

	imagesChannel = images.NewImagesChannel()

	imageStream = mjpeg.NewStream()

	go func() {
		imagesChannel.WatchImages(imageStream.UpdateJPEG)
	}()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/images", "templates/images")

	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"videoFeedURL": "/stream_image",
		})
	})
	r.GET("/stream_image", imageStream.ServeGinContent)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
