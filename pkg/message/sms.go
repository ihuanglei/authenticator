package message

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ihuanglei/authenticator/pkg/logger"
	"github.com/simplexwork/common"
)

var _ Message = SMSMessage{}

var _httpClient = &http.Client{Timeout: time.Second * 5}

// SMSMessage 短信消息
type SMSMessage struct {
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Querys map[string]string `json:"query"`
	Params map[string]string `json:"body"`
	To     string
	Body   string
}

const _logPrefix = "[SMS]"

// Dump .
func (m SMSMessage) Dump() {
	logger.Debugf(fmt.Sprintf("%s to:%s body: %s", _logPrefix, m.To, m.Body))
}

// Send 发送
func sendSMS(msg *SMSMessage) (err error) {
	defer msg.Dump()
	var req *http.Request
	method := common.ToUpper(msg.Method)
	body := url.Values{}
	for k, v := range msg.Params {
		body.Add(k, v)
	}
	req, err = http.NewRequest(method, msg.URL, strings.NewReader(body.Encode()))
	for k, v := range msg.Header {
		req.Header.Add(k, v)
	}
	query := url.Values{}
	for k, v := range msg.Querys {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	if err != nil {
		err = fmt.Errorf("%s send sms message error:%v", _logPrefix, err)
		return
	}
	logger.Debugf("%s send sms message request: %v", _logPrefix, req)
	var resp *http.Response
	resp, err = _httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("%s send sms message error:%v", _logPrefix, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s send sms message error:%d %s", _logPrefix, resp.StatusCode, resp.Status)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("%s send sms message error:%v", _logPrefix, err)
		return
	}
	logger.Debugf("%s send sms message response: %v", _logPrefix, string(respBody))
	return
}
