package users

import (
	"encoding/base64"
	"fmt"
	"strings"

	"api.template.com.br/helpers"
	"github.com/sirupsen/logrus"

	"api.template.com.br/app/types"
	"api.template.com.br/app/users/contracts"
	"api.template.com.br/libraries"
	"api.template.com.br/services"
)

type UserModel struct {
	Logger *logrus.Logger
}

func (model *UserModel) CheckIfUserExists(email string) (bool, error) {
	log := model.Logger
	log.WithField("modelLayer", map[string]any{"func": "CheckIfUserExists", "email": email})

	repository := UserRepository{Logger: log}
	if user, err := repository.FindOneByEmail(email); err != nil {
		return false, err
	} else {
		return user != nil, nil
	}
}

func (model *UserModel) ListAll(limit, skip int64, active *bool) (*[]types.UserDTO, error) {
	log := model.Logger
	log.WithField("modelLayer", map[string]any{"func": "ListAll", "active": active})

	repository := UserRepository{Logger: log}
	if result, err := repository.FindAll(limit, skip, active); err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func (model *UserModel) Signup(payload contracts.UserSignupRequest) (*types.UserDTO, error) {
	log := model.Logger
	log.WithField("modelLayer", map[string]any{"func": "Signup", "signupEmail": payload.Email})

	encryptedPass, err := libraries.Encrypt(payload.Password)
	if err != nil {
		return nil, err
	}

	user := types.UserDTO{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Type:      "common",
		Password:  encryptedPass,
		RepositoryGeneralProps: types.RepositoryGeneralProps{
			Active: true,
		},
		Verified: false,
	}

	repository := UserRepository{Logger: log}
	if result, err := repository.Create(user); err != nil {
		return nil, err
	} else {
		user.Id = *result

		sendEmailVerification(user.Email, user.FirstName)

		return &user, nil
	}
}

func (model *UserModel) ResendEmailVerification(email string) (bool, error) {
	log := model.Logger
	log.WithField("modelLayer", map[string]any{"func": "ResendEmailVerification", "args": map[string]any{"email": email}})

	repository := UserRepository{Logger: log}
	if user, err := repository.FindOneByEmail(email); err != nil {
		return false, fmt.Errorf("user not found")
	} else {
		if ok, err := sendEmailVerification(user.Email, user.FirstName); err != nil {
			log.WithField("error", err).Errorf("error on sending email verification: %s", err.Error())
			return false, fmt.Errorf("email verification not sent")
		} else {
			return ok, nil
		}
	}
}

func sendEmailVerification(email, name string) (bool, error) {
	environment, _ := helpers.GetEnvironment()

	if verificationHash, err := libraries.Encrypt(email); err != nil {
		return false, fmt.Errorf("error on generating hash to verify account")
	} else {
		b64UrlSafe := base64.RawURLEncoding.EncodeToString([]byte(verificationHash))
		templateData := map[string]string{
			"Name":           name,
			"VerficationURL": fmt.Sprintf("%s/verification?hash=%s", environment.WEBSITE_URL, b64UrlSafe),
		}

		go services.SendEmail(services.SendEmailPayload{
			RecipientEmail: email,
			RecipientName:  name,
			TemplateData:   templateData,
			Subject:        "Welcome to the Template!",
			TemplateType:   "welcome",
		})

		return true, nil
	}
}

func (model *UserModel) VerifyAccount(payload contracts.UserVerifyRequest) (bool, error) {
	log := model.Logger
	log.WithField("modelLayer", map[string]any{"func": "VerifyAccount", "args": map[string]any{"payload": payload}})

	repository := UserRepository{Logger: log}

	decryptedHash, err := libraries.Decrypt(strings.Trim(payload.VerifyCode, " "))
	if err != nil {
		log.WithField("error", err).Errorf("error on decrypting hash: %s", err.Error())
		return false, fmt.Errorf("decrypt: %s", err.Error())
	}

	user, err := repository.FindOneByEmail(decryptedHash)
	if err != nil {
		log.WithField("error", err).Errorf("error on decrypting hash: %s", err.Error())
		return false, fmt.Errorf("find-user: %s", err.Error())
	}

	verified, err := repository.VerifyAccount(user.Id.Hex())
	if err != nil {
		log.WithField("error", err).Errorf("error on verifying user: %s", err.Error())
		return false, fmt.Errorf("verify-user: %s", err.Error())
	}

	return verified, nil
}

func (model *UserModel) Signin(payload contracts.UserSigninRequest) (*types.UserDTO, error) {
	log := model.Logger
	repository := UserRepository{Logger: log}
	user, err := repository.FindOneByEmail(payload.Email)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, fmt.Errorf("invalid-credentials")
	}

	decryptedPass, err := libraries.Decrypt(user.Password)
	if err != nil {
		return nil, err
	} else if decryptedPass != payload.Password {
		return nil, fmt.Errorf("invalid-credentials")
	}

	if !user.Verified {
		return nil, fmt.Errorf("not-verified")
	}

	return user, nil
}
