package main

import (
	"os"
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
	r.Static("/images", "templates/images")

	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		if os.Getenv("DEMO_IMAGES") != "" {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Devices": []gin.H{
					gin.H{
						"Name":                    "drnic-laptop",
						"RawImageURL":             "/images/raw-drnic-laptop.png",
						"ObjectDetectionImageURL": "/images/objectdetection-drnic-laptop.png",
					},
					gin.H{
						"Name":                    "drnic-pi",
						"RawImageURL":             "/images/raw-drnic-laptop.png",
						"ObjectDetectionImageURL": "/images/objectdetection-drnic-laptop.png",
					},
				},
			})
		} else {
			devices := make([]gin.H, 2)
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Devices": devices,
			})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:PORT
}
