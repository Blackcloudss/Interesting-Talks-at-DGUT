package types

// @Title        deepseek.go
// @Description
// @Create       XdpCs 2025-03-06 上午1:27
// @Update       XdpCs 2025-03-06 上午1:27

// 定义API请求结构体
type ChatReq struct {
	Model    string    `json:"model"`    //模型名称
	Messages []Message `json:"messages"` //历史消息
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 定义API响应结构体
type ChatResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
