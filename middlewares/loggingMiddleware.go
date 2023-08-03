package middlewares

import (
	"net/http/httputil"
	"time"

	"github.com/Chista-Framework/Chista/logger"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Logs each HTTP Request
func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Starting time request
		startTime := time.Now()

		// Processing request
		ctx.Next()

		// End Time request
		endTime := time.Now()

		// execution time
		latencyTime := endTime.Sub(startTime)

		// Request method
		reqMethod := ctx.Request.Method

		// Request route
		reqUri := ctx.Request.RequestURI

		// status code
		statusCode := ctx.Writer.Status()

		// Request IP
		clientIP := ctx.ClientIP()

		//Use global logger as middleware logger
		logger.Log.WithFields(log.Fields{
			"METHOD":    reqMethod,
			"URI":       reqUri,
			"STATUS":    statusCode,
			"LATENCY":   latencyTime,
			"CLIENT_IP": clientIP,
		}).Info("HTTP REQUEST")

		// If request dumping enabled, dump requests for debugging
		dump_mode_bool := logger.Dump_Mode
		if dump_mode_bool == true {
			requestDump, err := httputil.DumpRequest(ctx.Request, true)
			if err != nil {
				logger.Log.Debugln("Raw requst dump error:", err)
			}
			logger.Log.Debugln(string(requestDump))
		}

		ctx.Next()
	}
}
