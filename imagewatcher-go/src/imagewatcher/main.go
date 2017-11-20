package main

import (
	"fmt"
	"net/http"
	"os"

	"imagewatcher/status"
)

func webHandlerRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}

var statusChannel *status.StatusChannel

func main() {
	statusChannel = status.NewStatusChannel()
	statusChannel.PostStatus("starting")

	http.HandleFunc("/", webHandlerRoot)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.ListenAndServe(":"+port, nil)
}
