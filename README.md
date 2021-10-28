# xmus-email
simple email client using net/SMTP


sample usage :
```GOLANG
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
	"<%MESSAGE%>": "Lorem ipsu dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia mollit anim id est laborum.",
	"<%FOOTER%>":  "copyright my company",
}
err := e.Send(msessage, "subject", xmusemail.SampleHolloweenTemplate, "amupxm@gmail.com")
if err != nil {
	log.Error(err)
}
```