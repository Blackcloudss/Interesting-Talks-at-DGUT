package logic

import (
	"context"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// 定义内部逻辑错误
var (
	codeCommentNotFound     = response.MsgCode{Code: 40031, Msg: "评论不存在"}
	codeCommentCreateFailed = response.MsgCode{Code: 40032, Msg: "创建评论失败"}
	codeCommentDeleteFailed = response.MsgCode{Code: 40033, Msg: "删除评论失败"}
	codeCommentLikeFailed   = response.MsgCode{Code: 40034, Msg: "点赞或取消点赞评论失败"}
	codeGetCommentFail      = response.MsgCode{Code: 40041, Msg: "获取评论失败"}
)

type CommentLogic struct{}

// NewCommentLogic 创建评论逻辑层实例
func NewCommentLogic() *CommentLogic {
	return &CommentLogic{}
}

func (l *CommentLogic) CreateComment(ctx context.Context, req types.CreateCommentReq, userID int64) (*types.CreateCommentResp, error) {
	comment := &model.Comment{
		UserID:  userID,
		BlogID:  req.BlogID,
		Content: req.Content,
	}

	if req.ParentID != 0 {
		comment.ParentID = req.ParentID
	}

	if err := repo.NewCommentRepo(global.DB).CreateComment(comment); err != nil {
		zlog.CtxErrorf(ctx, "创建评论失败: %v", err)
		return nil, response.ErrResp(err, codeCommentCreateFailed)
	}

	return &types.CreateCommentResp{
		CommentID: comment.ID,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (l *CommentLogic) DeleteComment(ctx context.Context, req types.DeleteCommentReq, userID int64) (*types.DeleteCommentResp, error) {
	defer utils.RecordTime(time.Now())
	err := repo.NewCommentRepo(global.DB).DeleteComment(req.CommentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxErrorf(ctx, "评论不存在或用户无权限 (commentID: %d)", req.CommentID)
			return nil, response.ErrResp(err, codeCommentNotFound)
		}
		zlog.CtxErrorf(ctx, "删除评论失败: %v", err)
		return nil, response.ErrResp(err, codeCommentDeleteFailed)
	}

	zlog.CtxInfof(ctx, "评论删除成功 (commentID: %d)", req.CommentID)
	return &types.DeleteCommentResp{}, nil
}

// LikeComment 点赞或取消点赞评论
func (l *CommentLogic) LikeComment(ctx context.Context, commentID int64, userID int64) (resp *types.LikeCommentResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 用 redis 加锁
	lockKey := fmt.Sprintf("comment:like:lock:user:%d:comment:%d", userID, commentID)
	locked, err := global.Rdb.SetNX(ctx, lockKey, 1, 1*time.Second).Result()
	if err != nil {
		zlog.CtxErrorf(ctx, "Redis 上锁失败: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	if !locked {
		// 未获取到锁，说明该操作正在被其他请求处理
		zlog.CtxInfof(ctx, "点赞/取消赞操作正被 user: %d, comment: %d 使用，请稍等 1 s", userID, commentID)
		return nil, response.ErrResp(err, response.USER_OPERATION_LOCKED) // 用户操作被锁定
	}
	defer global.Rdb.Del(ctx, lockKey)

	resp, err = repo.NewCommentRepo(global.DB).LikeComment(userID, commentID)
	if err != nil {
		zlog.CtxErrorf(ctx, "ToggleLikeComment failed: %v", err)
		return nil, response.ErrResp(err, codeCommentLikeFailed)
	}
	return resp, nil
}

func (l *CommentLogic) GetCommentList(ctx context.Context, req types.GetCommentListReq) (*types.GetCommentListResp, error) {
	comments, total, err := repo.NewCommentRepo(global.DB).GetCommentList(
		req.BlogID,
		req.Page,
		req.PageSize,
		req.SortBy,
	)
	if err != nil {
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	return &types.GetCommentListResp{
		TotalCount: total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Comments:   comments,
	}, nil
}

func (l *CommentLogic) GetRepliesList(ctx context.Context, req types.GetRepliesListReq) (*types.GetRepliesListResp, error) {
	replies, total, err := repo.NewCommentRepo(global.DB).GetReplies(
		req.ParentID,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	return &types.GetRepliesListResp{
		TotalCount: total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Replies:    replies,
	}, nil
}
