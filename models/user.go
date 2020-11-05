package models

import (
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/convert"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/common"

	"xorm.io/builder"
)

// CreateUserWithThird 三方注册登录
func CreateUserWithThird(register *st.RegisterDto) (common.ID, error) {
	has, err := HasUserByThird(register.TP, register.OpenID)
	if err != nil {
		return 0, err
	}
	if has {
		return 0, errors.ErrUserExist
	}
	user := &user{Mode: consts.Third, Mobile: register.Mobile}
	userInfo := &userInfo{
		Nickname: register.Nickname,
		Avatar:   register.Avatar,
		Province: register.Province,
		City:     register.City,
		IP:       register.IP,
	}
	thirdUser := &userThird{OpenID: register.OpenID, Type: register.TP}
	if err := createUser(user, userInfo, thirdUser); err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// CreateUserWithName 用户名密码创建用户
func CreateUserWithName(register *st.RegisterDto) (common.ID, error) {
	name := common.Trim(register.LoginName)
	// 必须包含一个字母或下划线
	if b, err := regexp.MatchString(`^\w{5,20}$`, name); !b || err != nil {
		return 0, errors.ErrName
	}
	if b, err := regexp.MatchString(`^\d{5,20}$`, name); b || err != nil {
		return 0, errors.ErrName
	}

	if !common.IsSimplePassword(register.Password) {
		return 0, errors.ErrPassword
	}
	has, err := HasUserByName(name)
	if err != nil {
		return 0, err
	}
	if has {
		return 0, errors.ErrUserExist
	}
	user := &user{Name: name, Password: register.Password, Mode: consts.Name}
	userInfo := &userInfo{IP: register.IP}
	if err := createUser(user, userInfo, nil); err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// CreateUserWithEmail 邮箱密码创建用户
func CreateUserWithEmail(register *st.RegisterDto) (common.ID, string, error) {
	if !common.IsEmail(register.Email) {
		return 0, "", errors.ErrEmail
	}
	if !common.IsSimplePassword(register.Password) {
		return 0, "", errors.ErrPassword
	}
	has, err := HasUserByEmail(register.Email)
	if err != nil {
		return 0, "", err
	}
	if has {
		return 0, "", errors.ErrUserEmailExist
	}
	activateCode := common.MD5(fmt.Sprintf("%v%s", time.Now().UnixNano(), common.RandomString(24)))
	user := &user{Email: register.Email, Password: register.Password, Activate: consts.UnActivated, ActivateCode: activateCode, Mode: consts.Email}
	userInfo := &userInfo{IP: register.IP}
	if err := createUser(user, userInfo, nil); err != nil {
		return 0, "", err
	}
	return user.UserID, activateCode, nil
}

// CreateUserWithMobile 创建用户
func CreateUserWithMobile(register *st.RegisterDto) (common.ID, error) {
	if !common.IsMobile(register.Mobile) {
		return 0, errors.ErrMobile
	}
	if !common.IsSimplePassword(register.Password) {
		return 0, errors.ErrPassword
	}
	has, err := HasUserByMobile(register.Mobile)
	if err != nil {
		return 0, err
	}
	if has {
		return 0, errors.ErrUserMobileExist
	}
	user := &user{Mobile: register.Mobile, Mode: consts.Mobile}
	userInfo := &userInfo{IP: register.IP}
	if err := createUser(user, userInfo, nil); err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// ActivateUser 激活用户
func ActivateUser(userID common.ID, activateCode string) error {
	user, err := getUserByID(userID)
	if err != nil {
		return err
	}
	if user.IsActivated() {
		return errors.ErrUserAlreadyActivate
	}
	if user.ActivateCode != activateCode {
		return errors.ErrActiveCode
	}
	return updateActivateForUser(userID)
}

// ActivateUserWithOutCode 直接激活用户
func ActivateUserWithOutCode(userID common.ID) error {
	user, err := getUserByID(userID)
	if err != nil {
		return err
	}
	if user.IsActivated() {
		return errors.ErrUserAlreadyActivate
	}
	return updateActivateForUser(userID)
}

// UpdateActivateCodeForUser 重新发送激活码
func UpdateActivateCodeForUser(userID common.ID) (string, error) {
	user, err := getUserByID(userID)
	if err != nil {
		return "", err
	}
	if user.IsActivated() {
		return "", errors.ErrUserAlreadyActivate
	}
	activateCode := common.MD5(fmt.Sprintf("%v%s", time.Now().UnixNano(), common.RandomString(24)))
	return activateCode, updateActivateCodeForUser(userID, activateCode)
}

// UpdateAvatarForUser 更新头像
func UpdateAvatarForUser(userID common.ID, avatar string) error {
	if common.IsEmpty(avatar) {
		return errors.ErrAvatar
	}
	if _, err := getUserByID(userID); err != nil {
		return err
	}
	return updateAvatarForUser(userID, avatar)
}

// UpdateNicknameForUser 更新昵称
func UpdateNicknameForUser(userID common.ID, nickname string) error {
	if common.IsEmpty(nickname) || utf8.RuneCountInString(nickname) > 15 {
		return errors.ErrNickname
	}
	if _, err := getUserByID(userID); err != nil {
		return err
	}
	return updateNicknameForUser(userID, nickname)
}

// UpdateGenderForUser 更新性别
func UpdateGenderForUser(userID common.ID, gender consts.Gender) error {
	_, err := getUserByID(userID)
	if err != nil {
		return err
	}
	return updateGenderForUser(userID, gender)
}

// UpdateEmailForUser 更新邮箱
func UpdateEmailForUser(userID common.ID, email string) error {
	if !common.IsEmail(email) {
		return errors.ErrEmail
	}
	if _, err := getUserByID(userID); err != nil {
		return err
	}
	has, err := HasUserByEmail(email)
	if err != nil {
		return err
	}
	if has {
		return errors.ErrUserEmailExist
	}
	return updateEmailForUser(userID, email)
}

// UpdateMobileForUser 更新手机号
func UpdateMobileForUser(userID common.ID, mobile string) error {
	if !common.IsMobile(mobile) {
		return errors.ErrMobile
	}
	if _, err := getUserByID(userID); err != nil {
		return err
	}
	has, err := HasUserByMobile(mobile)
	if err != nil {
		return err
	}
	if has {
		return errors.ErrUserMobileExist
	}
	return updateMobileForUser(userID, mobile)
}

// UpdateMy 更新信息
func UpdateMy(userID common.ID, userInfoDto *st.UserInfoDto) error {
	if _, err := getUserByID(userID); err != nil {
		return err
	}
	userInfo := new(userInfo)
	if err := convert.Map(userInfoDto, userInfo); err != nil {
		return err
	}
	return updateMy(userID, userInfo)
}

// UpdatePasswordForUser 修改密码
func UpdatePasswordForUser(userID common.ID, oldPassword, password string) error {
	if !common.IsSimplePassword(password) {
		return errors.ErrPassword
	}
	user, err := getUserByID(userID)
	if err != nil {
		return err
	}
	if common.MD5(oldPassword+user.Salt) != user.Password {
		return errors.ErrInvalidPassword
	}
	return updatePasswordForUser(userID, password)
}

// UpdatePassword1ForUser 直接修改密码
func UpdatePassword1ForUser(userID common.ID, password string) error {
	if !common.IsSimplePassword(password) {
		return errors.ErrPassword
	}
	return updatePasswordForUser(userID, password)
}

// UpdateForbiddenForUser 更新禁用状态
func UpdateForbiddenForUser(userID common.ID, forbidden consts.Forbidden) error {
	user, err := getUserByID(userID)
	if err != nil {
		return err
	}
	if user.Forbidden == forbidden {
		return nil
	}
	return updateForbiddenForUser(userID, forbidden)
}

// UpdateLoginErrorForUser 增加登录错误次数
func UpdateLoginErrorForUser(userID common.ID) error {
	user, err := getUserByID(userID)
	if err != nil {
		return err
	}
	if user.Error >= consts.LoginErrorCount {
		return nil
	}
	return updateLoginErrorForUser(userID, 1)
}

// ResetLoginErrorForUser 重新设置错误次数为0
func ResetLoginErrorForUser(userID common.ID) error {
	user, err := getUserByID(userID)
	if err != nil {
		return err
	}
	if user.Error == 0 {
		return nil
	}
	return updateLoginErrorForUser(userID, 0)
}

// GetUserByMobile 根据手机号查询用户信息
func GetUserByMobile(mobile string) (*st.UserDto, error) {
	if !common.IsMobile(mobile) {
		return nil, errors.ErrMobile
	}
	user, err := getUserByMobile(mobile)
	if err != nil {
		return nil, err
	}
	userDto := new(st.UserDto)
	err = convert.Map(user, userDto)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}

// GetUserByEmail 根据邮箱获取用户信息
func GetUserByEmail(email string) (*st.UserDto, error) {
	if !common.IsEmail(email) {
		return nil, errors.ErrEmail
	}
	user, err := getUserByEmail(email)
	if err != nil {
		return nil, err
	}
	userDto := new(st.UserDto)
	err = convert.Map(user, userDto)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}

// GetUserByID 根据用户ID获取用户信息
func GetUserByID(userID common.ID) (*st.UserDto, error) {
	user, err := getUserByID(userID)
	if err != nil {
		return nil, err
	}
	userDto := new(st.UserDto)
	err = convert.Map(user, userDto)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}

// GetUserInfoByID 根据用户ID获取用户详细信息
func GetUserInfoByID(userID common.ID) (*st.UserInfoDto, error) {
	userInfo, err := getUserInfoByID(userID)
	if err != nil {
		return nil, err
	}
	userInfoDto := new(st.UserInfoDto)
	err = convert.Map(userInfo, userInfoDto)
	if err != nil {
		return nil, err
	}
	return userInfoDto, nil
}

// GetUsers 获取用户列表
func GetUsers(query st.UserQuery) (int64, []*st.UserDto, error) {
	page := query.Page
	limit := query.Limit

	cond := builder.NewCond()

	keyword := query.Keyword
	if common.Trim(keyword) != "" {
		keyword = keyword + "%"
		cond = builder.Or(builder.Like{"mobile", keyword}, builder.Like{"email", keyword}, builder.Like{"name", keyword})
	}

	if common.Trim(query.Name) != "" {
		cond = cond.And(builder.Like{"name", query.Name + "%"})
	}
	if common.Trim(query.Email) != "" {
		cond = cond.And(builder.Like{"email", query.Email + "%"})
	}
	if common.Trim(query.Mobile) != "" {
		cond = cond.And(builder.Like{"mobile", query.Mobile + "%"})
	}

	genderValues := query.Genders
	if len(genderValues) > 0 {
		var genders []consts.Gender
		for _, value := range genderValues {
			genders = append(genders, consts.NewGender(value))
		}
		cond = cond.And(builder.In("gender", genders))
	}

	forbiddenValues := query.Forbiddens
	if len(forbiddenValues) > 0 {
		var forbiddens []consts.Forbidden
		for _, value := range forbiddenValues {
			forbiddens = append(forbiddens, consts.NewForbidden(value))
		}
		cond = cond.And(builder.In("forbidden", forbiddens))
	}

	activateValues := query.Activates
	if len(activateValues) > 0 {
		var activates []consts.Activate
		for _, value := range activateValues {
			activates = append(activates, consts.NewActivate(value))
		}
		cond = cond.And(builder.In("activate", activates))
	}

	count, users, err := getUsers(cond, page, limit)
	if err != nil {
		return 0, nil, err
	}
	var usersDtos = make([]*st.UserDto, len(users))
	err = convert.Map(&users, &usersDtos)
	if err != nil {
		return 0, nil, err
	}
	return count, usersDtos, nil
}

// GetUserLoginByID 根据用户ID获取用户登录历史
func GetUserLoginByID(userID common.ID, query st.EmptyQuery) (int64, []*st.UserLoginDto, error) {
	page := query.Page
	limit := query.Limit
	cond := builder.Eq{"user_id": userID}
	count, userLogins, err := getUserLogins(cond, page, limit)
	if err != nil {
		return 0, nil, err
	}
	var userLoginDtos = make([]*st.UserLoginDto, len(userLogins))
	if err := convert.Map(&userLogins, &userLoginDtos); err != nil {
		return 0, nil, err
	}
	return count, userLoginDtos, nil
}

// HasUserByEmail 是否存在有效邮箱
func HasUserByEmail(email string) (bool, error) {
	if !common.IsEmail(email) {
		return false, errors.ErrEmail
	}
	count, err := getUserCount(builder.Eq{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasUserByMobile 是否存在有效手机号
func HasUserByMobile(mobile string) (bool, error) {
	if !common.IsMobile(mobile) {
		return false, errors.ErrMobile
	}
	count, err := getUserCount(builder.Eq{"mobile": mobile})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasUserByName 是否存在有效用户名
func HasUserByName(name string) (bool, error) {
	count, err := getUserCount(builder.Eq{"name": name})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasUserByThird 是否存在第三方绑定
func HasUserByThird(tp, openID string) (bool, error) {
	_, err := getUserByTypeAndOpenID(tp, openID)
	if err != nil {
		if e, ok := err.(errors.Error); ok && e.Code() == errors.ErrUserNotExist.Code() {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckUserByMobile 根据手机号检查用户状态
func CheckUserByMobile(mobile string) error {
	if !common.IsMobile(mobile) {
		return errors.ErrMobile
	}
	user, err := getUserByMobile(mobile)
	if err != nil {
		return err
	}
	if err := checkState(user); err != nil {
		return err
	}
	return nil
}

// Login 登录
func Login(loginDto *st.LoginDto) (*st.UserDto, error) {
	type loginType int
	const (
		ltName loginType = 1 << iota
		ltEmail
		ltMobile
	)

	guessLoginType := ltName

	var user *user
	var err error

	if common.IsEmail(loginDto.LoginName) {
		guessLoginType = ltEmail
	}

	if common.IsMobile(loginDto.LoginName) {
		guessLoginType = ltMobile
	}

	if guessLoginType&ltEmail == ltEmail {
		user, err = getUserByEmail(loginDto.LoginName)
		if err == nil {
			goto check
		}
	}

	if guessLoginType&ltMobile == ltMobile {
		user, err = getUserByMobile(loginDto.LoginName)
		if err == nil {
			goto check
		}
	}

	if guessLoginType&ltName == ltName {
		user, err = getUserByName(loginDto.LoginName)
		if err == nil {
			goto check
		}
	}

	if err != nil {
		return nil, errors.ErrUserNotExist
	}

check:

	if err := checkState(user); err != nil {
		return nil, err
	}

	if common.MD5(loginDto.Password+user.Salt) != user.Password {
		UpdateLoginErrorForUser(user.UserID)
		return nil, errors.ErrInvalidPassword
	}
	return login(user, loginDto.IP)
}

// LoginByMobile 手机验证码登录
func LoginByMobile(loginDto *st.LoginDto) (*st.UserDto, error) {
	user, err := getUserByMobile(loginDto.Mobile)
	if err != nil {
		return nil, errors.ErrUserNotExist
	}
	return login(user, loginDto.IP)
}

// LoginByOpenID 根据用户OpenID获取用户信息
func LoginByOpenID(loginDto *st.LoginDto) (*st.UserDto, error) {
	user, err := getUserByTypeAndOpenID(loginDto.Type, loginDto.OpenID)
	if err != nil {
		return nil, err
	}
	return login(user, loginDto.IP)
}

func login(user *user, ip string) (*st.UserDto, error) {
	if err := updateLoginForUser(user.UserID, ip); err != nil {
		return nil, err
	}
	userDto := new(st.UserDto)
	if err := convert.Map(user, userDto); err != nil {
		return nil, err
	}
	return userDto, nil
}

func checkState(user *user) error {
	if !user.canLogin(consts.LoginErrorCount) {
		return errors.ErrUserLocked
	}
	if user.IsForbidden() {
		return errors.ErrUserForbidden
	}
	if user.IsDelete() {
		return errors.ErrUserNotExist
	}
	if !user.IsActivated() {
		return errors.ErrUserUnActivated
	}
	return nil
}
