package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/matcornic/hermes/v2"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

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
	port := config.MustGetSMTPPort()
	host := config.GetSMTPHost()
	user_name := config.GetSMTPUsername()
	password := config.GetSMTPPassword()
	fmt.Println(port, host, user_name, password)

	//smtpAuth := smtp.PlainAuth("", user_name, password, host)
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			log.Printf("Error closing client: %v", cerr)
		}
	}()

	// Start TLS (upgrade the connection to TLS)
	if err = client.StartTLS(&tls.Config{
		ServerName: host,
	}); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate
	if err = client.Auth(LoginAuth(user_name, password)); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	fmt.Println("Here")
	if err = client.Mail(user_name); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err = client.Rcpt(data.Email); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// Construct the email content with headers
	emailContent := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\n%s\n%s",
		user_name, data.Email, data.Subject, mime, data.Message,
	)
	_, err = writer.Write([]byte(emailContent))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return client.Quit()

	/*
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
	*/
}

func SendEmail(data *SendEmailData) error {
	/*
		h := hermes.Hermes{}
		email := hermes.Email{
			Body: *data.Message,
		}
	*/
	port := config.MustGetSMTPPort()
	host := config.GetSMTPHost()
	user_name := config.GetSMTPUsername()
	password := config.GetSMTPPassword()
	fmt.Println(port, host, user_name, password)
	/*
	 */
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	/*err = client.Hello("my_localhost")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}*/
	fmt.Println("Again Again")
	err = client.StartTLS(&tls.Config{
		//InsecureSkipVerify: true,
		ServerName: host,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	smtpAuth := smtp.PlainAuth("", user_name, password, host)
	fmt.Println("Again Hereoooooooooooooooooooooooooooooooooooooooooooooooooooooo")
	err = client.Auth(smtpAuth)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Here")
	err = client.Mail(user_name)

	if err != nil {
		return err
	}

	err = client.Rcpt(data.Email)
	if err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	str := fmt.Sprintf("%s\r\n%s\r\n", data.Subject, "hello world")
	by := bytes.NewBufferString(str).Bytes()
	_, err = writer.Write(by)
	if err != nil {
		return err
	}
	return nil

	/*
		str := fmt.Sprintf("%s\r\n%s\r\n", data.Subject, "hello world")
		by := bytes.NewBufferString(str).Bytes()
		smtpAuth := smtp.PlainAuth("", user_name, password, host)
		err = smtp.SendMail(fmt.Sprintf("%s:%d", host, port), smtpAuth, user_name, []string{data.Email}, by)
		return err
	*/
	/*
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
	*/
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
