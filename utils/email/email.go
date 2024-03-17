package email

import (
	"fmt"
	"strings"
	"time"

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

type SendEmailVerificationEmailData struct {
	LastName  string
	FirstName string
	Email     string
	Subject   string
	Token     string
	TokenTTL  time.Time
}

func SendEmailVerificationEmail(dataInput *SendEmailVerificationEmailData) {

	body := hermes.Body{
		Name: strings.Join([]string{dataInput.FirstName, dataInput.LastName}, " "),
		Intros: []string{
			"Welcome to Arravo Hackathon!",
			"Please verify your email to complete the registration process.",
		},
		Actions: []hermes.Action{
			{
				Instructions: "To verify your email, click the button below:",
				Button: hermes.Button{
					Color: "#22BC66",
					Text:  "Verify Email",
					Link:  "https://arravo.com/verify-email",
				},
			},
			{
				Instructions: "Alternatively, you can use the following token to verify your email:",
				InviteCode:   fmt.Sprintf("Token: %s", dataInput.Token),
			},
		},
		Outros: []string{
			"If you have any questions, feel free to contact us at support@arravo.co",
			"Thank you for joining Arravo Hackathon!",
		},
	}

	SendEmail(&SendEmailData{
		Email:   dataInput.Email,
		Message: &body,
		Subject: dataInput.Subject,
	})
}

type SendEmailVerificationCompleteEmailData struct {
	LastName  string
	FirstName string
	Email     string
	Subject   string
}

func SendEmailVerificationCompleteEmail(dataInput *SendEmailVerificationCompleteEmailData) {
	body := hermes.Body{
		Name: strings.Join([]string{}, " "),
		Intros: []string{
			"Welcome to Arravo Hackathon!",
			"Your email has been successfully verified.",
		},
		Outros: []string{
			"If you have any questions, feel free to contact us at support@arravo.co",
			"Thank you for joining Arravo Hackathon!",
		},
	}

	SendEmail(&SendEmailData{
		Email:   dataInput.Email,
		Message: &body,
		Subject: dataInput.Subject,
	})
}

type SendPasswordRecoveryEmailData struct {
	LastName  string
	FirstName string
	Email     string
	Subject   string
	Token     string
	TokenTTL  time.Time
}

func SendPasswordRecoveryEmail(dataInput *SendPasswordRecoveryEmailData) {

	body := hermes.Body{
		Name: strings.Join([]string{dataInput.FirstName, dataInput.LastName}, " "),
		Intros: []string{
			"Welcome to Arravo Hackathon!",
			"Please verify your email to complete the registration process.",
		},
		Actions: []hermes.Action{
			{
				Instructions: "To recover your password, click the button below:",
				Button: hermes.Button{
					Color: "#22BC66",
					Text:  "Verify Email",
					Link:  "https://arravo.com/verify-email",
				},
			},
			{
				Instructions: "Alternatively, you can use the following token to recover your password:",
				InviteCode:   fmt.Sprintf("Token: %s", dataInput.Token),
			},
		},
		Outros: []string{
			"If you have any questions, feel free to contact us at support@arravo.co",
			"Thank you for joining Arravo Hackathon!",
		},
	}

	SendEmail(&SendEmailData{
		Email:   dataInput.Email,
		Message: &body,
		Subject: dataInput.Subject,
	})
}

type SendPasswordRecoveryCompleteEmailData struct {
	LastName  string
	FirstName string
	Email     string
	Subject   string
}

func SendPasswordRecoveryCompleteEmail(dataInput *SendPasswordRecoveryCompleteEmailData) {
	body := hermes.Body{
		Name: strings.Join([]string{}, " "),
		Intros: []string{
			"Welcome to Arravo Hackathon!",
			"Your password has been successfully recovered.",
		},
		Outros: []string{
			"If you have any questions, feel free to contact us at support@arravo.co",
			"Thank you for joining Arravo Hackathon!",
		},
	}

	SendEmail(&SendEmailData{
		Email:   dataInput.Email,
		Message: &body,
		Subject: dataInput.Subject,
	})
}
