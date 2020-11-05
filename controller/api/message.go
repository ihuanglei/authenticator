package api

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/logger"
	"github.com/ihuanglei/authenticator/pkg/message"
	"github.com/simplexwork/common"
)

const (
	// 通知类邮箱配置
	cateEmail = "email"
	tpEmail   = "notice"

	// 验证码短信配置
	cateSMS = "sms"
	tpSMS   = "identify"

	// 模板类型
	cateTmpl = "tmpl"
	// 邮件注册激活码模板
	tmplEmailReg = "email_activate"
	// 邮件绑定验证码模板
	tmplEmailBind = "email_bind"
	// 邮件忘记密码
	tmplEmailForgot = "email_forgot"
	// 短信注册验证码模板
	tmplMobileReg = "mobile_reg"
)

// 邮箱注册激活
func sendActivateMessageWithEmail(userID common.ID, activateCode, email string) {
	defer messageRecover()
	ev, err := buildMailMessage(tmplEmailReg)
	if err != nil {
		logger.Error(err)
		return
	}
	activateCode1 := fmt.Sprintf("%v@%s", userID, activateCode)
	ev.Body = strings.ReplaceAll(ev.Body, "{code}", common.Base64Encode([]byte(activateCode1)))
	ev.To = []string{email}
	message.SendMessage(ev)
}

// 邮箱绑定\更新
func sendBindMessageWithEmail(email, code string) {
	defer messageRecover()
	ev, err := buildMailMessage(tmplEmailBind)
	if err != nil {
		logger.Error(err)
		return
	}
	ev.Body = strings.ReplaceAll(ev.Body, "{code}", code)
	ev.To = []string{email}
	message.SendMessage(ev)
}

// 邮箱忘记密码
func sendForgotMessageWithEmail(email, code string) {
	defer messageRecover()
	ev, err := buildMailMessage(tmplEmailForgot)
	if err != nil {
		logger.Error(err)
		return
	}
	ev.Body = strings.ReplaceAll(ev.Body, "{code}", code)
	ev.To = []string{email}
	message.SendMessage(ev)
}

func buildMailMessage(tmpl string) (*message.MailMessage, error) {
	// FIXME: 是否要缓存当前邮件数据?
	dictDto, err := models.GetOneDict(cateEmail, tpEmail)
	if err != nil {
		return nil, err
	}
	var ev *message.MailMessage
	if err := common.FromJSON([]byte(dictDto.Value), &ev); err != nil {
		return nil, err
	}
	dictDto1, err := models.GetOneDict(cateTmpl, tmpl)
	if err != nil {
		return nil, err
	}
	if err := common.FromJSON([]byte(dictDto1.Value), &ev); err != nil {
		return nil, err
	}
	return ev, nil
}

/***** ******/

// 短信验证码，统一格式
func sendMessageWithMobile(mobile, code string) {
	replaceVarFunc := func(s, mobile, code, body string) string {
		s = strings.ReplaceAll(s, "{to}", mobile)
		s = strings.ReplaceAll(s, "{code}", code)
		s = strings.ReplaceAll(s, "{body}", body)
		return s
	}
	defer messageRecover()
	ev, err := buildSMSMessage(tmplMobileReg)
	if err != nil {
		logger.Error(err)
		return
	}
	body := strings.ReplaceAll(ev.Body, "{code}", code)
	for k, v := range ev.Params {
		ev.Params[k] = replaceVarFunc(v, mobile, code, body)
	}
	for k, v := range ev.Querys {
		ev.Querys[k] = replaceVarFunc(v, mobile, code, body)
	}
	ev.To = mobile
	ev.Body = body
	message.SendMessage(ev)
}

func buildSMSMessage(tmpl string) (*message.SMSMessage, error) {
	// FIXME: 是否要缓存当前短信数据?
	dictDto, err := models.GetOneDict(cateSMS, tpSMS)
	if err != nil {
		return nil, err
	}
	var ev *message.SMSMessage
	if err := common.FromJSON([]byte(dictDto.Value), &ev); err != nil {
		return nil, err
	}
	dictDto1, err := models.GetOneDict(cateTmpl, tmpl)
	if err != nil {
		return nil, err
	}
	ev.Body = dictDto1.Value
	return ev, nil
}

/***** ******/

func messageRecover() {
	if err := recover(); err != nil {
		var buf [1024]byte
		n := runtime.Stack(buf[:], false)
		logger.Error(string(buf[:n]))
	}
}
