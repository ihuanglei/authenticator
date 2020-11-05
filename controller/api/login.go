package api

import (
	"fmt"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/simplexwork/cache"
	"github.com/simplexwork/common"

	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/config"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/convert"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/ihuanglei/authenticator/pkg/third"
	"github.com/ihuanglei/authenticator/pkg/third/github"
	"github.com/ihuanglei/authenticator/pkg/third/qq"
	"github.com/ihuanglei/authenticator/pkg/third/weibo"
	"github.com/ihuanglei/authenticator/pkg/third/weixin"
)

const (
	tokenKey      string = "__token_%v"
	tokenThirdKey string = "__token_thrid_%v"

	sessionKeyWithWeiXinMP = "__session_key_weixinmp_%v"
)

// Login 手机号\邮箱\用户名和密码登录
// @tags 前端 - 用户登录
// @Summary 手机号\邮箱\用户名和密码登录
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param login_name formData string false "[手机号|邮箱|用户名]"
// @Param password formData string false "密码"
// @Router /api/login [post]
func Login(form st.LoginForm, ctx *context.Context) {
	token, err := login(&form, ctx, func(loginDto *st.LoginDto) (*st.UserDto, error) {
		return models.Login(loginDto)
	})
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSON(string(token))
}

// LoginByMobile 手机号和验证码登录
// @tags 前端 - 用户登录
// @Summary 手机号和验证码登录
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param mobile formData string false "手机号"
// @Param code formData string false "验证码"
// @Router /api/login/mobile [post]
func LoginByMobile(form st.LoginWithMobileAndCodeForm, config *config.Config, cache cache.Cache, ctx *context.Context) {
	key := fmt.Sprintf(codeKeyWithLogin, form.Mobile)
	tmpCode, err := cache.GetString(key)
	if err != nil || tmpCode != form.Code {
		ctx.BadRequestByError(errors.ErrCode)
		return
	}
	token, err := login(&form, ctx, func(loginDto *st.LoginDto) (*st.UserDto, error) {
		return models.LoginByMobile(loginDto)
	})
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	cache.Del(key)
	ctx.JSON(string(token))
}

// LoginByThirdCode 第三方使用code登录
// @tags 前端 - 用户登录
// @Summary 第三方QQ，微信，微博使用code登录
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param id path string true "第三方编号"
// @Param code formData string false "第三方授权后返回的code"
// @Param state formData string false "第三方授权后返回的state"
// @Router /api/login/th/{id} [post]
func LoginByThirdCode(form st.LoginWithThirdCodeForm, cache cache.Cache, ctx *context.Context) {
	id := ctx.Params("id")
	thirdUser, err := getThirdUser(id, form.Code)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	form.OpenID = thirdUser.OpenID
	form.Type = thirdUser.TP
	token, err := login(&form, ctx, func(loginDto *st.LoginDto) (*st.UserDto, error) {
		return models.LoginByOpenID(loginDto)
	})
	// 三方没有注册过返回三方用户信息
	if err != nil {
		if e, ok := err.(errors.Error); ok && e.Code() == errors.ErrUserNotExist.Code() {
			if err := cache.Set(fmt.Sprintf(tokenThirdKey, form.Code), thirdUser, time.Minute*10); err != nil {
				ctx.Error(err)
				return
			}
			ctx.JSONByCode(e.Code(), map[string]string{"nickname": thirdUser.Nickname, "avatar": thirdUser.Avatar})
		} else if err != nil {
			ctx.Error(err)
		}
		return
	}
	ctx.JSON(string(token))
}

// LoginByWeiXinMPCode 微信小程序code获取session
// @tags 前端 - 用户登录
// @Summary 微信小程序登录
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param id path string true "第三方编号"
// @Param code formData string false "微信小程序通过wx.login获取的临时登录凭证"
// @Param only_session formData bool false "是否只登录微信小程序"
// @Router /api/login/th/weixinmp/{id} [post]
func LoginByWeiXinMPCode(form st.LoginWithWeiXinMPCodeForm, cache cache.Cache, ctx *context.Context) {

	buildSession := func(openID, tp, sessionKey string) (string, error) {
		s := common.RandomString(32)
		retry := 1
		for _, err := cache.Get(s); err == nil; {
			if retry > 10 {
				return "", errors.ErrUnknown
			}
			s = common.RandomString(32)
			retry++
		}
		var thirdUser third.User
		thirdUser.OpenID = openID
		thirdUser.TP = tp
		thirdUser.Ext = sessionKey
		if err := cache.Set(fmt.Sprintf(sessionKeyWithWeiXinMP, s), thirdUser, time.Hour*24); err != nil {
			return "", err
		}
		return s, nil
	}

	id := ctx.Params("id")
	mp, err := getThird(id)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	v, ok := mp.(*weixin.MinPro)
	if !ok {
		ctx.BadRequestByError(errors.ErrDictNotFound)
		return
	}
	openID, sessionKey, err := v.GetSession(form.WeiXinMPCode)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}

	if form.OnlySession {
		s, err := buildSession(openID, mp.GetType(), sessionKey)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.JSON(s)
		return
	}

	form.OpenID = openID
	form.Type = mp.GetType()
	token, err := login(&form, ctx, func(loginDto *st.LoginDto) (*st.UserDto, error) {
		return models.LoginByOpenID(loginDto)
	})
	if err != nil {
		if e, ok := err.(errors.Error); ok && e.Code() == errors.ErrUserNotExist.Code() {
			s, err := buildSession(openID, mp.GetType(), sessionKey)
			if err != nil {
				ctx.Error(err)
				return
			}
			var thirdUser third.User
			thirdUser.OpenID = openID
			thirdUser.TP = mp.GetType()
			thirdUser.Ext = sessionKey
			if err := cache.Set(fmt.Sprintf(sessionKeyWithWeiXinMP, s), thirdUser, time.Hour*3); err != nil {
				ctx.Error(err)
				return
			}
			ctx.JSONByCode(e.Code(), s)
		} else if err != nil {
			ctx.Error(err)
		}
		return
	}
	ctx.JSON(string(token))
}

// RedirectURLForThird 第三方登录
// @tags 前端 - 用户登录
// @Summary 第三方QQ，微信，微博登录地址
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} 第三方请求地址
// @Param id path string true "第三方编号"
// @Param state query string false "第三方授权后返回的state"
// @Router /api/login/th/{id} [get]
func RedirectURLForThird(form st.LoginWithThirdForm, ctx *context.Context) {
	id := ctx.Params("id")
	url, err := getAuthorizeURL(id, form.State)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSON(url)
}

func getAuthorizeURL(id, state string) (string, error) {
	third, err := getThird(id)
	if err != nil {
		return "", err
	}
	return third.GetAuthorizeURL(state), nil
}

func getThirdUser(id, code string) (*third.User, error) {
	third, err := getThird(id)
	if err != nil {
		return nil, err
	}
	thirdUser, err := third.GetUser(code)
	if err != nil {
		return nil, err
	}
	return thirdUser, nil
}

func getThird(id string) (third.Third, error) {
	dictDto, err := models.GetDictByID(common.StrToID(id))
	if err != nil {
		return nil, err
	}
	tp := dictDto.TP
	val := dictDto.Value
	if dictDto.TP == "" {
		return nil, errors.ErrDictNotFound
	}
	switch tp {
	case "qq":
		var qq qq.QQ
		err := common.FromJSON([]byte(val), &qq)
		if err != nil {
			return nil, err
		}
		return &qq, nil
	case "weibo":
		var weibo weibo.WeiBo
		err := common.FromJSON([]byte(val), &weibo)
		if err != nil {
			return nil, err
		}
		return &weibo, nil
	case "weixin":
		var weixin weixin.WeiXin
		err := common.FromJSON([]byte(val), &weixin)
		if err != nil {
			return nil, err
		}
		return &weixin, nil
	case "weixinmp":
		var weixinMP weixin.MinPro
		err := common.FromJSON([]byte(val), &weixinMP)
		if err != nil {
			return nil, err
		}
		return &weixinMP, nil
	case "github":
		var github github.Github
		err := common.FromJSON([]byte(val), &github)
		if err != nil {
			return nil, err
		}
		return &github, nil
	}
	return nil, errors.ErrDictNotFound
}

func login(form interface{}, ctx *context.Context, handle func(loginDto *st.LoginDto) (*st.UserDto, error)) (string, error) {
	loginDto := new(st.LoginDto)
	if err := convert.Map(form, loginDto); err != nil {
		return "", err
	}
	loginDto.IP = ctx.IP
	userDto, err := handle(loginDto)
	if err != nil {
		return "", err
	}
	token, err := getUserAndCreateJWTToken(userDto.UserID, ctx.Secret, ctx.Expire)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Authenticator %v", string(token)), nil
}

func getUserAndCreateJWTToken(userID common.ID, secret string, expire int64) ([]byte, error) {
	userInfoDto, err := models.GetUserInfoByID(userID)
	if err != nil {
		return nil, err
	}
	exp := time.Hour * 24 * time.Duration(expire)
	if err != nil {
		return nil, err
	}
	subjectMap := map[string]interface{}{}
	subjectMap["user_id"] = userInfoDto.UserID
	subjectMap["avatar"] = userInfoDto.Avatar
	subjectMap["nickname"] = userInfoDto.Nickname
	b, err := common.ToJSON(subjectMap)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	hs256 := jwt.NewHS256([]byte(secret))
	p := jwt.Payload{
		Issuer:         "authenticator",
		Audience:       []string{"authenticator"},
		Subject:        b,
		ExpirationTime: &jwt.Time{Time: now.Add(exp)},
		IssuedAt:       &jwt.Time{Time: now},
	}
	bs, err := jwt.Sign(p, hs256)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
