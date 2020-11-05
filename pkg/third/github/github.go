package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ihuanglei/authenticator/pkg/third"
	"github.com/simplexwork/common"
)

var _ third.Third = Github{}

const (
	// 授权
	_AuthorizeURL = "https://github.com/login/oauth/authorize?client_id=%v&state=%v&redirect_uri=%v"

	// 获取token
	_AccessTokenURL = "https://github.com/login/oauth/access_token"

	// 用户信息
	_UserURL = "https://api.github.com/user"
)

// Github .
type Github struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

// 接口返回的用户字段
type user struct {
	ID       int    `json:"id"`
	Nickname string `json:"name"`
	Avatar   string `json:"avatar_url"`
	Province string `json:"location"`
}

type accessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (github Github) getAccessToken(code string) (*accessToken, error) {
	data := url.Values{}
	data.Add("client_id", github.ClientID)
	data.Add("client_secret", github.ClientSecret)
	data.Add("redirect_uri", github.RedirectURL)
	data.Add("code", code)

	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest("POST", _AccessTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var accessToken accessToken
	err = common.FromJSON(body, &accessToken)
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}

func (github Github) getUserInfo(accessToken *accessToken) (*user, error) {
	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest("GET", _UserURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", accessToken.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var user user
	err = common.FromJSON(body, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUser .
func (github Github) GetUser(code string) (*third.User, error) {
	accessToken, err := github.getAccessToken(code)
	if err != nil {
		return nil, err
	}
	user, err := github.getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}
	u := third.User{
		OpenID:   common.IntToStr(user.ID),
		TP:       github.GetType(),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Province: user.Province,
	}
	return &u, nil
}

// GetAuthorizeURL .
func (github Github) GetAuthorizeURL(state string) string {
	return fmt.Sprintf(_AuthorizeURL, github.ClientID, state, url.QueryEscape(github.RedirectURL))
}

// GetType .
func (Github) GetType() string {
	return "github"
}
