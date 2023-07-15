package services

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"api.template.com.br/helpers"
	"api.template.com.br/libraries"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendEmailPayload struct {
	RecipientName  string
	RecipientEmail string
	Subject        string
	TemplateData   any
	TemplateType   string // welcome, forgot-password, message, event, privacy-policy
}

func SendEmail(payload SendEmailPayload) {
	log := libraries.GetLogger(nil, nil)
	environment, _ := helpers.GetEnvironment()

	from := mail.NewEmail(environment.MAIL_USER_NAME, environment.MAIL_USER_EMAIL)
	subject := payload.Subject
	to := mail.NewEmail(payload.RecipientName, payload.RecipientEmail)

	// creating template by type
	// Specify the path to the HTML file
	filePath := fmt.Sprintf("./templates/%s.html", payload.TemplateType)

	// Read the file content
	contentFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	htmlTemplate, err := template.New(payload.TemplateType).Parse(string(contentFile))
	if err != nil {
		log.Error("error getting template")
		return
	}

	var templateOutput bytes.Buffer
	if err := htmlTemplate.ExecuteTemplate(&templateOutput, payload.TemplateType, payload.TemplateData); err != nil {
		log.Error("error on template generating")
		return
	}

	message := mail.NewSingleEmail(from, subject, to, "", templateOutput.String())
	client := sendgrid.NewSendClient(environment.SENDGRID_API_KEY)
	response, err := client.Send(message)
	if err != nil {
		log.WithField("payload", payload).Error("error on sending email")
		return
	} else {
		log.WithField("payload", payload).WithField("response", response).Info("email notification send successfully")
		return
	}
}
