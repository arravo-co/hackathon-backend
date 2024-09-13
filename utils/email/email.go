package email

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

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
	TeamName string
}

type SendTeamMemberWelcomeEmailData struct {
	SendWelcomeEmailData
	TeamName     string
	TeamLeadName string
}
type SendIndividualWelcomeEmailData struct {
	SendWelcomeEmailData
}

type SendJudgeWelcomeEmailData struct {
	JudgeName string
	Email     string
	Subject   string
	TTL       int
	Token     string
	Link      string
}

type SendJudgeCreatedByAdminWelcomeEmailData struct {
	Name        string
	Email       string
	Subject     string
	TTL         int
	Token       string
	Link        string
	InviterName string
	Password    string
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
	AdminName   string
	InviterName string
	Email       string
	Subject     string
	TTL         int
	Token       string
	Password    string
	Link        string
}

type SendEmailVerificationEmailData struct {
	Name    string
	Email   string
	Subject string
	TTL     int
	Token   string
	Link    string
}

type SendEmailVerificationCompleteEmailData struct {
	Name    string
	Email   string
	Subject string
	TTL     int
	Token   string
	Link    string
}

type SendPasswordRecoveryEmailData struct {
	LastName  string
	FirstName string
	Email     string
	Subject   string
	Token     string
	Link      string
	TTL       uint32
}

type SendPasswordRecoveryCompleteEmailData struct {
	LastName  string
	FirstName string
	Email     string
	Subject   string
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

func SendTeamLeadWelcomeEmail(data *SendTeamLeadWelcomeEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/team_lead_welcome_and_verify_email.go.tmpl"))
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

func SendTeamMemberWelcomeEmail(data *SendTeamMemberWelcomeEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/team_member_welcome_email.go.tmpl"))
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

func SendJudgeWelcomeEmail(data *SendJudgeWelcomeEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/judge_welcome_and_verify_email.go.tmpl"))
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

func SendEmailVerificationEmail(dataInput *SendEmailVerificationEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/verify_email.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, dataInput)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	body := buf.String()
	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   dataInput.Email,
		Message: body,
		Subject: dataInput.Subject,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func SendEmailVerificationCompleteEmail(dataInput *SendEmailVerificationCompleteEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/verify_email_success.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, dataInput)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	body := buf.String()
	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   dataInput.Email,
		Message: body,
		Subject: dataInput.Subject,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func SendPasswordRecoveryEmail(data *SendPasswordRecoveryEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/password_recovery.go.tmpl"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	body := buf.String()
	//Subject: Password Recovery for Arravo Hackathon Account

	err = SendEmailHtml(&SendEmailHtmlData{
		Email:   data.Email,
		Message: body,
		Subject: data.Subject,
	})
	return err
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

func SendJudgeCreatedByAdminWelcomeEmail(data *SendJudgeCreatedByAdminWelcomeEmailData) error {
	tmpl := template.Must(template.ParseFiles("templates/judge_created_by_admin_welcome_and_verify_email.go.tmpl"))
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
