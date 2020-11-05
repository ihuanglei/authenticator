package weixin

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/third"
	"github.com/simplexwork/common"
)

var _ third.Third = MinPro{}

const (
	// 小程序appid和secret获取token
	_MPAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v"

	// 微信小程序获取openid code2Session
	_MPCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%v&secret=%v&js_code=%v&grant_type=authorization_code"
)

// MinPro .
type MinPro struct {
	AppID     string `json:"client_id"`
	AppSecret string `json:"client_secret"`
}

// 小程序用户
type mpUser struct {
	OpenID    string `json:"openId"`
	Nickname  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Avatar    string `json:"avatarUrl"`
	UnionID   string `json:"unionId"`
	Watermark struct {
		Appid     string `json:"appid"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
	Error string `json:"errmsg"`
	Code  int32  `json:"errcode"`
}

type mpUserPhone struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		Appid     string `json:"appid"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
}

// AccessToken 令牌
type mpAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
	Error       string `json:"errmsg"`
	Code        int32  `json:"errcode"`
}

// Session 小程序返回的session
type mpSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	Error      string `json:"errmsg"`
	Code       int32  `json:"errcode"`
}

// 小程序code换session
func (mp MinPro) code2Session(code string) (*mpSession, error) {
	resp, err := http.Get(fmt.Sprintf(_MPCode2SessionURL, mp.AppID, mp.AppSecret, code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var session mpSession
	err = common.FromJSON(body, &session)
	if err != nil {
		return nil, err
	}
	if session.Code != 0 {
		return nil, fmt.Errorf("error: %v", session.Error)
	}
	return &session, nil
}

// AuthorizeOnInternal 微信内授权页面
// func (wx MinPro) authorizeOnInternal(state string, isAuth bool) string {
// 	scope := "snsapi_base"
// 	if isAuth {
// 		scope = "snsapi_userinfo"
// 		// wx.isAuth = true
// 	}
// 	return fmt.Sprintf(_AuthorizeOnInternalURL, wx.AppID, wx.RedirectURL, scope, state)
// }

// AccessToken 授权
func (mp MinPro) getAccessToken() (*accessToken, error) {
	resp, err := http.Get(fmt.Sprintf(_MPAccessTokenURL, mp.AppID, mp.AppSecret))
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

// GetMPUser .
func (mp MinPro) GetMPUser(sessionKey, iv, encryptedData string) (*third.User, error) {
	bytes, err := mp.decrypt(sessionKey, iv, encryptedData)
	if err != nil {
		return nil, err
	}
	var user mpUser
	err = common.FromJSON(bytes, &user)
	if err != nil {
		return nil, err
	}
	var gender string
	switch user.Gender {
	case 0:
		gender = "unknown"
	case 1:
		gender = "male"
	case 2:
		gender = "female"
	}
	u := third.User{
		OpenID:   user.OpenID,
		TP:       mp.GetType(),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   gender,
		Province: user.Province,
		City:     user.City,
	}
	return &u, nil
}

// GetMPUserPhone .
func (mp MinPro) GetMPUserPhone(sessionKey, iv, encryptedData string) (string, error) {
	bytes, err := mp.decrypt(sessionKey, iv, encryptedData)
	if err != nil {
		return "", err
	}
	var userPhone mpUserPhone
	err = common.FromJSON(bytes, &userPhone)
	if err != nil {
		return "", err
	}
	return userPhone.PurePhoneNumber, nil
}

// GetSession 获取session
func (mp MinPro) GetSession(code string) (string, string, error) {
	session, err := mp.code2Session(code)
	if err != nil {
		return "", "", err
	}
	return session.OpenID, session.SessionKey, nil
}

// GetUser .
func (MinPro) GetUser(code string) (*third.User, error) {
	return nil, errors.ErrUnknown
}

// GetAuthorizeURL .
func (MinPro) GetAuthorizeURL(state string) string {
	return ""
}

// GetType .
func (MinPro) GetType() string {
	return "weixinmp"
}

func (MinPro) decrypt(sessionKey, iv, encryptedData string) ([]byte, error) {
	sessionKeyBytes, err := common.Base64Decode(sessionKey)
	if err != nil {
		return nil, err
	}
	ivBytes, err := common.Base64Decode(iv)
	if err != nil {
		return nil, err
	}
	encryptedBytes, err := common.Base64Decode(encryptedData)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(sessionKeyBytes)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, ivBytes)
	origData := make([]byte, len(encryptedBytes))
	blockMode.CryptBlocks(origData, encryptedBytes)
	for i, ch := range origData {
		if ch == '\x0e' {
			origData[i] = ' '
		}
	}
	return origData, nil
}
