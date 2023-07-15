package users

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api.template.com.br/app/users/contracts"
	"api.template.com.br/helpers"
	"api.template.com.br/libraries"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (controller *UserController) ListAll(context *gin.Context) {
	log := libraries.GetLogger(context, nil)

	var err error

	model := UserModel{Logger: log}

	active := true

	var limit int
	limitQuery := context.Query("limit")
	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			helpers.GetHTTPResponse(http.StatusBadRequest, context, helpers.HTTPResponse{
				Result:  false,
				Message: "Limite de dados inválido",
			})
		}
	} else {
		limit = 10
	}

	var skip int
	skipQuery := context.Query("skip")
	if skipQuery != "" {
		skip, err = strconv.Atoi(skipQuery)
		if err != nil {
			helpers.GetHTTPResponse(http.StatusBadRequest, context, helpers.HTTPResponse{
				Result:  false,
				Message: "Skip de dados inválido",
			})
		}
	} else {
		skip = 0
	}

	if users, err := model.ListAll(int64(limit), int64(skip), &active); err != nil {
		helpers.GetHTTPResponse(http.StatusInternalServerError, context, helpers.HTTPResponse{
			Result:  false,
			Message: err.Error(),
		})
	} else {
		var userListResponse []contracts.User

		for _, user := range *users {
			userListResponse = append(userListResponse, contracts.User{
				Id:        user.Id.Hex(),
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Type:      user.Type,
				CreatedAt: user.CreatedAt,
			})
		}

		helpers.GetHTTPResponse(http.StatusOK, context, helpers.HTTPResponse{
			Result:  true,
			Message: "Listagem de usuários ativos",
			Data:    userListResponse,
		})
	}

	return

}

func (controller *UserController) Signup(context *gin.Context) {
	log := libraries.GetLogger(context, nil)

	model := UserModel{Logger: log}
	var userSignupRequest contracts.UserSignupRequest

	if err := context.BindJSON(&userSignupRequest); err != nil {
		log.
			WithField("error", err).
			WithField("requestHeaders", context.Request.Header).
			Error("Error on binding json")
		helpers.GetHTTPResponse(http.StatusBadRequest, context, helpers.HTTPResponse{
			Result:  false,
			Message: "Preencha todos os campos obrigatórios",
		})
		return
	}

	if exists, err := model.CheckIfUserExists(userSignupRequest.Email); err != nil {
		log.WithField("error", err).Error("Error on checking if user already exists")
		helpers.GetHTTPResponse(http.StatusServiceUnavailable, context, helpers.HTTPResponse{
			Result:  false,
			Message: "Ops...ocorreu um erro",
		})
		return
	} else {
		if exists {
			helpers.GetHTTPResponse(http.StatusConflict, context, helpers.HTTPResponse{
				Result:  false,
				Message: "Você não pode se cadastrar utilizando esse e-mail",
			})
			return
		} else {
			if _, err := model.Signup(userSignupRequest); err != nil {
				log.WithField("error", err).Error("Error on signup user")
				helpers.GetHTTPResponse(http.StatusInternalServerError, context, helpers.HTTPResponse{
					Result:  false,
					Message: "Ops...ocorreu um erro",
				})
				return
			} else {
				helpers.GetHTTPResponse(http.StatusOK, context, helpers.HTTPResponse{
					Result: true,
					Message: fmt.Sprintf(
						"Parabéns, %s! Sua conta foi criada com sucesso e enviamos uma mensagem para seu e-mail para que você confirme seu cadastro.",
						userSignupRequest.FirstName,
					),
				})

				return
			}
		}
	}

}

func (controller *UserController) ResendEmailVerification(context *gin.Context) {
	log := libraries.GetLogger(context, nil)

	model := UserModel{Logger: log}
	var emailVerificationRequest contracts.EmailVerificationRequest

	if err := context.BindJSON(&emailVerificationRequest); err != nil {
		log.
			WithField("error", err).
			WithField("requestHeaders", context.Request.Header).
			Error("Error on binding json")
		helpers.GetHTTPResponse(http.StatusBadRequest, context, helpers.HTTPResponse{
			Result:  false,
			Message: "Preencha todos os campos obrigatórios",
		})
		return
	}

	if ok, err := model.ResendEmailVerification(emailVerificationRequest.Email); err != nil {
		log.WithField("error", err).Error("Error on resending email user")
		helpers.GetHTTPResponse(http.StatusInternalServerError, context, helpers.HTTPResponse{
			Result:  false,
			Message: "Ops...an error occourred",
		})
		return
	} else if !ok {
		helpers.GetHTTPResponse(http.StatusConflict, context, helpers.HTTPResponse{
			Result:  false,
			Message: fmt.Sprintf("O e-mail de verificação não foi enviado para %s", emailVerificationRequest.Email),
			Data:    nil,
		})

		return
	} else {
		helpers.GetHTTPResponse(http.StatusOK, context, helpers.HTTPResponse{
			Result:  true,
			Message: fmt.Sprintf("Foi enviado um e-mail de verificação para %s", emailVerificationRequest.Email),
			Data:    nil,
		})

		return
	}

}

func (controller *UserController) Signin(context *gin.Context) {
	log := libraries.GetLogger(context, nil)

	model := UserModel{Logger: log}
	var userSigninRequest contracts.UserSigninRequest

	if err := context.BindJSON(&userSigninRequest); err != nil {
		log.
			WithField("error", err).
			WithField("requestHeaders", context.Request.Header).
			Error("Error on binding json")
		helpers.GetHTTPResponse(http.StatusBadRequest, context, helpers.HTTPResponse{
			Result:  false,
			Message: "Preencha todos os campos obrigatórios",
		})
		return
	}

	if user, err := model.Signin(userSigninRequest); err != nil {
		log.WithField("error", err).Error("Error on signin user")

		var statusCode int
		var errorMessage string

		switch err.Error() {
		case "invalid-credentials":
			statusCode = http.StatusUnauthorized
			errorMessage = "Credenciais invãlidas. Tente novamente"

		case "not-verified":
			statusCode = http.StatusUnauthorized
			errorMessage = "Sua conta ainda não foi verificada. Cheque seu e-mail e o verifique através da nossa mensagem."

		default:
			statusCode = http.StatusInternalServerError
			errorMessage = "Ops..um erro ocorreu no servidor"
		}

		helpers.GetHTTPResponse(statusCode, context, helpers.HTTPResponse{
			Result:  false,
			Message: errorMessage,
		})
		return
	} else {
		if accessToken, err := libraries.GenerateUserToken(user); err != nil {
			helpers.GetHTTPResponse(http.StatusInternalServerError, context, helpers.HTTPResponse{
				Result:  false,
				Message: err.Error(),
			})
		} else {
			helpers.GetHTTPResponse(http.StatusOK, context, helpers.HTTPResponse{
				Result:  true,
				Message: "Usuário logado",
				Data: contracts.UserSigninResponse{
					Id:          user.Id.Hex(),
					FirstName:   user.FirstName,
					LastName:    user.LastName,
					Email:       user.Email,
					Level:       user.Level,
					AccessToken: accessToken,
				},
			})
		}

		return
	}

}

func (controller *UserController) VerifyAccount(context *gin.Context) {
	log := libraries.GetLogger(context, nil)
	model := UserModel{Logger: log}

	var userVerifyRequest contracts.UserVerifyRequest

	if err := context.BindJSON(&userVerifyRequest); err != nil {
		log.
			WithField("error", err).
			WithField("requestHeaders", context.Request.Header).
			Error("Error on binding json")
		helpers.GetHTTPResponse(http.StatusBadRequest, context, helpers.HTTPResponse{
			Result:  false,
			Message: "Preencha todos os campos obrigatórios",
		})
		return
	}

	if verified, err := model.VerifyAccount(userVerifyRequest); err != nil {
		log.WithField("error", err).Error("Error on signin user")

		var statusCode int
		var errorMessage string

		if strings.Contains(err.Error(), "decrypt") {
			statusCode = http.StatusInternalServerError
			errorMessage = "Ops...ocorreu um erro no servidor"
		} else if strings.Contains(err.Error(), "find-user") {
			statusCode = http.StatusInternalServerError
			errorMessage = "Ops...ocorreu um erro no servidor"
		} else if strings.Contains(err.Error(), "verify-user") {
			statusCode = http.StatusInternalServerError
			errorMessage = "Ops...ocorreu um erro no servidor"
		} else {
			statusCode = http.StatusInternalServerError
			errorMessage = "Ops...ocorreu um erro no servidor"
		}

		helpers.GetHTTPResponse(statusCode, context, helpers.HTTPResponse{
			Result:  false,
			Message: errorMessage,
		})
		return
	} else {
		if verified {
			helpers.GetHTTPResponse(http.StatusOK, context, helpers.HTTPResponse{
				Result:  true,
				Message: "Conta de usuário verificada",
			})
			return
		} else {
			helpers.GetHTTPResponse(http.StatusConflict, context, helpers.HTTPResponse{
				Result:  false,
				Message: "Conta de usuário não verificada",
			})
			return
		}
	}

}
