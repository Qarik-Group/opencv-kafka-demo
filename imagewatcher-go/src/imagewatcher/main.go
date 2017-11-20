package main

import (
	"imagewatcher/images"
	"imagewatcher/mjpeg"
	"imagewatcher/status"

	"net/http"

	"github.com/gin-gonic/gin"
)

// func webHandlerRoot(w http.ResponseWriter, r *http.Request) {
//   tmpl, err := template.New("index").Parse(...)
// // Error checking elided
// err = tmpl.Execute(w, data)
//
// }

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
	r.LoadHTMLGlob("templates/*")

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
