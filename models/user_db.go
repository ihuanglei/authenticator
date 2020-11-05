package models

import (
	"time"

	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/logger"
	"github.com/ihuanglei/authenticator/pkg/region"
	"github.com/simplexwork/common"
	"xorm.io/builder"
)

// 用户登录历史
func getUserLogins(cond builder.Cond, page, limit int) (int64, []*userLogin, error) {
	if limit <= 0 {
		limit = consts.PageSize
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	var userLogins = make([]*userLogin, 0)
	count, err := _Engine.Desc("create_time").Where(cond).Limit(limit, start).FindAndCount(&userLogins)
	if err != nil {
		return 0, nil, err
	}
	return count, userLogins, nil
}

// 获取用户列表
func getUsers(cond builder.Cond, page, limit int) (int64, []*user, error) {
	if limit <= 0 {
		limit = consts.PageSize
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	var users = make([]*user, 0)
	count, err := _Engine.Omit("id", "password", "salt").Desc("create_time").Where(cond).Limit(limit, start).FindAndCount(&users)
	if err != nil {
		return 0, nil, err
	}
	return count, users, nil
}

// 根据用户ID获取用户信息
func getUserByID(userID common.ID) (*user, error) {
	return getUser(builder.Eq{"user_id": userID})
}

// 根据用户名获取用户信息
func getUserByName(name string) (*user, error) {
	return getUser(builder.Eq{"name": name})
}

// 根据邮箱获取用户信息
func getUserByEmail(email string) (*user, error) {
	return getUser(builder.Eq{"email": email})
}

// 根据手机号获取用户信息
func getUserByMobile(mobile string) (*user, error) {
	return getUser(builder.Eq{"mobile": mobile})
}

// 获取第三方绑定的账号
func getUserByTypeAndOpenID(t string, openID string) (*user, error) {
	third := new(userThird)
	has, err := _Engine.Where("type = ? AND open_id = ? AND status = ?", t, openID, consts.Normal).Get(third)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.ErrUserNotExist
	}
	user, err := getUserByID(third.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 用户数量
func getUserCount(cond builder.Cond) (int64, error) {
	user := new(user)
	return _Engine.Where(cond).Count(user)
}

// 获取用户信息
func getUser(cond builder.Cond) (*user, error) {
	user := new(user)
	has, err := _Engine.Where(cond).Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.ErrUserNotExist
	}
	return user, nil
}

// 根据用户ID获取用户详细信息
func getUserInfoByID(userID common.ID) (*userInfo, error) {
	userInfo := new(userInfo)
	has, err := _Engine.Where("user_id = ?", userID).Get(userInfo)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.ErrUserNotExist
	}
	return userInfo, nil
}

// 更新登录错误次数
func updateLoginErrorForUser(userID common.ID, numb int) error {
	session := _Engine.NewSession()
	defer session.Close()
	var err error
	if err = session.Begin(); err != nil {
		return err
	}
	user := &user{LastErrorTime: common.Now()}
	if numb == 0 {
		user.Error = 0
		_, err = session.Cols("error", "last_error_time").Where("user_id = ?", userID).Update(user)
	} else {
		_, err = session.Incr("error").Cols("error", "last_error_time").Where("user_id = ?", userID).Update(user)
	}
	if err != nil {
		return err
	}
	return session.Commit()
}

// 更新用户详细信息
func updateMy(userID common.ID, userInfo *userInfo) error {
	return updateUserInfo(userID, userInfo, "nickname", "gender", "qq", "weixin")
}

// 更新头像
func updateAvatarForUser(userID common.ID, avatar string) error {
	userInfo := &userInfo{Avatar: avatar}
	return updateUserInfo(userID, userInfo, "avatar")
}

// 更新昵称
func updateNicknameForUser(userID common.ID, nickname string) error {
	userInfo := &userInfo{Nickname: nickname}
	return updateUserInfo(userID, userInfo, "nickname")
}

// 更新性别
func updateGenderForUser(userID common.ID, gender consts.Gender) error {
	userInfo := &userInfo{Gender: gender}
	return updateUserInfo(userID, userInfo, "gender")
}

// 更新密码
func updatePasswordForUser(userID common.ID, password string) error {
	salt := common.RandomString(6)
	password = common.MD5(password + salt)
	user := &user{Password: password, Salt: salt, UpdateTime: common.Now()}
	return updateUser(userID, user, "salt", "password", "update_time")
}

// 更新禁用状态
func updateForbiddenForUser(userID common.ID, forbidden consts.Forbidden) error {
	user := &user{Forbidden: forbidden, ForbiddenTime: common.Now()}
	return updateUser(userID, user, "forbidden", "forbidden_time")
}

// 激活用户
func updateActivateForUser(userID common.ID) error {
	user := &user{Activate: consts.Activated, ActivateCode: "", ActivateTime: common.Now()}
	return updateUser(userID, user, "activate", "activate_code", "activate_time")
}

// 重新生成激活码
func updateActivateCodeForUser(userID common.ID, activateCode string) error {
	user := &user{ActivateCode: activateCode}
	return updateUser(userID, user, "activate_code")
}

// 更新登录信息
func updateLoginForUser(userID common.ID, ip string) error {
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	user := user{Error: 0}
	if _, err := session.Cols("error").Where("user_id = ?", userID).Update(&user); err != nil {
		return err
	}
	userLogin := userLogin{UserID: userID, IP: ip, CreateTime: common.Now()}
	if region, err := region.IP2Region(ip); err == nil {
		userLogin.City = region.City
		userLogin.Country = region.Country
		userLogin.Province = region.Province
		userLogin.Region = region.Region
	}
	if _, err := session.Insert(&userLogin); err != nil {
		return err
	}
	return session.Commit()
}

// 更新手机号
func updateMobileForUser(userID common.ID, mobile string) error {
	user := &user{Mobile: mobile, UpdateTime: common.Now()}
	return updateUser(userID, user, "mobile", "update_time")
}

// 更新邮箱
func updateEmailForUser(userID common.ID, email string) error {
	user := &user{Email: email, UpdateTime: common.Now()}
	return updateUser(userID, user, "email", "update_time")
}

// 更新用户信息
func updateUser(userID common.ID, user *user, columns ...string) error {
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	if _, err := session.Cols(columns...).Where("user_id = ?", userID).Update(user); err != nil {
		return err
	}
	return session.Commit()
}

func updateUserInfo(userID common.ID, userInfo *userInfo, columns ...string) error {
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	if _, err := session.Cols(columns...).Where("user_id = ?", userID).Update(userInfo); err != nil {
		return err
	}
	user := &user{UpdateTime: common.Now()}
	if _, err := session.Cols("update_time").Where("user_id = ?", userID).Update(user); err != nil {
		return err
	}
	return session.Commit()
}

func createUser(user *user, userInfo *userInfo, userThird *userThird) error {
	if user == nil || userInfo == nil {
		return errors.ErrArgument
	}
	uid, err := _IDWorker.Next()
	if err != nil {
		return err
	}

	// 默认时间处理
	nullDate := common.DateTime(time.Date(1, 1, 1, 0, 0, 0, 0, time.Local))
	user.CreateTime = common.Now()
	user.ForbiddenTime = nullDate
	user.LastErrorTime = nullDate
	user.ActivateTime = nullDate
	user.UpdateTime = nullDate

	user.UserID = uid

	user.Salt = common.RandomString(6)
	// 无密码注册，自动生成密码
	if user.Password == "" {
		user.Password = common.RandomNumber(12)
	}
	user.Password = common.MD5(user.Password + user.Salt)
	user.Status = consts.Normal
	user.Forbidden = consts.Available
	if user.Activate != consts.UnActivated {
		// 非需激活用户设置激活状态、时间和最后一次登录时间
		user.Activate = consts.Activated
		user.ActivateTime = user.CreateTime
	}

	// 无内容必填字段处理
	if user.Email == "" {
		user.Email = user.UserID.Str()
	}
	if user.Mobile == "" {
		user.Mobile = user.UserID.Str()
	}
	if user.Name == "" {
		user.Name = user.UserID.Str()
	}

	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	if _, err := session.Insert(user); err != nil {
		return err
	}

	// 用户详情
	userInfo.UserID = user.UserID
	if _, err := session.Insert(userInfo); err != nil {
		return err
	}

	// 三方注册
	if userThird != nil {
		userThird.UserID = user.UserID
		userThird.Status = consts.Normal
		userThird.CreateTime = user.CreateTime
		userThird.UpdateTime = nullDate
		if _, err := session.Insert(userThird); err != nil {
			return err
		}
	}

	userLogin := userLogin{UserID: user.UserID, IP: userInfo.IP, CreateTime: common.Now()}
	if region, err := region.IP2Region(userInfo.IP); err == nil {
		userLogin.City = region.City
		userLogin.Country = region.Country
		userLogin.Province = region.Province
		userLogin.Region = region.Region
	}

	// 注册即登录
	if _, err := session.Insert(&userLogin); err != nil {
		logger.Error(userLogin)
		return err
	}

	return session.Commit()
}
