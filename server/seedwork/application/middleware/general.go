package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware that logs the request details
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// After request
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logMsg := fmt.Sprintf("[GIN] %s | %s | %s | %s | Status: %d | %s | %s\n",
			time.Now().Format(time.RFC3339),
			method,
			path,
			c.ClientIP(),
			statusCode,
			latency.String(),
			c.GetString("error"))

		gin.DefaultWriter.Write([]byte(logMsg))

		// Log errors if any
		if len(c.Errors) > 0 {
			gin.DefaultErrorWriter.Write([]byte(c.Errors.String()))
		}

		// If we have a slow request, log it differently
		if latency > time.Second*5 {
			gin.DefaultWriter.Write([]byte("SLOW REQUEST: " + path + " took " + latency.String() + "\n"))
		}
	}
}

// CORS middleware to handle Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandler is a middleware that handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			c.JSON(c.Writer.Status(), gin.H{
				"errors": c.Errors.Errors(),
			})
		}
	}
}
