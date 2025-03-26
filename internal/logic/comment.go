package logic

import (
	"context"
	"errors"
	"time"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"gorm.io/gorm"
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
	codeGetCommentFail      = response.MsgCode{Code: 40041, Msg: "获取评论失败"}
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
	if req.ParentID == 0 { // 一级评论
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
			UserID:   UserID,
			BlogID:   req.BlogID,
			Content:  req.Content,
			ParentID: req.ParentID,
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
func (l *CommentLogic) DeleteComment(ctx context.Context, req types.DeleteCommentReq) (*types.DeleteCommentResp, error) {
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

// GetCommentList 获取一级评论列表，并显示部分二级评论
func (l *CommentLogic) GetCommentList(ctx context.Context, req types.GetCommentListReq) (resp *types.GetCommentListResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 调用 repo 层获取一级评论列表
	comments, err := repo.NewCommentRepo(global.DB).GetFirstCommentList(req.BlogID, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取一级评论列表失败: %v", err)
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	// 获取总评论数
	var totalCount int64
	if err := global.DB.Model(&model.FirstComment{}).
		Where("blog_id = ? AND deleted_at IS NULL", req.BlogID).
		Count(&totalCount).Error; err != nil {
		zlog.CtxErrorf(ctx, "获取评论总数失败: %v", err)
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	// 构造响应数据
	resp = &types.GetCommentListResp{
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Comments:   comments,
	}

	return resp, nil
}

// GetSecondCommentList 获取更多二级评论
func (l *CommentLogic) GetSecondCommentList(ctx context.Context, req types.GetSecondCommentListReq) (resp *types.GetSecondCommentListResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 调用 repo 层获取二级评论列表
	comments, totalCount, err := repo.NewCommentRepo(global.DB).GetSecondCommentList(req)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取二级评论列表失败: %v", err)
		return nil, response.ErrResp(err, codeGetCommentFail)
	}

	// 构造响应数据
	resp = &types.GetSecondCommentListResp{
		TotalCount:     totalCount,
		Page:           req.Page,
		PageSize:       req.PageSize,
		SecondComments: comments,
	}

	return resp, nil
}
