package logic

import (
	"context"
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"gorm.io/gorm"
	"time"
)

// 定义内部逻辑错误
var (
	codeCommentNotFound     = response.MsgCode{Code: 40031, Msg: "评论不存在"}
	codeCommentCreateFailed = response.MsgCode{Code: 40032, Msg: "创建评论失败"}
	codeCommentDeleteFailed = response.MsgCode{Code: 40033, Msg: "删除评论失败"}
	codeCommentLikeFailed   = response.MsgCode{Code: 40034, Msg: "点赞评论失败"}
	codeCommentUnlikeFailed = response.MsgCode{Code: 40035, Msg: "取消点赞评论失败"}
	codeCommentAlreadyLiked = response.MsgCode{Code: 40036, Msg: "已经点赞过该评论"}
	codeCommentNotLiked     = response.MsgCode{Code: 40037, Msg: "未点赞该评论"}
)

type CommentLogic struct{}

// NewCommentLogic 创建评论逻辑层实例
func NewCommentLogic() *CommentLogic {
	return &CommentLogic{}
}

// CreateComment 创建评论
func (l *CommentLogic) CreateComment(ctx context.Context, req types.CreateCommentReq, UserID int64) (*types.CreateCommentResp, error) {
	defer utils.RecordTime(time.Now())()
	// 构建评论对象
	if req.RootParentID == 0 { // 一级评论
		comment := model.FirstComment{
			UserID:  UserID,
			BlogID:  req.BlogID,
			Content: req.Content,
		}
		// 创建一级评论
		commentRepo := repo.NewCommentRepo(global.DB)
		err := commentRepo.CreateFirstComment(&comment)
		if err != nil {
			zlog.CtxErrorf(ctx, "CreateComment failed: %v", err)
			return nil, response.ErrResp(err, codeCommentCreateFailed)
		}
		// 返回创建成功的一级评论信息
		resp := &types.CreateCommentResp{
			CommentID: comment.ID,
			CreatedAt: comment.CreatedAt,
		}
		return resp, nil
	} else { // 回复（二级评论）
		comment := model.SecondComment{
			UserID:       UserID,
			BlogID:       req.BlogID,
			Content:      req.Content,
			ParentID:     req.ParentID,
			RootParentID: req.RootParentID,
		}
		// 创建二级评论
		commentRepo := repo.NewCommentRepo(global.DB)
		err := commentRepo.CreateSecondComment(&comment)
		if err != nil {
			zlog.CtxErrorf(ctx, "CreateComment failed: %v", err)
			return nil, response.ErrResp(err, codeCommentCreateFailed)
		}
		// 返回创建成功的二级评论信息
		resp := &types.CreateCommentResp{
			CommentID: comment.ID,
			CreatedAt: comment.CreatedAt,
		}
		return resp, nil
	}
}

// DeleteComment 删除评论
func (l *CommentLogic) DeleteComment(ctx context.Context, req types.DeleteCommentReq, UserID int64) (*types.DeleteCommentResp, error) {
	defer utils.RecordTime(time.Now())()

	// 检查评论是否存在
	commentRepo := repo.NewCommentRepo(global.DB)
	firstComment, secondComment, err := commentRepo.GetCommentByID(req.CommentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "DeleteComment failed: comment not found (commentID: %d)", req.CommentID)
			return nil, response.ErrResp(err, codeCommentNotFound)
		}
		zlog.CtxErrorf(ctx, "GetCommentByID failed: %v", err)
		return nil, response.ErrResp(err, codeCommentDeleteFailed)
	}

	// 判断是一级评论还是二级评论
	isFirstComment := firstComment != nil

	var blogID int64
	if isFirstComment {
		blogID = firstComment.BlogID
	} else {
		blogID = secondComment.BlogID
	}

	// 删除评论并更新帖子的评论数
	err = commentRepo.DeleteComment(req.CommentID, blogID, isFirstComment)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteComment failed: %v", err)
		return nil, response.ErrResp(err, codeCommentDeleteFailed)
	}

	zlog.CtxInfof(ctx, "Comment deleted successfully (commentID: %d, blogID: %d)", req.CommentID, blogID)
	return &types.DeleteCommentResp{}, nil
}

// LikeComment 点赞评论
func (l *CommentLogic) LikeComment(ctx context.Context, req types.LikeCommentReq, UserID int64) (*types.LikeCommentResp, error) {
	defer utils.RecordTime(time.Now())()
	err := repo.NewCommentRepo(global.DB).LikeComment(UserID, req.CommentID)
	if err != nil {
		if errors.Is(err, errors.New("already liked this comment")) {
			return nil, response.ErrResp(err, codeCommentAlreadyLiked)
		}
		zlog.CtxErrorf(ctx, "LikeComment failed: %v", err)
		return nil, response.ErrResp(err, codeCommentLikeFailed)
	}
	return &types.LikeCommentResp{}, nil
}

// UnlikeComment 取消点赞评论
func (l *CommentLogic) UnlikeComment(ctx context.Context, req types.UnlikeCommentReq, UserID int64) (*types.UnlikeCommentResp, error) {
	defer utils.RecordTime(time.Now())()

	err := repo.NewCommentRepo(global.DB).UnlikeComment(UserID, req.CommentID)
	if err != nil {
		if errors.Is(err, errors.New("not liked this comment")) {
			return nil, response.ErrResp(err, codeCommentNotLiked)
		}
		zlog.CtxErrorf(ctx, "UnlikeComment failed: %v", err)
		return nil, response.ErrResp(err, codeCommentUnlikeFailed)
	}
	return &types.UnlikeCommentResp{}, nil
}

/*
func (l *CommentLogic) GetCommentList(ctx context.Context, req types.GetCommentListReq) (*types.GetCommentListResp, error) {
	defer utils.RecordTime(time.Now())()
	commentRepo := repo.NewCommentRepo(global.DB)
	comments, err := commentRepo.GetCommentList(req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCommentList failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	// 递归填充子回复
	for i := range comments {
		comments[i].Replies, err = commentRepo.GetRepliesList(comments[i].ID)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to get replies for comment %d: %v", comments[i].ID, err)
			return nil, response.ErrResp(err, response.INTERNAL_ERROR)
		}
	}
	zlog.CtxInfof(ctx, "Comment list with replies retrieved successfully (blogID: %d, count: %d)", req.BlogID, len(comments))
	resp := &types.GetCommentListResp{
		Comments: comments,
	}
	return resp, nil
}
*/
/*
func (l *CommentLogic) GetRepliesList(ctx context.Context, req types.GetRepliesListReq) (*types.GetRepliesListResp, error) {
	defer utils.RecordTime(time.Now())()
	commentRepo := repo.NewCommentRepo(global.DB)
	replies, err := commentRepo.GetRepliesList(req.CommentID)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetRepliesList failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	// 递归填充子回复
	for i := range replies {
		replies[i].Replies, err = commentRepo.GetRepliesList(replies[i].ID)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to get replies for reply %d: %v", replies[i].ID, err)
			return nil, response.ErrResp(err, response.INTERNAL_ERROR)
		}
	}
	zlog.CtxInfof(ctx, "Replies list retrieved successfully (commentID: %d, count: %d)", req.CommentID, len(replies))
	resp := &types.GetRepliesListResp{
		Replies: replies,
	}
	return resp, nil
}

// GetCommentListByLikes 获取评论列表（按点赞数排序）
func (l *CommentLogic) GetCommentListByLikes(ctx context.Context, req types.GetCommentListReq) (*types.GetCommentListResp, error) {
	defer utils.RecordTime(time.Now())()
	commentRepo := repo.NewCommentRepo(global.DB)
	comments, err := commentRepo.GetCommentListByLikes(req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCommentListByLikes failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	zlog.CtxInfof(ctx, "Comment list by likes retrieved successfully (blogID: %d, count: %d)", req.BlogID, len(comments))
	resp := &types.GetCommentListResp{
		Comments: comments,
	}
	return resp, nil
}*/
