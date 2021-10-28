package main

import (
	"crypto/tls"

	xmusemail "github.com/amupxm/xmus-email"
	xmuslogger "github.com/amupxm/xmus-logger"
)

// SSL/TLS Email Example

func main() {
	auth := xmusemail.Auth{
		Host:     "mail.yourdomain.com",
		Port:     465,
		From:     "target@anotherdomain.com",
		Password: "passWordHere",
	}
	logger := xmuslogger.CreateLogger(&xmuslogger.LoggerOptions{LogLevel: xmuslogger.Error})
	log := logger.BeginWithPrefix("[xmus-email]")
	e := xmusemail.NewSMTP(auth,
		&tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "mail.yourdomain.com",
		},
		log,
		true,
	)
	msessage := xmusemail.Message{
		"<%PLACE%>":   "in google meet =)",
		"<%HEADER%>":  "some fancy header here",
		"<%MESSAGE%>": "Lorem ipsu dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		"<%FOOTER%>":  "copyright my company",
	}
	err := e.Send(msessage, "subject", xmusemail.SampleHolloweenTemplate, "amupxm@gmail.com")
	if err != nil {
		log.Error(err)
	}
}
