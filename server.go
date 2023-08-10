package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/connection/", checkWebConnection)
	r.GET("/connection/loop", checkWebConnectionLoop)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func createLogger() *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)

	return log
}

func checkConnection(logger *logrus.Logger) *connectionResult {
	response, err := http.Get("https://example.com")
	if err != nil {
		logger.Debugf("http get request: %v", err)
		return &connectionResult{
			IsAlive: false,
			Message: "request error",
		}
	}
	defer response.Body.Close()

	// check status code
	if response.StatusCode != http.StatusOK {
		logger.WithFields(logrus.Fields{
			"status_code": response.StatusCode,
		}).Debug("http get request: status_code is not 200")
		return &connectionResult{
			IsAlive: false,
			Message: "example.com is not available",
		}
	}

	return &connectionResult{
		IsAlive: true,
		Message: "success",
	}
}

type connectionResult struct {
	IsAlive bool
	Message string
}

func checkWebConnection(c *gin.Context) {
	logger := createLogger()

	result := checkConnection(logger)
	if !result.IsAlive {
		c.JSON(500, gin.H{"message": result.Message})
		return
	}
	c.JSON(200, gin.H{"message": result.Message})
}

func checkWebConnectionLoop(c *gin.Context) {
	logger := createLogger()

	successCount := 0
	for {
		result := checkConnection(logger)
		if !result.IsAlive {
			c.JSON(500, gin.H{"message": result.Message})
			return
		}

		successCount += 1
		if successCount > 10 {
			c.JSON(200, gin.H{"message": result.Message})
			return
		}
	}
}
