package libraries

import (
	b "bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"fmt"
	"net/http"

	"api.template.com.br/app/types"
	"api.template.com.br/helpers"
)

func AuthMiddleware(context *gin.Context) {
	environment, _ := helpers.GetEnvironment()
	accessToken := context.GetHeader("access_token")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
		}

		return []byte(environment.JWT_SECRET), nil
	})

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// context.Set("id", claims["id"])
		// context.Set("name", claims["name"])
		// context.Set("type", claims["type"])
		// context.Set("email", claims["email"])

		context.Next()
	} else {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"result":  false,
			"message": "Sem autorização",
		})
		return
	}
}

func OwnerAuthMiddleware(context *gin.Context) {
	environment, _ := helpers.GetEnvironment()
	accessToken := context.GetHeader("access_token")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
		}

		return []byte(environment.JWT_SECRET), nil
	})

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var authUser struct {
			Id string `json:"userId"`
		}

		// getting body's bytes to reuse the ReadCloser
		body, bodyErr := io.ReadAll(context.Request.Body)

		if err := json.NewDecoder(io.NopCloser(b.NewReader(body))).Decode(&authUser); authUser.Id == claims["id"] && err == nil && bodyErr == nil {
			context.Set("id", claims["id"])
			context.Set("level", claims["level"])
			context.Set("firstName", claims["firstName"])
			context.Set("lastName", claims["lastName"])
			context.Set("email", claims["email"])

			// create a new Reader for binding into controller functions
			context.Request.Body = io.NopCloser(b.NewReader(body))

			context.Next()
			return
		}

	}

	context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"result":  false,
		"message": "No authorization",
	})
}

func GenerateUserToken(user *types.UserDTO) (string, error) {
	if user != nil {
		environment, _ := helpers.GetEnvironment()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":        user.Id.Hex(),
			"level":     user.Level,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
		})

		tokenString, err := token.SignedString([]byte(environment.JWT_SECRET))

		return tokenString, err
	} else {
		return "", fmt.Errorf("user null")
	}
}
