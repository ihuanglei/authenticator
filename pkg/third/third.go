package third

// User .
type User struct {
	OpenID   string `json:"open_id"`  // id
	TP       string `json:"tp"`       // 类型
	Nickname string `json:"nickname"` // 昵称
	Mobile   string `json:"mobile"`   // 手机号
	Avatar   string `json:"avatar"`   // 头像
	Gender   string `json:"gender"`   // 性别
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	Ext      string `json:"ext"`      //扩展字段
}

// Third .
type Third interface {
	GetAuthorizeURL(state string) string
	GetUser(code string) (*User, error)
	GetType() string
}
