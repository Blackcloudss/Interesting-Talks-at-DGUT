package model

// @Title        usermsg.go
// @Description
// @Create       XdpCs 2025-03-10 上午1:32
// @Update       XdpCs 2025-03-10 上午1:32

type UserMsg struct {
	CommonModel
	OpenId   string `gorm:"open_id"`                                  // 用户在不同类型产品中的身份id，不同产品互不相同
	Nickname string `gorm:"default:游客" json:"nickname"`               // 默认昵称为“游客”
	Avatar   string `gorm:"default:default_avatar.png" json:"avatar"` // 默认头像

	Name      string `gorm:"column:name;type:varchar(30);comment:'真实姓名'"`
	StudentId string `gorm:"column:student_id;type:char(13);comment:'学号'"`
	Academy   string `gorm:"column:academy;type:varchar(30);comment:'学院'"`
	Grade     int64  `gorm:"column:grade;type:bigint;comment:'年级'"`
	Major     string `gorm:"column:major;type:varchar(30);comment:'专业'"`
	Phone     string `gorm:"column:phone;type:char(11);comment:'手机号'"`
	Email     string `gorm:"column:email;type:varchar(30);comment:'邮箱'"` // xxxx@email.com

	Sex      string `gorm:"column:sex;type:char(2);comment:'性别'"`
	Birthady string `gorm:"column:birthday;type:varchar(15);comment:'出生日期'"` //****年**月**日
	Sign     string `gorm:"column:sign;type:varchar(50);comment:'个性签名'"`
	Tag      string `gorm:"column:tag;type:varchar(50);comment:'用户标签'"`
}
