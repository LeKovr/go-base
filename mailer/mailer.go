package mailer

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"log"
)

// -----------------------------------------------------------------------------

// Flags is a package flags sample
// in form ready for use with github.com/jessevdk/go-flags
type Flags struct {
	Port     int    `long:"smtp_port"  default:"25"        description:"SMTP server port"`
	Host     string `long:"smtp_host"  default:"localhost" description:"SMTP server host"`
	NoTLS    bool   `long:"smtp_notls" description:"Disable TLS cert checking"`
	Login    string `long:"smtp_login" description:"SMTP sender login"`
	Pass     string `long:"smtp_pass"  description:"SMTP sender password"`
	Copy     string `long:"smtp_copy"  description:"SMTP bcopy address"`
	From     string `long:"smtp_from"  description:"SMTP sender email"`
	FromName string `long:"smtp_fromname"  description:"SMTP sender name"`
}

// App is a package general type
type App struct {
	Log    *log.Logger
	Config *Flags
}

// New creates mailer object
// Configuration should be set via functional options
func New(logger *log.Logger, cfg *Flags, options ...func(a *App) error) (*App, error) {
	a := App{Config: cfg, Log: logger}
	for _, option := range options {
		err := option(&a)
		if err != nil {
			return nil, err
		}
	}
	a.Log.Printf("debug: Mail config: %+v", cfg)
	return &a, nil
}

// Send used to send emails
func (a App) Send(email, name, subject, buf string, files []string) (err error) {

	if a.Config.From == "" {
		a.Log.Printf("warn: Mail config field FROM does not set. Skip mailing")
		return
	}
	msg := gomail.NewMessage()
	if a.Config.FromName == "" {
		msg.SetHeader("From", a.Config.From)
	} else {
		msg.SetAddressHeader("From", a.Config.From, a.Config.FromName)
	}
	msg.SetHeader("To", email)
	if a.Config.Copy != "" {
		msg.SetHeader("Bcc", a.Config.Copy)
	}
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/html", buf)

	for _, name := range files {
		msg.Attach(name)
	}
	mailer := gomail.NewDialer(a.Config.Host, a.Config.Port, a.Config.Login, a.Config.Pass)

	if a.Config.NoTLS {
		mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	err = mailer.DialAndSend(msg)
	return
}
