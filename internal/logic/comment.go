package logic

import (
	"context"
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

// 定义内部逻辑错误
var (
	codeCommentNotFound     = response.MsgCode{Code: 40031, Msg: "评论不存在"}
	codeCommentCreateFailed = response.MsgCode{Code: 40032, Msg: "创建评论失败"}
	codeCommentDeleteFailed = response.MsgCode{Code: 40033, Msg: "删除评论失败"}
	codeUnauthorized        = response.MsgCode{Code: 40034, Msg: "用户未授权"}
	codeUserNotLoggedIn     = response.MsgCode{Code: 20001, Msg: "用户未登录"}
)

type CommentLogic struct{}

// NewCommentLogic 创建评论逻辑层实例
func NewCommentLogic() *CommentLogic {
	return &CommentLogic{}
}

// CreateComment 创建评论
func (l *CommentLogic) CreateComment(ctx context.Context, req types.CreateCommentReq) (*types.CreateCommentResp, error) {
	/*// 检查用户是否登录
	if req.UserID == 0 {
		zlog.CtxErrorf(ctx, "CreateComment failed: user not logged in")
		return nil, response.ErrResp(nil, codeUserNotLoggedIn)
	}*/

	// 检查帖子是否存在
	blogRepo := repo.NewBlogRepo(global.DB)
	_, err := blogRepo.GetBlogByID(req.BlogID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "CreateComment failed: blog not found (blogID: %d)", req.BlogID)
			return nil, response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "CreateComment failed: %v", err)
		return nil, response.ErrResp(err, codeCommentCreateFailed)
	}

	// 创建评论并更新帖子的评论数
	commentRepo := repo.NewCommentRepo(global.DB)
	comment, err := commentRepo.CreateComment(req)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateComment failed: %v", err)
		return nil, response.ErrResp(err, codeCommentCreateFailed)
	}

	zlog.CtxInfof(ctx, "Comment created successfully (commentID: %d, blogID: %d)", comment.ID, comment.BlogID)
	return &types.CreateCommentResp{CommentID: comment.ID}, nil
}

// DeleteComment 删除评论
func (l *CommentLogic) DeleteComment(ctx context.Context, req *types.DeleteCommentReq) (*types.DeleteCommentResp, error) {
	// 检查评论是否存在
	commentRepo := repo.NewCommentRepo(global.DB)
	comment, err := commentRepo.GetCommentByID(req.CommentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "DeleteComment failed: comment not found (commentID: %d)", req.CommentID)
			return nil, response.ErrResp(err, codeCommentNotFound)
		}
		zlog.CtxErrorf(ctx, "GetCommentByID failed: %v", err)
		return nil, response.ErrResp(err, codeCommentDeleteFailed)
	}

	// 删除评论并更新帖子的评论数
	if err := commentRepo.DeleteComment(req.CommentID, comment.BlogID); err != nil {
		zlog.CtxErrorf(ctx, "DeleteComment failed: %v", err)
		return nil, response.ErrResp(err, codeCommentDeleteFailed)
	}

	zlog.CtxInfof(ctx, "Comment deleted successfully (commentID: %d, blogID: %d)", req.CommentID, comment.BlogID)
	return &types.DeleteCommentResp{}, nil
}

// GetCommentList 获取评论列表
func (l *CommentLogic) GetCommentList(ctx context.Context, req types.GetCommentListReq) (*types.GetCommentListResp, error) {
	commentRepo := repo.NewCommentRepo(global.DB)
	comments, err := commentRepo.GetCommentList(req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCommentList failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	zlog.CtxInfof(ctx, "Comment list retrieved successfully (blogID: %d, count: %d)", req.BlogID, len(comments))
	resp := &types.GetCommentListResp{
		Comments: comments,
	}
	return resp, nil
}
