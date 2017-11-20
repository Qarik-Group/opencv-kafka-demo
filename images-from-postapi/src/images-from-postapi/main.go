package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"images-from-postapi/images"
	"images-from-postapi/status"
)

func webHandlerRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}

func webHandlerReceiveImage(w http.ResponseWriter, r *http.Request) {
	image, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to ReadAll image")
	}
	imagesChannel.PostImage(&image)
	fmt.Fprintf(w, "{}")
}

var statusChannel *status.StatusChannel
var imagesChannel *images.ImagesChannel

func main() {
	statusChannel = status.NewStatusChannel()
	statusChannel.PostStatus("starting")

	imagesChannel = images.NewImagesChannel()

	http.HandleFunc("/", webHandlerRoot)
	http.HandleFunc("/image", webHandlerReceiveImage)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.ListenAndServe(":"+port, nil)
}
