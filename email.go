package xmusemail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

type (
	//SMTP is the basic SMTP client
	SMTP interface {
		Send(message Message, subject, template string, to ...string) error
	}
	smtpClient struct {
		auth      Auth
		tlsCongif *tls.Config
		logger    LeveledLoggerInterface
		log       bool
	}
	//Auth is the authentication information for the SMTP server
	Auth struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		From     string `json:"from"`
	}
	//Message is the basic message
	Message map[EngineKey]EngineValue
	//EngineValue is the value of the message
	EngineValue interface{}
	//EngineKey is the key of the message
	EngineKey string
	// LeveledLoggerInterface is the interface for leveled logger
	LeveledLoggerInterface interface {
		// Debugf logs a debug message using Printf conventions.
		Debugf(format string, v ...interface{})

		// Errorf logs a warning message using Printf conventions.
		Errorf(format string, v ...interface{})

		// Infof logs an informational message using Printf conventions.
		Infof(format string, v ...interface{})

		// Warnf logs a warning message using Printf conventions.
		Warnf(format string, v ...interface{})
	}
)

// NewSMTP creates a new SMTP client
func NewSMTP(auth Auth, tlsConfig *tls.Config, logger LeveledLoggerInterface, log bool) SMTP {
	core := &smtpClient{
		auth: Auth{
			Host:     auth.Host,
			Port:     auth.Port,
			Password: auth.Password,
			From:     auth.From,
		},
	}
	core.logger = logger
	core.tlsCongif = tlsConfig
	core.log = log
	return core
}

func (sc smtpClient) CreateMessageFromTemplate(parts ...string) string {
	return ""
}
func (sc smtpClient) createHeaders(from mail.Address, to, subject string) string {
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""
	header := ""
	for k, v := range headers {
		header += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return header
}
func (sc smtpClient) createMessageFromTemplate(template string, message Message) string {
	for k, v := range message {
		template = strings.Replace(template, fmt.Sprint(k), fmt.Sprintf("%v", v), -1)
	}
	return template
}
func (sc smtpClient) createEmail(from mail.Address, subject, to, template string, message Message) string {
	return fmt.Sprintf("%s%s", sc.createHeaders(from, to, subject), sc.createMessageFromTemplate(template, message))
}
func (sc smtpClient) Send(message Message, subject, template string, to ...string) error {
	serverName := fmt.Sprintf("%s:%d", sc.auth.Host, sc.auth.Port)
	host, _, _ := net.SplitHostPort(serverName)
	auth := smtp.PlainAuth("", sc.auth.From, sc.auth.Password, host)
	for _, v := range to {
		from := mail.Address{"", sc.auth.From}
		to := mail.Address{"", v}
		emailMessage := sc.createEmail(from, subject, v, template, message)
		err := sc.send(auth, from, to, serverName, host, emailMessage)
		if err != nil {
			if sc.log {
				sc.logger.Errorf("error on sending to %s with message %v:", to.String, err)
			}
		}
		return err
	}
	return nil
}
func (sc smtpClient) send(auth smtp.Auth, from, to mail.Address, serverName, host, message string) error {
	conn, err := tls.Dial("tcp", serverName, sc.tlsCongif)
	if err != nil {
		return err
	}
	smtpclient, err := smtp.NewClient(conn, host)
	if err != nil {
		return ErrClientIsNotValid
	}
	defer smtpclient.Quit()
	if err = smtpclient.Auth(auth); err != nil {
		return ErrUnauthorized
	}
	if err = smtpclient.Mail(from.Address); err != nil {
		return ErrYourEmailIsInvalid
	}
	if err = smtpclient.Rcpt(to.Address); err != nil {
		return ErrRCPTConnection
	}
	w, err := smtpclient.Data()
	if err != nil {
		return ErrOnDataCommand
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		return ErrOnWriteData
	}
	err = w.Close()
	if err != nil {
		return ErrOnCloseDataPipe
	}
	return nil
}
