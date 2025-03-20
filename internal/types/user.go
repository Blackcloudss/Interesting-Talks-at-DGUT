package types

import "mime/multipart"

// @Title        middle.go
// @Description
// @Create       XdpCs 2025-03-09 下午4:49
// @Update       XdpCs 2025-03-09 下午4:49

// 前端调用 wx.getUserProfile() 获取 用户加密信息 传给后端
type UserInfoReq struct {
	RawData       string `json:"rawData"`       //不包括敏感信息的原始数据字符串，用于计算签名
	Signature     string `json:"signature"`     //使用 sha1( rawData + sessionkey ) 得到字符串，用于校验用户信息
	EncryptedData string `json:"encryptedData"` //包括敏感数据在内的完整用户信息的加密数据
	Iv            string `json:"iv"`            //加密算法的初始向量
}

// 后端解密后获取昵称，头像
type UserInfo struct {
	Nickname string `json:"nickname"` //  微信默认昵称
	Avatar   string `json:"avatar"`   //  微信默认头像
}

type UserInfoResp struct {
}

// 获取用户基本信息 入参
type GetCommonProfileReq struct {
}

// 获取用户基本信息 出参
type GetCommonProfileResp struct {
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Tag      string `json:"tag"`      // 标签
	Sex      string `json:"sex"`      // 性别
	Birthday string `json:"birthday"` // 生日
	Sign     string `json:"sign"`     // 签名
}

// 获取他人基本信息 入参
type GetOtherProfileReq struct {
	OtherID int64 `json:"otherid"` // 他人ID
}

// 获取他人基本信息 出参
type GetOtherProfileResp struct {
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Tag      string `json:"tag"`      // 标签
	Sex      string `json:"sex"`      // 性别
	Birthday string `json:"birthday"` // 生日
	Sign     string `json:"sign"`     // 签名
}

// 更改用户基本信息 入参
type UpdateCommonProfileReq struct {
	Nickname string                `json:"nickname"` // 昵称
	Avatar   *multipart.FileHeader `form:"avatar"`   // 头像图片文件
	Sex      string                `json:"sex"`      // 性别
	Birthday string                `json:"birthday"` // 生日
	Sign     string                `json:"sign"`     // 签名
}

// 更改用户基本信息 出参
type UpdateCommonProfileResp struct {
}

// 获取用户私人信息 入参
type GetPrivateProfileReq struct {
}

// 获取用户私人信息 出参
type GetPrivateProfileResp struct {
	Name      string `json:"name"`      // 真实姓名
	StudentId string `json:"studentId"` // 学号
	Academy   string `json:"academy"`   // 学院
	Grade     int64  `json:"grade"`     // 年级
	Major     string `json:"major"`     // 专业
	Phone     string `json:"phone"`     // 手机号
}

// 更改用户隐私信息 入参
type UpdatePrivateProfileReq struct {
	Name      string `json:"name"`      // 真实姓名
	StudentId string `json:"studentId"` // 学号
	Academy   string `json:"academy"`   // 学院
	Grade     int64  `json:"grade"`     // 年级
	Major     string `json:"major"`     // 专业
	Phone     string `json:"phone"`     // 手机号
}

// 更改用户隐私信息 出参
type UpdatePrivateProfileResp struct {
	Role string `json:"role"` // 用户身份
}

// 更改其他用户身份 入参
type UpdateOtherRoleReq struct {
	OtherID int64 `json:"otherid"` //其他用户ID
}

// 更改其他用户身份 出参
type UpdateOtherRoleResp struct {
	Role string `json:"role"` // 当前用户身份
}
