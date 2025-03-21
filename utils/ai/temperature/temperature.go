package temperature

import "github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/ai/prompt"

// @Title        temperature.go
// @Description
// @Create       XdpCs 2025-03-21 下午7:45
// @Update       XdpCs 2025-03-21 下午7:45

/*
	Temperature设置
	代码生成/数学解题	0.0
	数据抽取/分析	1.0
	通用对话	1.3
	翻译	1.3
	创意类写作/诗歌创作	1.5
*/
// Temperature设置
const (
	CODE_TEMPERATURE      = 0.0 //代码生成/数学解题
	DATA_TEMPERATURE      = 1.0 //数据抽取/分析
	COMMON_TEMPERATURE    = 1.3 //通用对话
	TRANSLATE_TEMPERATURE = 1.3 //翻译
	POETRY_TEMPERATURE    = 1.5 //创意类写作/诗歌创作	1.5
)

func JudgeTemperature(Type string) float64 {
	switch Type {
	case prompt.OUTLINE:
		return DATA_TEMPERATURE
	case prompt.TRANSLATION:
		return TRANSLATE_TEMPERATURE
	case prompt.SLOGAN_PROMPT:
		return POETRY_TEMPERATURE
	default:

		return COMMON_TEMPERATURE
	}
}
