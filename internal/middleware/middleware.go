package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001", 
			"https://frontend-tormenta20.vercel.app",
			"https://backend-tormenta20.fly.dev",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true, // Habilitar credentials com origens específicas
		MaxAge:           12 * time.Hour,
	})
}

func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func ErrorHandler() gin.HandlerFunc {
	return gin.RecoveryWithWriter(os.Stdout, func(c *gin.Context, recovered interface{}) {
		if neterr, ok := recovered.(*net.OpError); ok {
			if neterr.Timeout() {
				c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
				c.Abort()
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		c.Abort()
	})
}
