package utils

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"time"
)

// @Title        time.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:34
// @Update       XdpCs 2025-02-24 下午11:34
func RecordTime(start time.Time) func() {
	return func() {
		end := time.Now()
		zlog.Debugf("use time:%d", end.Unix()-start.Unix())
	}
}
