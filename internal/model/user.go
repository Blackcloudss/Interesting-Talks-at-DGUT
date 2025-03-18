package model

import "time"

// @Title        user.go
// @Description
// @Create       XdpCs 2025-03-10 上午1:32
// @Update       XdpCs 2025-03-10 上午1:32

// 用户展示表
type UserDisplay struct {
	CommonModel        // id 为 用户id
	OpenId      string `gorm:"column:openid;type:varchar(255);uniqueIndex"`                 // 添加唯一索引                     // 用户在不同类型产品中的身份id，不同产品互不相同
	Role        string `gorm:"column:role;type:varchar(20);default:tourist;comment:'用户身份'"` //  角色身份 默认身份为游客
	//建立联合索引
	Avatar   string `gorm:"column:avatar;type:varchar(255);index:idx_profile,priority:3;comment:'微信头像'"`               // 头像
	Nickname string `gorm:"column:nickname;type:varchar(50);default:微信用户;index:idx_profile,priority:2;comment:'微信昵称'"` // 昵称
	Tag      string `gorm:"column:tag;type:varchar(50);default:大一新生;index:idx_profile,priority:1;comment:'用户标签'"`      // 用户标签
}

func (t *UserDisplay) TableName() string {
	return "user_display"
}

// 用户基本信息表
type UserCommon struct {
	CommonModel
	Sex      string    `gorm:"column:sex;type:char(2);comment:'性别'"`
	Birthday time.Time `gorm:"column:birthday;type:date;comment:'出生日期'"` //****年**月**日
	Sign     string    `gorm:"column:sign;type:varchar(50);comment:'个性签名'"`

	UserID int64 `gorm:"column:user_id;type:bigint;comment:'用户ID'"`
	//外键关联
	UserDisplay UserDisplay `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *UserCommon) TableName() string {
	return "user_common"
}

// 用户私人信息表
type UserPrivate struct {
	CommonModel
	Name      string `gorm:"column:name;type:varchar(30);comment:'真实姓名'"`
	StudentId string `gorm:"column:student_id;type:char(13);comment:'学号'"`
	Academy   string `gorm:"column:academy;type:varchar(30);comment:'学院'"`
	Grade     int64  `gorm:"column:grade;type:bigint;comment:'年级'"`
	Major     string `gorm:"column:major;type:varchar(30);comment:'专业'"`
	Phone     string `gorm:"column:phone;type:char(11);comment:'手机号'"`

	UserID int64 `gorm:"column:user_id;type:bigint;comment:'用户ID'"`
	//外键关联
	UserDisplay UserDisplay `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *UserPrivate) TableName() string {
	return "user_private"
}

// 用户授权表
type UserAuth struct {
	CommonModel
	Is_agreed_phone    bool `gorm:"column:is_agreed_phone;type:tinyint(1);comment:'是否同意授权手机号'"`     // 0 未授权 1 已授权
	Is_agreed_avatar   bool `gorm:"column:is_agreed_avatar;type:tinyint(1);comment:'是否同意授权微信头像'"`   // 0 未授权 1 已授权
	Is_agreed_nickname bool `gorm:"column:is_agreed_nickname;type:tinyint(1);comment:'是否同意授权微信昵称'"` // 0 未授权 1 已授权

	UserID int64 `gorm:"column:user_id;type:bigint;comment:'用户ID'"`
	//外键关联
	UserDisplay UserDisplay `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *UserAuth) TableName() string {
	return "user_auth"
}
