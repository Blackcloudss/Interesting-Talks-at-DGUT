package model

// @Title        test.go
// @Description
// @Create       XdpCs 2025-02-16 上午1:31
// @Update       XdpCs 2025-02-16 上午1:31

// Test 一个实体类，对应数据库一张表
type Test struct {
	CommonModel
	UserID string
	Name   string
	Age    int
}

func (t *Test) TableName() {

}
