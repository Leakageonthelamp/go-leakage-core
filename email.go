package core

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"io"
	"net/http"

	"gopkg.in/gomail.v2"
)

var EmailServiceSendFailed = Error{
	Status:  http.StatusBadGateway,
	Code:    "SEND_EMAIL_ERROR",
	Message: "can not send email",
}

var EmailServiceParserError = Error{
	Status:  http.StatusInternalServerError,
	Code:    "EMAIL_PARSING_ERROR",
	Message: "can not parse email",
}

type IEmail interface {
	SendHTML(from string, to []string, subject string, body string) IError
	SendHTMLWithAttach(from string, to []string, subject string, body string, file []byte, fileName string) IError
	SendText(from string, to []string, subject string, body string) IError
	ParseHTMLToString(path string, data interface{}) (string, IError)
}

type email struct {
	ctx IContext
}

func (s email) ParseHTMLToString(path string, data interface{}) (string, IError) {
	if data == nil {
		data = make(map[string]interface{})
	}

	t := template.New("")
	t, err := t.ParseFiles(path)
	if err != nil {
		return "", s.ctx.NewError(err, EmailServiceParserError)
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "template", data); err != nil {
		return "", s.ctx.NewError(err, EmailServiceParserError)
	}

	return tpl.String(), nil
}

func (s email) SendHTML(from string, to []string, subject string, body string) IError {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.send(m)
}

func (s email) SendHTMLWithAttach(from string, to []string, subject string, body string, fileByte []byte, fileName string) IError {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if len(fileByte) > 0 {
		m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(fileByte)
			return err
		}),
			gomail.SetHeader(map[string][]string{"Content-Type": {"application/pdf"}}),
		)
	}

	return s.send(m)
}

func (s email) SendText(from string, to []string, subject string, body string) IError {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return s.send(m)
}

func (s email) send(msg *gomail.Message) IError {
	dialer := gomail.NewDialer(
		s.ctx.ENV().Config().EmailServer,
		s.ctx.ENV().Config().EmailPort,
		s.ctx.ENV().Config().EmailUsername,
		s.ctx.ENV().Config().EmailPassword)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := dialer.DialAndSend(msg)
	if err != nil {
		return s.ctx.NewError(err, EmailServiceSendFailed, msg)
	}

	return nil
}

func NewEmail(ctx IContext) IEmail {
	return &email{ctx: ctx}
}
