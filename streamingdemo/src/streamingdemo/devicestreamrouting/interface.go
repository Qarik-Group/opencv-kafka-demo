package devicestreamrouting

import (
	"github.com/gin-gonic/gin"
)

// DeviceStreamRouting describes how an HTML template will be rendered to show live streams
type DeviceStreamRouting interface {
	HTMLLinks() (links []gin.H)
	RegisterGinRouting(r *gin.Engine)
}
