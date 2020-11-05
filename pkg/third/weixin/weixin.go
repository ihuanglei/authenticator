package weixin

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ihuanglei/authenticator/pkg/third"
	"github.com/simplexwork/common"
)

var _ third.Third = WeiXin{}

const (
	// AuthorizeURL 授权
	_AuthorizeURL = "https://open.weixin.qq.com/connect/qrconnect?appid=%v&redirect_uri=%v&response_type=code&scope=snsapi_login&state=%v#wechat_redirect"

	// _AuthorizeOnInternalURL 微信内授权
	// _AuthorizeOnInternalURL = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%v&redirect_uri=%v&response_type=code&scope=%v&state=%v#wechat_redirect"

	// _AccessTokenURL 获取token
	_AccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%v&secret=%v&code=%v&grant_type=authorization_code"

	// _UserURL 用户信息
	_UserURL = "https://api.weixin.qq.com/sns/userinfo?access_token=%v&openid=%v"
)

// WeiXin .
type WeiXin struct {
	AppID       string `json:"client_id"`
	AppSecret   string `json:"client_secret"`
	RedirectURL string `json:"redirect_url"`
	Display     string `json:"display"`
}

type accessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	Error        string `json:"errmsg"`
	Code         int32  `json:"errcode"`
}

type user struct {
	ID       string `json:"openid"`
	UnionID  string `json:"unionid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"headimgurl"`
	Gender   string `json:"sex"`
	Province string `json:"province"`
	City     string `json:"city"`
	Error    string `json:"errmsg"`
	Code     int32  `json:"errcode"`
}

// AccessToken 授权
func (wx WeiXin) getAccessToken(code string) (*accessToken, error) {
	resp, err := http.Get(fmt.Sprintf(_AccessTokenURL, wx.AppID, wx.AppSecret, code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var accessToken accessToken
	err = common.FromJSON(body, &accessToken)
	if err != nil {
		return nil, err
	}
	if accessToken.Code != 0 {
		return nil, fmt.Errorf("error: %v", accessToken.Error)
	}
	return &accessToken, nil
}

// GetAuthorizeURL .
func (wx WeiXin) GetAuthorizeURL(state string) string {
	return fmt.Sprintf(_AuthorizeURL, wx.AppID, wx.RedirectURL, state)
}

// UserInfo 用户信息
func (wx WeiXin) getUserInfo(accessToken *accessToken) (*user, error) {
	url := fmt.Sprintf(_UserURL, accessToken.AccessToken, accessToken.OpenID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var user user
	err = common.FromJSON(body, &user)
	if err != nil {
		return nil, err
	}
	if user.Code != 0 {
		return nil, fmt.Errorf("error: %v", user.Error)
	}
	return &user, nil
}

// GetUser .
func (wx WeiXin) GetUser(code string) (*third.User, error) {
	accessToken, err := wx.getAccessToken(code)
	if err != nil {
		return nil, err
	}
	user, err := wx.getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}
	u := third.User{
		OpenID:   accessToken.OpenID,
		TP:       wx.GetType(),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Province: user.Province,
		City:     user.City,
	}
	return &u, nil
}

// GetType .
func (WeiXin) GetType() string {
	return "weixin"
}
