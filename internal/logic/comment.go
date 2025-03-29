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

func NewCommentLogic() *CommentLogic {
	return &CommentLogic{}
}

func (l *CommentLogic) CreateComment(ctx context.Context, req types.CreateCommentReq, userID int64) (*types.CreateCommentResp, error) {
	defer utils.RecordTime(time.Now())()

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

	zlog.CtxInfof(ctx, "评论创建成功 (commentID: %d, userID: %d)", comment.ID, userID)
	return &types.CreateCommentResp{
		CommentID: comment.ID,
		CreatedAt: comment.CreatedAt,
	}, nil
}

// DeleteComment 删除评论（优化版）
func (l *CommentLogic) DeleteComment(ctx context.Context, req types.DeleteCommentReq, userID int64) (*types.DeleteCommentResp, error) {
	defer utils.RecordTime(time.Now())()

	// 1. 获取评论信息
	comment, err := repo.NewCommentRepo(global.DB).GetCommentByID(req.CommentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "评论不存在 ID:%d", req.CommentID)
			return nil, response.ErrResp(err, codeCommentNotFound)
		}
		zlog.CtxErrorf(ctx, "获取评论失败 ID:%d 错误:%v", req.CommentID, err)
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	// 2. 执行删除
	if err := repo.NewCommentRepo(global.DB).DeleteComment(comment); err != nil {
		zlog.CtxErrorf(ctx, "删除评论失败 ID:%d 错误:%v", req.CommentID, err)
		return nil, response.ErrResp(err, codeCommentDeleteFailed)
	}

	zlog.CtxInfof(ctx, "评论删除成功 ID:%d", req.CommentID)
	return &types.DeleteCommentResp{}, nil
}

func (l *CommentLogic) LikeComment(ctx context.Context, commentID int64, userID int64) (*types.LikeCommentResp, error) {
	defer utils.RecordTime(time.Now())()

	// 用 redis 加锁防止重复操作
	lockKey := fmt.Sprintf("comment:like:lock:user:%d:comment:%d", userID, commentID)
	locked, err := global.Rdb.SetNX(ctx, lockKey, 1, 1*time.Second).Result()
	if err != nil {
		zlog.CtxErrorf(ctx, "Redis 上锁失败: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	if !locked {
		zlog.CtxInfof(ctx, "操作频繁，请稍后重试 (userID: %d, commentID: %d)", userID, commentID)
		return nil, response.ErrResp(errors.New("操作过于频繁"), response.USER_OPERATION_LOCKED)
	}
	defer global.Rdb.Del(ctx, lockKey)

	// 检查评论是否存在
	if _, err := repo.NewCommentRepo(global.DB).GetCommentByID(commentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "评论不存在 (commentID: %d)", commentID)
			return nil, response.ErrResp(err, codeCommentNotFound)
		}
		zlog.CtxErrorf(ctx, "获取评论失败: %v", err)
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	resp, err := repo.NewCommentRepo(global.DB).LikeComment(userID, commentID)
	if err != nil {
		zlog.CtxErrorf(ctx, "点赞操作失败: %v", err)
		return nil, response.ErrResp(err, codeCommentLikeFailed)
	}

	zlog.CtxInfof(ctx, "点赞状态更新成功 (commentID: %d, action: %s)", commentID, resp.IsLiked)
	return resp, nil
}

func (l *CommentLogic) GetCommentList(ctx context.Context, req types.GetCommentListReq) (*types.GetCommentListResp, error) {
	defer utils.RecordTime(time.Now())()

	comments, total, err := repo.NewCommentRepo(global.DB).GetComments(
		req.BlogID,
		req.Page,
		req.PageSize,
		req.SortBy,
	)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取评论列表失败: %v", err)
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	zlog.CtxInfof(ctx, "获取评论列表成功 (blogID: %d, count: %d)", req.BlogID, len(comments))
	return &types.GetCommentListResp{
		TotalCount: total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Comments: comments,
	}, nil
}
func (l *CommentLogic) GetRepliesList(ctx context.Context, req types.GetRepliesListReq) (*types.GetRepliesListResp, error) {
	defer utils.RecordTime(time.Now())()

	replies, total, err := repo.NewCommentRepo(global.DB).GetReplies(
		req.ParentID,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取回复列表失败: %v", err)
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	zlog.CtxInfof(ctx, "获取回复列表成功 (parentID: %d, count: %d)", req.ParentID, len(replies))
	return &types.GetRepliesListResp{
		TotalCount: total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Replies: replies,
	}, nil
}
