package libraries

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogFields map[string]interface{}

func GetLogger(routerContext *gin.Context, fields *LogFields) *logrus.Logger {

	level, levelErr := logrus.ParseLevel("trace")
	if levelErr != nil {
		fmt.Println("Incorrect logger level", levelErr)
	}

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if routerContext != nil {
		logrus.WithContext(routerContext)
	}

	if fields != nil {
		logrus.WithFields(logrus.Fields{
			"Context": fields,
		})
	}

	return logrus.StandardLogger()
}
