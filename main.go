package main

import (
	"fmt"

	"api.template.com.br/helpers"
	"api.template.com.br/libraries"
	"github.com/gin-gonic/gin"
)

func main() {
	log := libraries.GetLogger(nil, nil)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	env, _ := helpers.GetEnvironment()

	InitRoutes(router)

	log.Debug(fmt.Sprintf("Running server on %s port...", env.PORT))

	router.Run(fmt.Sprintf(":%s", env.PORT))
}
