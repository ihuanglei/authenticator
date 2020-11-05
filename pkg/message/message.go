package message

import "github.com/ihuanglei/authenticator/pkg/logger"

// Message 消息主体
type Message interface {
	Dump()
}

// SendMessage 发送消息
func SendMessage(message Message) {
	var err error
	switch msg := message.(type) {
	case *MailMessage:
		err = sendMail(msg)
	case *SMSMessage:
		err = sendSMS(msg)
	default:
		logger.Warn("message type not found")
	}
	if err != nil {
		logger.Panic(err)
	}
}
