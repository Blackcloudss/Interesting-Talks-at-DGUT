package types

// @Title        deepseek.go
// @Description
// @Create       XdpCs 2025-03-06 上午1:27
// @Update       XdpCs 2025-03-06 上午1:27

// 定义API请求结构体
type ChatReq struct {
	Model            string    `json:"model"`             //模型名称
	Messages         []Message `json:"messages"`          //历史消息
	MaxTokens        int       `json:"max_tokens"`        //最大输出字数
	FrequencyPenalty float64   `json:"frequency_penalty"` //频率惩罚
	PresencePenalty  float64   `json:"presence_penalty"`  //存在惩罚
	Temperature      float64   `json:"temperature"`       //温度
	TopP             float64   `json:"top_p"`             //贪婪度
	Stream           bool      `json:"stream"`            //是否流式输出
	ResponseFormat   Format    `json:"response_format"`   //响应格式
	Stop             []string  `json:"stop"`              //停止词
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Format struct {
	Type string `json:"type"`
}

// 定义API响应结构体
type ChatResp struct {
	Id      string   `json:"id"`
	Choices []Choice `json:"choices""`
	Error   Error    `json:"error"`
}

type Choice struct {
	Message struct {
		Content string `json:"content"`
		Role    string `json:"role"`
	} `json:"message"`
}
type Error struct {
	Message string `json:"message"`
}
