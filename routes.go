package main

import (
	"encoding/json"
	"net/http"
	"time"

	"api.template.com.br/app/users"
	"api.template.com.br/libraries"
	"github.com/gin-gonic/gin"
)

type logEntry struct {
	TimeStamp  string        `json:"timestamp"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency"`
	ClientIP   string        `json:"client_ip"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
}

func routerLogger(context *gin.Context) {
	log := libraries.GetLogger(context, nil)
	start := time.Now()

	context.Next()

	latency := time.Since(start)
	status := context.Writer.Status()

	entry := logEntry{
		TimeStamp:  start.Format(time.RFC3339),
		StatusCode: status,
		Latency:    latency,
		ClientIP:   context.ClientIP(),
		Method:     context.Request.Method,
		Path:       context.Request.URL.Path,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		log.Errorf("Could not marshal log entry: %v", err)
		return
	}

	log.Info(string(data))
}

func InitRoutes(router *gin.Engine) {

	router.Use(routerLogger)

	router.GET("/health-check", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"result": true,
		})
	})

	userRouter := users.UserRouter{}
	userRouter.InitRoutes(router)
}
