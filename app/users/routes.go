package users

import (
	"fmt"

	"api.template.com.br/libraries"

	"github.com/gin-gonic/gin"
)

const (
	BASEPATH = "/users"
)

type UserRouter struct {
}

func (userRoutes *UserRouter) InitRoutes(router *gin.Engine) {

	var controller *UserController

	router.GET(fmt.Sprintf("%s/", BASEPATH), libraries.AuthMiddleware, controller.ListAll)

	router.POST(fmt.Sprintf("%s/signup", BASEPATH), controller.Signup)
	router.POST(fmt.Sprintf("%s/signin", BASEPATH), controller.Signin)
	router.POST(fmt.Sprintf("%s/verify", BASEPATH), controller.VerifyAccount)
	router.POST(fmt.Sprintf("%s/resendverification", BASEPATH), controller.ResendEmailVerification)
}
