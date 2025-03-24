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
	// 参数有效性验证，确保PresencePenalty在合法范围内
	if req.PresencePenalty < -2.0 || req.PresencePenalty > 2.0 {
		zlog.CtxErrorf(ctx, "PresencePenalty out of range")
		return
	}
	zlog.CtxInfof(ctx, "AIChat request: %v", req)

	// 设置SSE协议要求的响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 设置redis key
	MessagesKey := fmt.Sprintf(global.REDIS_MESSAGES_KEY, userid, req.Type)

	// 获取AI聊天流并确保最终关闭响应体
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

	// 流式数据处理
	var fullResponse strings.Builder
	flusher, _ := c.Writer.(http.Flusher)
	//创建一个带缓冲的读取器，
	//将HTTP响应体(httpResp.Body)包装为缓冲I/O接口，通过减少系统调用次数提高数据读取效率。
	reader := bufio.NewReader(httpResp.Body)
	for {
		//从reader读取字节，直到遇到'\n'换行符，返回包含分隔符的字节切片
		line, err := reader.ReadBytes('\n')
		if err != nil {
			//遇到文件结束（io.EOF）时返回已读取的数据
			if err == io.EOF {
				break
			}
			response.SendSSEError(c, "流式传输中断")
			return
		}

		// 解析 SSE 数据
		// 检查一个字节切片是否以特定的前缀开头，识别 SSE 协议中的有效数据行（以 "data: " 开头的行）
		if bytes.HasPrefix(line, []byte("data: ")) {
			var chuck types.AIChatStreamResp
			//移除 SSE 数据头的 data: 前缀(6字节长度),将剩余部分解析为预定义的结构体 types.AIChatStreamResp
			if json.Unmarshal(line[6:], &chuck) != nil {
				zlog.CtxWarnf(ctx, "Invalid SSE data chunk: %s", line)
				continue // 跳过无效数据块
			}
			//从解析后的响应中提取 AI 生成的增量内容(delta content)
			content := chuck.Choices[0].Delta.Content
			if content != "" {
				// 按 SSE 协议格式推送，SSE规范要求每个消息以"data: "开头，后跟数据，然后是两个换行符
				escapedContent := strings.ReplaceAll(content, "\n", "\\n") // 处理content本身包含的换行符
				_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", escapedContent)
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
	//	添加流结束标识 明确告知客户端 流已传输完毕，避免因无限等待超时。
	_, err = fmt.Fprintf(c.Writer, "event: end\ndata: stream completed\n\n")
	if err != nil {
		response.SendSSEError(c, "流式传输中断")
		return
	}
	flusher.Flush()

	// 持久化存储AI对话记录
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
