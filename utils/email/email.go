package email

import (
	"fmt"
	"strings"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/matcornic/hermes/v2"
	"github.com/resend/resend-go/v2"
)

type SendEmailData struct {
	Email   string
	Subject string
	Message *hermes.Body
}

func SendEmail(data *SendEmailData) error {
	h := hermes.Hermes{}
	email := hermes.Email{
		Body: *data.Message,
	}
	em, _ := h.GenerateHTML(email)
	apiKey := config.GetResendAPIKey()

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    config.GetResendFromEmail(),
		To:      []string{data.Email},
		Subject: data.Subject,
		Html:    em,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return err
	}
	fmt.Printf("%#v", sent)
	return nil
}

type SendWelcomeEmailData struct {
	LastName  string
	FirstName string
	Email     string
	Subject   string
}

func SendWelcomeEmail(data *SendWelcomeEmailData) {
	body := &hermes.Body{
		Name: strings.Join([]string{data.LastName, data.FirstName}, " "),
	}
	SendEmail(&SendEmailData{
		Email:   data.Email,
		Message: body,
		Subject: data.Subject,
	})
}

func SendEmailVerificationEmail(data *SendWelcomeEmailData) {
	body := &hermes.Body{
		Name: strings.Join([]string{data.LastName, data.FirstName}, " "),
	}
	SendEmail(&SendEmailData{
		Email:   data.Email,
		Message: body,
		Subject: data.Subject,
	})
}
