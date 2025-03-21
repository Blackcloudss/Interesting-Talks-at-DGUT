package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/ai/role"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// @Title        ai.go
// @Description
// @Create       XdpCs 2025-03-09 下午2:27
// @Update       XdpCs 2025-03-09 下午2:27

func AIChatStream(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	userid := jwt.GetUserId(c)
	req, err := types.BindReq[types.AIChatStreamReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "AIChat request error: %v", err)
		return
	}
	// 检查参数是否在范围内
	if req.PresencePenalty < -2.0 || req.PresencePenalty > 2.0 {
		zlog.CtxErrorf(ctx, "PresencePenalty out of range")
		return
	}
	zlog.CtxInfof(ctx, "AIChat request: %v", req)

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	MessagesKey := fmt.Sprintf(global.REDIS_MESSAGES_KEY, userid, req.Type)

	httpResp, messages, err := logic.NewAILogic().GetChatStream(ctx, MessagesKey, req)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zlog.CtxErrorf(ctx, "关闭Body失败: %v", err)
			return
		}
	}(httpResp.Body)
	if httpResp.StatusCode != http.StatusOK {
		zlog.CtxErrorf(ctx, "HTTP请求失败，状态码：%d", httpResp.StatusCode)
		return
	}

	// 流式处理
	var fullResponse strings.Builder
	flusher, _ := c.Writer.(http.Flusher)
	reader := bufio.NewReader(httpResp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			response.SendSSEError(c, "流式传输中断")
			return
		}

		// 解析 SSE 数据
		// 识别 SSE 协议中的有效数据行（以 "data: " 开头的行）
		if bytes.HasPrefix(line, []byte("data: ")) {
			var chuck types.AIChatStreamResp
			//移除 SSE 数据头的 data: 前缀(6字节长度)
			//将剩余部分解析为预定义的结构体 types.AIChatStreamResp
			if json.Unmarshal(line[6:], &chuck) == nil {
				//从解析后的响应中提取 AI 生成的增量内容(delta content)
				content := chuck.Choices[0].Delta.Content
				if content != "" {
					// 按 SSE 协议格式推送
					_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", content)
					if err != nil {
						response.SendSSEError(c, "流式传输中断")
						return
					}
					// 立即刷新缓冲区
					flusher.Flush()
					// 逐步收集所有流式片段，最终拼成完整的 AI 响应内容
					fullResponse.WriteString(content)
				}
			}
		}
	}

	// 保存对话内容
	messages = append(messages, types.Message{
		Role:    role.ASSISTANT,
		Content: fullResponse.String(),
	})
	err = logic.NewAILogic().SaveHistoryMessages(ctx, MessagesKey, messages)
	if err != nil {
		zlog.CtxErrorf(ctx, "SaveHistoryMessages error: %v", err)
		return
	}
}

func DeleteHistoryMessages(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.DeleteHistoryMessagesReq](c)
	userid := jwt.GetUserId(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Delete HistoryMessages err:%v", err)
		return
	}
	zlog.CtxInfof(ctx, "HistoryMessages request: %v", req)
	MessagesKey := fmt.Sprintf(global.REDIS_MESSAGES_KEY, userid, req.Type)
	resp, err := logic.NewAILogic().DeleteHistoryMessages(ctx, MessagesKey)
	response.Response(c, resp, err)
	return
}
