package eventsWeb

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
	"github.com/gin-gonic/gin"
)

func LimitHandler(lmt *config.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.MessageContentType, []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
