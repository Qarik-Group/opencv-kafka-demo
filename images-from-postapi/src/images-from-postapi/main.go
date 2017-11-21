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
	deviceID := r.FormValue("deviceID")
	if deviceID == "" {
		fmt.Fprintf(os.Stderr, "Missing 'deviceID' on /image")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Missing 'deviceID'"))
		return
	}
	image, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to ReadAll image")
	}
	imageStreamChannels.PostImage(&image, deviceID)
	fmt.Fprintf(w, "{}")
}

var statusChannel *status.StatusChannel
var imageStreamChannels *images.ImageStreamChannels

func main() {
	statusChannel = status.NewStatusChannel()
	statusChannel.PostStatus("starting")

	imageStreamChannels = images.NewImageStreamChannelsWithDiscovery()

	http.HandleFunc("/", webHandlerRoot)
	http.HandleFunc("/image", webHandlerReceiveImage)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.ListenAndServe(":"+port, nil)
}
