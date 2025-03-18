package initalize

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
)

func InitFilter() {
	// 加载本地敏感词字典
	err := global.Filter.LoadWordDict("utils/filter_swords/dict.txt")
	if err != nil {
		zlog.Errorf("加载敏感词库失败: %v", err)
	}
	zlog.Infof("敏感词库加载成功成功！")

	// 也可以加载网络敏感词字典
	// err = global.SensitiveFilter.LoadNetWordDict(" ")
	// if err != nil {
	//     zlog.Errorf("Failed to load net word dict: %v", err)
	// }
}
