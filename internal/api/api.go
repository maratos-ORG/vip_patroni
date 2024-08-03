package api

import (
	"time"
	"vip_patroni/internal/config"

	"vip_patroni/internal/ipmanager"
	log "vip_patroni/internal/logging"

	"github.com/gin-gonic/gin"
)

func NewAPI(conf *config.Config, ma *ipmanager.IPManager) {

	gin.SetMode(gin.ReleaseMode)
	webAPI := gin.New()
	webAPI.Use(Logger())
	webAPI.Use(gin.Recovery())
	webAPI.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"state":   Bool2int(ma.Configurer.QueryAddress()),
			"desired": Bool2int(ma.CurrentState),
		})
	})

	webAPI.GET("/config", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"config": conf,
		})
	})
	err := webAPI.Run(":" + conf.APIPort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		log.Fatal("unable to run webapi: %s", err)
	}
}

func Logger() gin.HandlerFunc {

	// Set the log format
	return func(c *gin.Context) {
		// Starting time
		startTime := time.Now().UnixMilli()
		// Processing request
		c.Next()
		// End Time
		endTime := time.Now().UnixMilli()
		// execution time
		latencyTime := endTime - startTime
		// Request method
		reqMethod := c.Request.Method
		// Request route
		reqURI := c.Request.RequestURI
		// status code
		statusCode := c.Writer.Status()
		// Request IP
		clientIP := c.ClientIP()

		log.WithFields(map[string]interface{}{
			"caller":     clientIP,
			"status":     statusCode,
			"latency ms": latencyTime,
			"method":     reqMethod,
			"url":        reqURI}).Info("API call")
	}
}

func Bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
