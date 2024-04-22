package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/ziadrahmatullah/ordent-test/config"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
	smtpTestServer    = "localhost:1025"
)

type SmtpGmail interface {
	SendEmail(string, string, bool) error
	SendEmailTest(string, string, bool) error
}

type smtpGmail struct {
	name       string
	address    string
	password   string
	prefixLink string
}

func NewSmtpGmail() SmtpGmail {
	emailConfig := config.NewEmailConfig()
	return &smtpGmail{
		name:       emailConfig.Name,
		address:    emailConfig.Address,
		password:   emailConfig.Password,
		prefixLink: emailConfig.PrefixLink,
	}
}

func emailVerifyContent(otp string) (subject, content string) {
	subject = "Account Verification"
	content = fmt.Sprintf(`
		<p>Please click the link to start your verification process!</p>
		<br />
		<a href="%s" style="text-decoration: none; color: white; background-color: #36A5B2; padding: 10px 20px; border-radius: 0.3rem; font-weight: bold; display: inline-block;">Verification Process</a>
		<br />
	`, otp)
	return subject, content
}

func emailForgotPasswordContent(otp string) (subject, content string) {
	subject = "Change Password Verification"
	content = fmt.Sprintf(`
		<p>We received a request to reset your password. Click the button below to reset it:</p>
		<br />
		<a href="%s" style="text-decoration: none; color: white; background-color: #36A5B2; padding: 10px 20px; border-radius: 0.3rem; font-weight: bold; display: inline-block;">Reset Password</a>
		<br />
	`, otp)
	return subject, content
}

func (r *smtpGmail) SendEmail(token, to string, isVerify bool) error {
	receiver := []string{to}
	link := r.prefixLink + token
	var subject, content string
	if isVerify {
		subject, content = emailVerifyContent(link)
	} else {
		subject, content = emailForgotPasswordContent(link)
	}
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", r.name, r.address)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = receiver
	smtpAuth := smtp.PlainAuth("", r.address, r.password, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func (r *smtpGmail) SendEmailTest(token, to string, isVerify bool) error {
	receiver := []string{to}
	link := r.prefixLink + token
	var subject, content string
	if isVerify {
		subject, content = emailVerifyContent(link)
	} else {
		subject, content = emailForgotPasswordContent(link)
	}
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, content)
	smtpAuth := smtp.PlainAuth("", "example@gmail.com", "1234", "localhost")
	return smtp.SendMail(smtpTestServer, smtpAuth, r.address, receiver, []byte(message))
}
