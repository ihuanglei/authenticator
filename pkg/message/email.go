package message

import (
	"fmt"

	"github.com/ihuanglei/authenticator/pkg/logger"
	"gopkg.in/gomail.v2"
)

var _ Message = MailMessage{}

// MailMessage 邮件消息
type MailMessage struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	To       []string
}

// Dump .
func (m MailMessage) Dump() {
	logger.Debugf(fmt.Sprintf("subject:%s to:%v body: %s", m.Subject, m.To, m.Body))
}

// Send 发送
func sendMail(msg *MailMessage) error {
	defer msg.Dump()
	m := gomail.NewMessage()
	m.SetHeader("From", msg.Username)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)
	d := gomail.NewDialer(msg.Host, msg.Port, msg.Username, msg.Password)
	return d.DialAndSend(m)
}
