package model

// LoginUserRes 登录返回
type LoginUserRes struct {
	Id           uint64 `orm:"id,primary"       json:"id"`         //
	UserName     string `orm:"username,unique" json:"userName"`    // 用户名
	Mobile       string `orm:"phone" json:"phone"`                 //手机号
	UserNickname string `orm:"nick_name"    json:"userNickname"`   // 用户昵称
	UserPassword string `orm:"password"    json:"-"`               // 登录密码;cmf_password加密
	UserSalt     string `orm:"user_salt"        json:"-"`          // 加密盐
	UserStatus   uint   `orm:"user_status"      json:"userStatus"` // 用户状态;0:禁用,1:正常,2:未验证
	IsAdmin      int    `orm:"is_admin"         json:"isAdmin"`    // 是否后台管理员 1 是  0   否
}
