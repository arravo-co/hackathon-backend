package email

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/matcornic/hermes/v2"
	"github.com/resend/resend-go/v2"
)

type SendEmailData struct {
	Email   string
	Subject string
	Message *hermes.Body
}

type SendEmailHtmlData struct {
	Email   string
	Subject string
	Message string
}

func SendEmailHtml(data *SendEmailHtmlData) error {
	apiKey := config.GetResendAPIKey()

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    config.GetResendFromEmail(),
		To:      []string{data.Email},
		Subject: data.Subject,
		Html:    data.Message,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return err
	}
	fmt.Printf("%#v", sent)
	return nil
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
	TTL       int
	Token     string
	Link      string
}

type SendTeamLeadWelcomeEmailData struct {
	SendWelcomeEmailData
}
type SendIndividualWelcomeEmailData struct {
	SendWelcomeEmailData
}

type SendJudgeWelcomeEmailData struct {
	JudgeName   string
	Email       string
	Subject     string
	TTL         int
	Token       string
	Link        string
	InviterName string
}

type SendTeamInviteEmailData struct {
	InviterName  string
	InviteeName  string
	InviterEmail string
	InviteeEmail string
	Subject      string
	TTL          int
	Link         string
}

type SendAdminWelcomeEmailData struct {
	FirstName string
	LastName  string
	Email     string
	Subject   string
	TTL       int
	Token     string
	Link      string
}

type SendAdminCreatedByAdminWelcomeEmailData struct {
	FirstName string
	LastName  string
	Email     string
	Subject   string
	TTL       int
	Token     string
	Password  string
}

func SendTeamLeadWelcomeEmail(data *SendTeamLeadWelcomeEmailData) {
	tmpl := template.Must(template.ParseFiles("templates/welcome_and_verify_email.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		exports.MySugarLogger.Error(err)
	}
	body := buf.String()
	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   data.Email,
		Message: body,
		Subject: data.Subject,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
	}
}

func SendIndividualParticipantWelcomeEmail(data *SendIndividualWelcomeEmailData) {
	tmpl := template.Must(template.ParseFiles("templates/welcome_and_verify_email.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		exports.MySugarLogger.Error(err)
	}
	body := buf.String()
	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   data.Email,
		Message: body,
		Subject: data.Subject,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
	}
}

func SendJudgeWelcomeEmail(data *SendJudgeWelcomeEmailData) {
	tmpl := template.Must(template.ParseFiles("templates/welcome_and_verify_email.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		exports.MySugarLogger.Error(err)
	}
	body := buf.String()
	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   data.Email,
		Message: body,
		Subject: data.Subject,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
	}
}

func SendWelcomeEmail(data *SendIndividualWelcomeEmailData) {
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
	Link      string
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

func SendInviteTeamMemberEmail(data *SendTeamInviteEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/email_invite.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	body := buf.String()
	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   data.InviteeEmail,
		Message: body,
		Subject: data.Subject,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func SendAdminWelcomeEmail(data *SendAdminWelcomeEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/admin_welcome_and_verify_email.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	body := buf.String()
	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   data.Email,
		Message: body,
		Subject: data.Subject,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func SendAdminCreatedByAdminWelcomeEmail(data *SendAdminCreatedByAdminWelcomeEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/admin_created_by_admin_welcome_and_verify_email.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	body := buf.String()
	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   data.Email,
		Message: body,
		Subject: data.Subject,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}
