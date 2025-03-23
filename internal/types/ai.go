package types

// @Title        ai.go
// @Description
// @Create       XdpCs 2025-03-06 上午1:27
// @Update       XdpCs 2025-03-06 上午1:27
//
// AI聊天请求体
type AIChatStreamReq struct {
	Type            string  `json:"type"`                   // 类型
	Content         string  `form:"content" json:"content"` // 对话内容
	PresencePenalty float64 `json:"presence_penalty"`       // 可能值  介于 -2.0 和 2.0 之间的数字。如果该值为正，
	// 那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚，从而增加模型谈论新主题的可能性。 降低模型重复相同内容的可能性
}

// API请求结构体
type AIAPIReq struct {
	Messages        []Message `json:"messages"`         //对话消息列表
	Model           string    `json:"model"`            //模型名称
	PresencePenalty float64   `json:"presence_penalty"` //可能值
	Temperature     float64   `json:"temperature"`      //温度
	ResponseFormat  Format    `json:"response_format"`  //响应格式
	Stream          bool      `json:"stream"`           //是否流式输出
}

// 对话消息结构体
type Message struct {
	Role    string `json:"role"` // 角色 ： 1.系统、 2.用户、 3.助手
	Content string `json:"content"`
}

// 响应格式  1.text 2.json_object
type Format struct {
	Type string `json:"type"`
}

// API响应结构体
type AIChatStreamResp struct {
	Choices []Choice `json:"choices"` //模型生成的 completion 的选择列表
}

// 模型生成的 completion 的选择列表
type Choice struct {
	Delta struct {
		Content string `json:"content"` // completion 增量的内容
		Role    string `json:"role"`    // 角色
	} `json:"delta"`
}

// 删除历史聊天记录 入参
type DeleteHistoryMessagesReq struct {
	Type string `json:"type"` //类型
}

// 删除历史聊天记录 出参
type DeleteHistoryMessagesResp struct {
}
