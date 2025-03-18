package initalize

func InitFilter() {
	// 加载本地敏感词字典
	//err := global.SensitiveFilter.LoadWordDict("utils/filter_swords/dict.txt")
	//if err != nil {
	//	zlog.Errorf("Failed to load word dict: %v", err)
	//}
	// 也可以加载网络敏感词字典
	// err = global.SensitiveFilter.LoadNetWordDict(" ")
	// if err != nil {
	//     zlog.Errorf("Failed to load net word dict: %v", err)
	// }
}
