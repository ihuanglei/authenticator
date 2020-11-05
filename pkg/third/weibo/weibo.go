package weibo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ihuanglei/authenticator/pkg/third"
	"github.com/simplexwork/common"
)

var _ third.Third = WeiBo{}

const (
	// 授权
	_AuthorizeURL = "https://api.weibo.com/oauth2/authorize?scope=all&client_id=%v&response_type=code&state=%v&display=%v&redirect_uri=%v"

	// 获取token
	_AccessTokenURL = "https://api.weibo.com/oauth2/access_token"

	// 用户信息
	_UserURL = "https://api.weibo.com/2/users/show.json?access_token=%v&uid=%v"

	// 省接口
	_ProvinceURL = "https://api.weibo.com/2/common/get_province.json?access_token=%v&country=%v"

	// 市接口
	_CityURL = "https://api.weibo.com/2/common/get_city.json?access_token=%v&province=%v"

	// 中国编号
	_CN = "001"
)

// WeiBo 微博
type WeiBo struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	Display      string `json:"display"`
}

type accessToken struct {
	AccessToken string `json:"access_token"`
	UID         string `json:"uid"`
	Error       string `json:"error"`
	Code        int32  `json:"error_code"`
}

// 微博接口返回的用户字段
type user struct {
	ID       string `json:"idstr"`
	Nickname string `json:"screen_name"`
	Avatar   string `json:"avatar_large"`
	Gender   string `json:"gender"`
	Province string `json:"province"`
	City     string `json:"city"`
	Location string `json:"location"`
	Error    string `json:"error"`
	Code     int32  `json:"error_code"`
}

func (wb WeiBo) getAccessToken(code string) (*accessToken, error) {
	data := url.Values{}
	data.Add("client_id", wb.ClientID)
	data.Add("client_secret", wb.ClientSecret)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", wb.RedirectURL)
	data.Add("code", code)
	resp, err := http.PostForm(_AccessTokenURL, data)
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

func (WeiBo) getUserInfo(accessToken *accessToken) (*user, error) {
	url := fmt.Sprintf(_UserURL, accessToken.AccessToken, accessToken.UID)
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

func (WeiBo) getProvince(accessToken *accessToken, provinceID string) (string, error) {
	url := fmt.Sprintf(_ProvinceURL, accessToken.AccessToken, _CN)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data []map[string]string
	err = common.FromJSON(body, &data)
	if err != nil {
		return "", fmt.Errorf("get province error: %v", string(body))
	}
	for _, m := range data {
		if v, ok := m[provinceID]; ok {
			return v, nil
		}
	}
	return "", nil
}

func (WeiBo) getCity(accessToken *accessToken, provinceID, cityID string) (string, error) {
	url := fmt.Sprintf(_CityURL, accessToken.AccessToken, provinceID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data []map[string]string
	err = common.FromJSON(body, &data)
	if err != nil {

		return "", fmt.Errorf("get city error: %v", string(body))
	}

	for _, m := range data {
		if v, ok := m[cityID]; ok {
			return v, nil
		}
	}
	return "", nil
}

// 实现third接口

// GetUser 获取weibo用户信息,流程如下:
// 1 获取code
// 2 获取操作令牌和用户编号
// 3 使用操作令牌获取用户信息
// 4 使用操作令牌获取省市数据
// 5 合并数据
func (wb WeiBo) GetUser(code string) (*third.User, error) {
	accessToken, err := wb.getAccessToken(code)
	if err != nil {
		return nil, err
	}
	user, err := wb.getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	provinceID := fmt.Sprintf("%s%03s", _CN, user.Province)
	province, err := wb.getProvince(accessToken, provinceID)
	if err != nil {
		return nil, err
	}

	cityID := fmt.Sprintf("%s%03s", provinceID, user.City)
	city, err := wb.getCity(accessToken, provinceID, cityID)
	if err != nil {
		return nil, err
	}
	u := third.User{
		OpenID:   accessToken.UID,
		TP:       wb.GetType(),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Province: province,
		City:     city,
	}
	return &u, nil
}

// GetAuthorizeURL .
func (wb WeiBo) GetAuthorizeURL(state string) string {
	return fmt.Sprintf(_AuthorizeURL, wb.ClientID, state, wb.Display, url.QueryEscape(wb.RedirectURL))
}

// GetType .
func (WeiBo) GetType() string {
	return "weibo"
}
