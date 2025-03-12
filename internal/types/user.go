package types

// @Title        tz_user.go
// @Description
// @Create       XdpCs 2025-03-09 下午4:49
// @Update       XdpCs 2025-03-09 下午4:49

// 获取前端传输的 昵称，头像  -- 前端调用 wx.getUserProfile() 获取
type UserInfoReq struct {
	Nickname string `json:"nickname"` //  微信默认昵称
	Avatar   string `json:"avatar"`   //  微信默认头像
}

type UserInfoResp struct {
}

// 用户基本信息 入参
type CommonProfileReq struct {
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Sex      string `json:"sex"`      // 性别
	Birthday string `json:"birthday"` // 生日
	Sign     string `json:"sign"`     // 签名
	Tag      string `json:"tag"`      // 标签
}

// 用户基本信息 出参
type CommonProfileResp struct {
}

type DetailProfileReq struct {
	Name      string `json:"name"`      // 真实姓名
	StudentId string `json:"studentId"` // 学号
	Academy   string `json:"academy"`   // 学院
	Grade     int64  `json:"grade"`     // 年级
	Major     string `json:"major"`     // 专业
	Email     string `json:"email"`     // 邮箱
	//Phone     string `json:"phone"` -- 手机号 后端自动获取
}
