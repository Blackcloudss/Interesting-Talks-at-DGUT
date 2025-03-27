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
	BLOG_ID      = "blog_id"
	IS_LIKED     = "is_liked"
	IS_COLLECTED = "is_collected"
	BE_COLLECTED = "be_collected"
	BE_LIKED     = "be_liked"
)
const blogSelectFields = `blog.id AS blog_id, blog.created_at AS create_at, blog.updated_at AS update_at, 
                          blog.title, blog.content, blog.be_liked, blog.be_collected, blog.comment_count, 
                          blog.blog_tag, blog.sub_tag, blog.view_permission, blog.user_id, 
                          user_display.nickname, user_display.avatar, user_display.tag`

// BlogRepo 帖子仓库
type BlogRepo struct {
	DB *gorm.DB
}

// NewBlogRepo 创建帖子仓库实例
func NewBlogRepo(db *gorm.DB) *BlogRepo {
	return &BlogRepo{
		DB: db,
	}
}

// CreateBlog 创建帖子
func (r *BlogRepo) CreateBlog(blog *model.Blog) error {
	return r.DB.Create(blog).Error
}

// UpdateBlog 更新帖子
func (r *BlogRepo) UpdateBlog(blog *model.Blog) error {
	return r.DB.Save(blog).Error
}

// DeleteBlog 删除帖子
func (r *BlogRepo) DeleteBlog(blogID int64) error {
	return r.DB.Delete(&model.Blog{}, blogID).Error
}

// CheckBlogExists 检查帖子是否存在，如果存在则返回帖子对象
func (r *BlogRepo) CheckBlogExists(blogID int64) (*model.Blog, error) {
	var blog model.Blog
	// 查询帖子是否存在
	err := r.DB.Model(&model.Blog{}).Where("id = ?", blogID).First(&blog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 如果帖子不存在，返回 nil
		}
		return nil, err // 如果发生其他错误，返回错误
	}
	return &blog, nil // 如果帖子存在，返回帖子对象
}

// GetBlogByID 根据ID获取帖子，帖子详情页面
func (r *BlogRepo) GetBlogByID(blogID int64) (types.BlogResp, error) {
	var blogResp types.BlogResp

	// 联合查询帖子和用户信息
	if err := r.DB.Model(&model.Blog{}).
		Select(blogSelectFields).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where("blog.id = ?", blogID).
		Scan(&blogResp).Error; err != nil {
		return types.BlogResp{}, err
	}

	return blogResp, nil
}

// GetBlogs 分页获取帖子列表（首页显示）
func (r *BlogRepo) GetBlogs(page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	query := r.DB.Model(&model.Blog{}).
		Select(blogSelectFields).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where("blog.deleted_at IS NULL").
		Order("blog.created_at DESC")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Scan(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}
func (r *BlogRepo) GetBlogsByTag(subTag string, page int, pageSize int) ([]types.BlogResp, int64, error) {
	if r.DB == nil {
		return nil, 0, errors.New("database connection is nil")
	}

	var blogs []types.BlogResp
	var total int64
	// 明确指定 deleted_at 列所属的表
	whereCondition := "sub_tag = ? AND blog.deleted_at IS NULL"

	// 查询符合条件的帖子总数
	if err := r.DB.Model(&model.Blog{}).
		Where(whereCondition, subTag).
		Count(&total).Error; err != nil {
		zlog.Errorf("Failed to count blogs by sub tag: %v", err)
		return nil, 0, err
	}

	// 分页查询帖子和用户信息
	if err := r.DB.Model(&model.Blog{}).
		Select(blogSelectFields).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where(whereCondition, subTag).
		Order("blog.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).Error; err != nil {
		zlog.Errorf("Failed to get blogs by sub tag: %v", err)
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetBlogsByUserID 根据用户ID分页获取帖子列表
func (r *BlogRepo) GetBlogsByUserID(userID int64, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	query := r.DB.Model(&model.Blog{}).
		Select(blogSelectFields).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where("blog.user_id = ? AND blog.deleted_at IS NULL", userID).
		Order("blog.created_at DESC")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Scan(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetCollectedBlogs 获取用户收藏的帖子（分页）
func (r *BlogRepo) GetCollectedBlogs(userID int64, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	query := r.DB.Model(&model.Blog{}).
		Select(blogSelectFields).
		Joins("INNER JOIN collection ON collection.blog_id = blog.id").
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where("collection.user_id = ?", userID).
		Order("blog.created_at DESC")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Scan(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// CollectBlog 收藏/取消收藏帖子
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
		Select(BE_COLLECTED).
		Where(fmt.Sprintf("%s = ? ", BLOG_ID), BlogID).
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
		Where(fmt.Sprintf("%s = ? ", BLOG_ID), BlogID).
		Update(BE_COLLECTED, collectCount).
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
		IsCollected:  !isCollected, // 返回操作后的新状态
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
		Select(BE_LIKED).
		Where("id = ?", BlogID).
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
		Where(fmt.Sprintf("%s = ? ", BLOG_ID), BlogID).
		Update(BE_LIKED, likeCount).
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
		IsLiked:   !isLiked, // 返回操作后的新状态
		LikeCount: likeCount,
	}, nil
}
