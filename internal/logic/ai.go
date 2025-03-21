package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/ai/prompt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/ai/role"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/ai/temperature"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

const (
	MaxHistoryRounds = 10 // 历史对话轮数 定义为 10轮
)

// @Title        ai.go
// @Description
// @Create       XdpCs 2025-03-09 下午2:27
// @Update       XdpCs 2025-03-09 下午2:27
type AILogic struct {
}

func NewAILogic() *AILogic {
	return &AILogic{}
}

// 获取API相应数据流
func (l *AILogic) GetChatStream(ctx context.Context, MessagesKey string, req types.AIChatStreamReq) (httpResp *http.Response, messages []types.Message, err error) {
	defer utils.RecordTime(time.Now())()
	//初始化或获取前十轮历史对话
	messages, err = l.GetHistoryMessages(ctx, MessagesKey, req.Type)

	//添加新的对话
	messages = append(messages, types.Message{
		Role:    role.USER,
		Content: req.Content,
	})

	// 构建请求体
	apiReq := l.NewApiReq(messages, req)

	//创建HTTP请求
	body, _ := json.Marshal(apiReq)
	httpReq, _ := http.NewRequest("POST", global.Config.AI.BaseUrl+global.Config.AI.FuncUrl, bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+global.Config.AI.ApiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	//发送HTTP请求 若服务器未在2分钟内响应，客户端会自动取消请求
	Client := &http.Client{Timeout: time.Minute * 2}
	httpResp, err = Client.Do(httpReq)
	if err != nil {
		zlog.CtxErrorf(ctx, "发送HTTP请求失败：%v", err)
		// 待补充
		return
	}
	return httpResp, messages, err
}

// 创建API请求体
func (l *AILogic) NewApiReq(messages []types.Message, req types.AIChatStreamReq) (apiReq types.AIAPIReq) {
	apiReq = types.AIAPIReq{
		Messages:        messages,
		Model:           global.Config.AI.Model,
		PresencePenalty: req.PresencePenalty,
		ResponseFormat:  types.Format{Type: "text"},
		Stream:          true,
		Temperature:     temperature.JudgeTemperature(req.Type),
	}
	return apiReq
}

// 获取历史对话
func (l *AILogic) GetHistoryMessages(ctx context.Context, MessagesKey string, Type string) (messages []types.Message, err error) {
	defer utils.RecordTime(time.Now())()
	// 获取系统消息
	system, err := global.Rdb.Get(ctx, MessagesKey+":system").Result()
	if err == redis.Nil {
		// 初始化系统消息
		switch Type {
		// 根据Type参数选择系统消息
		case prompt.OUTLINE:
			global.Rdb.Set(ctx, MessagesKey+":system", prompt.OUTLINE_PROMPT, 0)
			system = prompt.OUTLINE_PROMPT
		case prompt.SLOGAN:
			global.Rdb.Set(ctx, MessagesKey+":system", prompt.SLOGAN_PROMPT, 0)
			system = prompt.SLOGAN_PROMPT
		case prompt.TRANSLATION:
			global.Rdb.Set(ctx, MessagesKey+":system", prompt.TRANSLATION_PROMPT, 0)
			system = prompt.TRANSLATION_PROMPT
		case prompt.ITAD:
			global.Rdb.Set(ctx, MessagesKey+":system", prompt.ITAD_PROMPT, 0)
			system = prompt.ITAD_PROMPT
		case prompt.CIRNO:
			global.Rdb.Set(ctx, MessagesKey+":system", prompt.CIRNO_PROMPT, 0)
			system = prompt.CIRNO_PROMPT
		default:
			global.Rdb.Set(ctx, MessagesKey+":system", prompt.ITAD_PROMPT, 0)
			system = prompt.ITAD_PROMPT
		}
	} else if err != nil {
		return nil, err
	}

	// 获取对话记录
	msgs, err := global.Rdb.LRange(ctx, MessagesKey+":messages", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// 构造消息列表
	messages = []types.Message{{Role: role.SYSTEM, Content: system}}
	for _, msg := range msgs {
		var m types.Message
		if json.Unmarshal([]byte(msg), &m) == nil {
			messages = append(messages, m)
		}
	}
	return messages, nil
}

// 保存历史对话
func (l *AILogic) SaveHistoryMessages(ctx context.Context, MessagesKey string, messages []types.Message) (err error) {
	defer utils.RecordTime(time.Now())()
	// 排除系统消息
	var toSave []string
	for _, msg := range messages[1:] {
		data, _ := json.Marshal(msg)
		toSave = append(toSave, string(data))
	}

	// 保存并修剪历史
	pipe := global.Rdb.Pipeline()
	// 追加新消息到列表
	pipe.RPush(ctx, MessagesKey+":messages", toSave)
	pipe.LTrim(ctx, MessagesKey+":messages", -MaxHistoryRounds*2, -1)
	// 设置消息列表过期时间
	pipe.Expire(ctx, MessagesKey+":messages", global.MESSAGES_EFFECTIVE_TIME)
	// 设置系统键过期时间
	pipe.Expire(ctx, MessagesKey+":system", global.MESSAGES_EFFECTIVE_TIME)
	pipe.Exec(ctx)
	return
}
