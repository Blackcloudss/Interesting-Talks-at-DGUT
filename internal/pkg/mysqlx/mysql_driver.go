package mysqlx

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/pkg/database"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
}

// @Title        mysql_driver.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:25
// @Update       XdpCs 2025-02-24 下午11:25
func (m *Mysql) InitDataBase(config configs.Config) (*gorm.DB, error) {
	dsn := m.GetDsn(config)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		zlog.Panicf("MySQL无法连接数据库！: %v", err)
		return nil, err
	}
	zlog.Infof("MySQL连接数据库成功！")
	return db, nil
}
func (m *Mysql) GetDsn(config configs.Config) string {
	return config.DB.Dsn
}
func NewMySql() database.DataBase {
	return &Mysql{}
}
