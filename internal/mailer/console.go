package mailer

import "log"

type Console struct {
	Logger *log.Logger
}

func (c *Console) SendVerificationCode(to, code string) {
	l := c.Logger
	if l == nil {
		l = log.Default()
	}
	l.Printf("[mail] verification code for %s: %s\n", to, code)
}
