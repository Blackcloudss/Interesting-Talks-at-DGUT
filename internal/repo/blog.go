package repo

import (
	"errors"
	"fmt"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

const (
	BLOG_ID         = "blog_id" // 修改为 blog_id
	IS_LIKED        = "is_liked"
	IS_COLLECTED    = "is_collected"
	LIKE_COUNT      = "like_count"
	COLLECT_COUNT   = "collect_count"
	COMMENT_COUNT   = "comment_count"
	BLOG_TAG        = "blog_tag"
	SUB_TAG         = "sub_tag"
	VIEW_PERMISSION = "view_permission"
	TITLE           = "title"
	UPDATED_AT      = "updated_at"
	DELETED_AT      = "deleted_at"
	USER_TAG        = "tag"
)

const (
	BLOG_SELECT_FIELDS = `blog.id as blog_id, 
                         blog.created_at, 
                         blog.updated_at, 
                         blog.title, 
                         blog.content, 
                         blog.like_count, 
                         blog.collect_count, 
                         blog.comment_count, 
                         blog.blog_tag, 
                         blog.sub_tag, 
                         blog.view_permission, 
                         blog.user_id, 
                         user_display.nickname, 
                         user_display.avatar, 
                         user_display.tag`
)

type BlogRepo struct {
	DB *gorm.DB
}

// NewBlogRepo 创建帖子仓库实例
func NewBlogRepo(db *gorm.DB) *BlogRepo {
	return &BlogRepo{DB: db}
}

// CreateBlog 创建帖子
func (r *BlogRepo) CreateBlog(blog *model.Blog) error {
	if err := r.DB.Create(blog).Error; err != nil {
		zlog.Errorf(fmt.Sprintf("创建帖子失败: %v", err))
		return err
	}
	return nil
}

// UpdateBlog 更新帖子
func (r *BlogRepo) UpdateBlog(blog *model.Blog) error {
	if err := r.DB.Model(&model.Blog{}).
		Where(fmt.Sprintf("blog.%s = ?", ID), blog.ID).
		Updates(blog).Error; err != nil {
		zlog.Errorf(fmt.Sprintf("更新帖子失败: %v", err))
		return err
	}
	return nil
}

// DeleteBlog 删除帖子（软删除）
func (r *BlogRepo) DeleteBlog(blogID int64) error {
	if err := r.DB.Model(&model.Blog{}).
		Where(fmt.Sprintf("blog.%s = ?", ID), blogID).
		Update(DELETED_AT, gorm.Expr("NOW()")).Error; err != nil {
		zlog.Errorf(fmt.Sprintf("删除帖子失败: %v", err))
		return err
	}
	return nil
}

// CheckBlogExists 检查帖子是否存在
func (r *BlogRepo) CheckBlogExists(blogID int64) (*model.Blog, error) {
	var blog model.Blog
	err := r.DB.Model(&model.Blog{}).
		Where(fmt.Sprintf("blog.%s = ? AND %s IS NULL", ID, DELETED_AT), blogID).
		First(&blog).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		zlog.Errorf(fmt.Sprintf("检查帖子存在状态失败: %v", err))
		return nil, err
	}
	return &blog, nil
}

// GetBlogByID 获取单个帖子详情
func (r *BlogRepo) GetBlogByID(blogID int64) (types.BlogResp, error) {
	var blogResp types.BlogResp

	err := r.DB.Model(&model.Blog{}).
		Select(BLOG_SELECT_FIELDS).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where(fmt.Sprintf("blog.%s = ? AND blog.%s IS NULL", ID, DELETED_AT), blogID).
		Scan(&blogResp).
		Error

	if err != nil {
		zlog.Errorf(fmt.Sprintf("获取帖子详情失败: %v", err))
		return types.BlogResp{}, err
	}

	return blogResp, nil
}

// GetBlogs 分页获取帖子列表
func (r *BlogRepo) GetBlogs(page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	query := r.DB.Model(&model.Blog{}).
		Select(BLOG_SELECT_FIELDS).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where(fmt.Sprintf("blog.%s IS NULL", DELETED_AT)).
		Order(fmt.Sprintf("blog.%s DESC", CREATED_AT))

	if err := query.Count(&total).Error; err != nil {
		zlog.Errorf(fmt.Sprintf("统计帖子总数失败: %v", err))
		return nil, 0, err
	}

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).
		Error; err != nil {
		zlog.Errorf(fmt.Sprintf("获取帖子列表失败: %v", err))
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetBlogsByTag 根据标签获取帖子列表
func (r *BlogRepo) GetBlogsByTag(subTag string, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	query := r.DB.Model(&model.Blog{}).
		Select(BLOG_SELECT_FIELDS).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where(fmt.Sprintf("%s = ? AND blog.%s IS NULL", SUB_TAG, DELETED_AT), subTag).
		Order(fmt.Sprintf("blog.%s DESC", CREATED_AT))

	if err := query.Count(&total).Error; err != nil {
		zlog.Errorf(fmt.Sprintf("统计标签帖子总数失败: %v", err))
		return nil, 0, err
	}

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).
		Error; err != nil {
		zlog.Errorf(fmt.Sprintf("获取标签帖子列表失败: %v", err))
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetBlogsByUserID 获取用户发布的帖子列表
func (r *BlogRepo) GetBlogsByUserID(userID int64, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	query := r.DB.Model(&model.Blog{}).
		Select(BLOG_SELECT_FIELDS).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where(fmt.Sprintf("%s = ? AND blog.%s IS NULL", USER_ID, DELETED_AT), userID).
		Order(fmt.Sprintf("blog.%s DESC", CREATED_AT))

	if err := query.Count(&total).Error; err != nil {
		zlog.Errorf(fmt.Sprintf("统计用户帖子总数失败: %v", err))
		return nil, 0, err
	}

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).
		Error; err != nil {
		zlog.Errorf(fmt.Sprintf("获取用户帖子列表失败: %v", err))
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetCollectedBlogs 获取用户收藏的帖子列表
func (r *BlogRepo) GetCollectedBlogs(userID int64, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	query := r.DB.Model(&model.Blog{}).
		Select(BLOG_SELECT_FIELDS).
		Joins("INNER JOIN collection ON collection.blog_id = blog.id").
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where(fmt.Sprintf("collection.%s = ? AND blog.%s IS NULL", USER_ID, DELETED_AT), userID).
		Order(fmt.Sprintf("blog.%s DESC", CREATED_AT))

	if err := query.Count(&total).Error; err != nil {
		zlog.Errorf("统计收藏帖子总数失败: %v", err)
		return nil, 0, err
	}

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).
		Error; err != nil {
		zlog.Errorf("获取收藏帖子列表失败: %v", err)
		return nil, 0, err
	}

	return blogs, total, nil
}

func (r *BlogRepo) CollectBlog(UserID, BlogID int64) (*types.CollectBlogResp, error) {
	var isCollected bool = false

	tx := r.DB.Begin()

	// 查询收藏状态
	err := tx.Model(&model.Collection{}).
		Select(IS_COLLECTED).
		Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, BLOG_ID), UserID, BlogID).
		First(&isCollected).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 初始化收藏记录
			err = tx.Model(&model.Collection{}).
				Create(&model.Collection{
					UserID:      UserID,
					BlogID:      BlogID,
					IsCollected: false,
				}).Error
			if err != nil {
				tx.Rollback()
				zlog.Errorf("初始化收藏表失败：%v", err)
				return nil, err
			}
		} else {
			tx.Rollback()
			zlog.Errorf("查询收藏状态失败：%v", err)
			return nil, err
		}
	}

	// 查询当前收藏数
	var collectCount int64
	err = tx.Model(&model.Blog{}).
		Select(COLLECT_COUNT).
		Where(fmt.Sprintf("%s = ?", ID), BlogID).
		First(&collectCount).
		Error
	if err != nil {
		tx.Rollback()
		zlog.Errorf("查询收藏数失败：%v", err)
		return nil, err
	}

	// 切换收藏状态
	if !isCollected {
		// 执行收藏
		err = tx.Model(&model.Collection{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, BLOG_ID), UserID, BlogID).
			Update(IS_COLLECTED, true).
			Error
		if err != nil {
			tx.Rollback()
			zlog.Errorf("收藏操作失败：%v", err)
			return nil, err
		}
		collectCount++
	} else {
		// 取消收藏
		err = tx.Model(&model.Collection{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, BLOG_ID), UserID, BlogID).
			Update(IS_COLLECTED, false).
			Error
		if err != nil {
			tx.Rollback()
			zlog.Errorf("取消收藏失败：%v", err)
			return nil, err
		}
		collectCount--
	}

	// 更新博客收藏数
	err = tx.Model(&model.Blog{}).
		Where(fmt.Sprintf("%s = ?", ID), BlogID).
		Update(COLLECT_COUNT, collectCount).
		Error
	if err != nil {
		tx.Rollback()
		zlog.Errorf("更新收藏数失败：%v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &types.CollectBlogResp{
		IsCollected:  !isCollected,
		CollectCount: collectCount,
	}, nil
}

func (r *BlogRepo) LikeBlog(UserID, BlogID int64) (*types.LikeBlogResp, error) {
	var isLiked bool = false

	tx := r.DB.Begin()

	// 查询点赞状态
	err := tx.Model(&model.Like{}).
		Select(IS_LIKED).
		Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, BLOG_ID), UserID, BlogID).
		First(&isLiked).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 初始化点赞记录
			err = tx.Model(&model.Like{}).
				Create(&model.Like{
					UserID:  UserID,
					BlogID:  BlogID,
					IsLiked: false,
				}).Error
			if err != nil {
				tx.Rollback()
				zlog.Errorf("初始化点赞表失败：%v", err)
				return nil, err
			}
		} else {
			tx.Rollback()
			zlog.Errorf("查询点赞状态失败：%v", err)
			return nil, err
		}
	}

	// 查询当前点赞数
	var likeCount int64
	err = tx.Model(&model.Blog{}).
		Select(LIKE_COUNT).
		Where(fmt.Sprintf("%s = ?", ID), BlogID).
		First(&likeCount).
		Error
	if err != nil {
		tx.Rollback()
		zlog.Errorf("查询点赞数失败：%v", err)
		return nil, err
	}

	// 切换点赞状态
	if !isLiked {
		// 执行点赞
		err = tx.Model(&model.Like{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, BLOG_ID), UserID, BlogID).
			Update(IS_LIKED, true).
			Error
		if err != nil {
			tx.Rollback()
			zlog.Errorf("点赞操作失败：%v", err)
			return nil, err
		}
		likeCount++
	} else {
		// 取消点赞
		err = tx.Model(&model.Like{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, BLOG_ID), UserID, BlogID).
			Update(IS_LIKED, false).
			Error
		if err != nil {
			tx.Rollback()
			zlog.Errorf("取消点赞失败：%v", err)
			return nil, err
		}
		likeCount--
	}

	// 更新博客点赞数
	err = tx.Model(&model.Blog{}).
		Where(fmt.Sprintf("%s = ?", ID), BlogID).
		Update(LIKE_COUNT, likeCount).
		Error
	if err != nil {
		tx.Rollback()
		zlog.Errorf("更新点赞数失败：%v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &types.LikeBlogResp{
		IsLiked:   !isLiked,
		LikeCount: likeCount,
	}, nil
}
