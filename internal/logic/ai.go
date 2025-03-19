package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
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

func (l *AILogic) AIChat(ctx context.Context, req types.AIChatReq) (resp *types.AIChatResp, err error) {

	return
}
