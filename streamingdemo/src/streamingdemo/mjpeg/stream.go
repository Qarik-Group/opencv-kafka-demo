// Package mjpeg implements a simple MJPEG streamer.
//
// Evolved from https://github.com/saljam/mjpeg/blob/master/stream.go
//
// Stream objects implement the http.Handler interface, allowing to use them with the net/http package like so:
//	stream = mjpeg.NewStream()
//	http.Handle("/camera", stream)
// Then push new JPEG frames to the connected clients using stream.UpdateJPEG().
package mjpeg

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Stream represents a single video feed.
type Stream struct {
	requestChannels map[chan []byte]bool
	reqChannelsLock sync.Mutex
	frameBuffer     []byte
	FrameInterval   time.Duration
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" +
	"--" + boundaryWord + "\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

// ServeGinContent begins sending frames as a moving image (M-JPEG)
func (s *Stream) ServeGinContent(ctx *gin.Context) {
	log.Println("Stream: connected")
	ctx.Header("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)

	requestChannel := make(chan []byte)
	s.reqChannelsLock.Lock()
	s.requestChannels[requestChannel] = true
	s.reqChannelsLock.Unlock()

	for {
		time.Sleep(s.FrameInterval)
		frameBytes := <-requestChannel
		_, err := ctx.Writer.Write(frameBytes)
		if err != nil {
			break
		}
	}

	s.reqChannelsLock.Lock()
	delete(s.requestChannels, requestChannel)
	s.reqChannelsLock.Unlock()
	log.Println("Stream: disconnected")
}

// UpdateJPEG pushes a new JPEG frame onto the clients.
func (s *Stream) UpdateJPEG(jpeg []byte) {
	header := fmt.Sprintf(headerf, len(jpeg))
	if len(s.frameBuffer) < len(jpeg)+len(header) {
		s.frameBuffer = make([]byte, (len(jpeg)+len(header))*2)
	}

	copy(s.frameBuffer, header)
	copy(s.frameBuffer[len(header):], jpeg)

	s.reqChannelsLock.Lock()
	for requestChannel := range s.requestChannels {
		// Select to skip streams which are sleeping to drop frames.
		// This might need more thought.
		select {
		case requestChannel <- s.frameBuffer:
		default:
		}
	}
	s.reqChannelsLock.Unlock()
}

// NewStream initializes and returns a new Stream.
func NewStream() *Stream {
	return &Stream{
		requestChannels: make(map[chan []byte]bool),
		frameBuffer:     make([]byte, len(headerf)),
		FrameInterval:   50 * time.Millisecond,
	}
}
