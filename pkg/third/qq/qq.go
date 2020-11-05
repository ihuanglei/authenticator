package qq

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/ihuanglei/authenticator/pkg/third"

	"github.com/simplexwork/common"
)

var _ third.Third = QQ{}

const (
	// 授权
	_AuthorizeURL = "https://graph.qq.com/oauth2.0/authorize?scope=all&client_id=%v&response_type=code&state=%v&display=%v&redirect_uri=%v"

	// 获取token
	_AccessTokenURL = "https://graph.qq.com/oauth2.0/token?client_id=%v&client_secret=%v&code=%v&grant_type=authorization_code&redirect_uri=%v"

	// 获取用户openid
	_OpenIDURL = "https://graph.qq.com/oauth2.0/me?access_token=%v"

	// 用户信息
	_UserURL = "https://graph.qq.com/user/get_user_info?oauth_consumer_key=%v&access_token=%v&openid=%v"
)

// QQ .
type QQ struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	Display      string `json:"display"`
}

// qq接口返回的用户字段
type user struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"figureurl_qq_2"`
	Gender   string `json:"gender"`
	Province string `json:"province"`
	City     string `json:"city"`
	Error    string `json:"msg"`
	Code     int32  `json:"ret"`
}

func (qq QQ) getAccessToken(code string) (string, error) {
	v, err := qq.httpGet(fmt.Sprintf(_AccessTokenURL, qq.ClientID, qq.ClientSecret, code, qq.RedirectURL))
	if err != nil {
		return "", err
	}
	return v.Get("access_token"), nil
}

func (qq QQ) getOpenID(accessToken string) (string, error) {
	v, err := qq.httpGet(fmt.Sprintf(_OpenIDURL, accessToken))
	if err != nil {
		return "", err
	}
	return v.Get("openid"), nil
}

func (qq QQ) getUserInfo(accessToken, openID string) (*user, error) {
	var user user
	resp, err := http.Get(fmt.Sprintf(_UserURL, qq.ClientID, accessToken, openID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = common.FromJSON(body, &user)
	if err != nil {
		return nil, err
	}
	if user.Code != 0 {
		return nil, fmt.Errorf("error: %v", user.Error)
	}
	return &user, nil
}

func (QQ) parse2Query(str string) (string, error) {
	if strings.Contains(str, "callback") {
		reg, err := regexp.Compile("callback\\(|\\)|\\{|\\}|;|\"")
		if err != nil {
			return "", err
		}
		return common.Trim(strings.Replace(strings.Replace(reg.ReplaceAllString(str, ""), ":", "=", -1), ",", "&", -1)), nil
	}
	return str, nil
}

func (qq QQ) httpGet(uri string) (url.Values, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result, err := qq.parse2Query(string(body))
	if err != nil {
		return nil, err
	}
	v, err := url.ParseQuery(result)
	if err != nil {
		return nil, err
	}
	if v.Get("error") != "" {
		return nil, fmt.Errorf("error: %v %v", v.Get("error"), v.Get("error_description"))
	}
	return v, nil
}

// 实现third接口

// GetUser 获取QQ用户信息,流程如下:
// 1 获取code
// 2 获取操作令牌
// 3 使用操作令牌获取用户编号
// 4 使用操作令牌和用户编号获取用户信息
func (qq QQ) GetUser(code string) (*third.User, error) {
	accessToken, err := qq.getAccessToken(code)
	if err != nil {
		return nil, err
	}
	openID, err := qq.getOpenID(accessToken)
	if err != nil {
		return nil, err
	}
	user, err := qq.getUserInfo(accessToken, openID)
	if err != nil {
		return nil, err
	}

	u := third.User{
		OpenID:   openID,
		TP:       qq.GetType(),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Province: user.Province,
		City:     user.City,
	}
	return &u, nil
}

// GetAuthorizeURL .
func (qq QQ) GetAuthorizeURL(state string) string {
	return fmt.Sprintf(_AuthorizeURL, qq.ClientID, state, qq.Display, url.QueryEscape(qq.RedirectURL))
}

// GetType .
func (QQ) GetType() string {
	return "qq"
}
