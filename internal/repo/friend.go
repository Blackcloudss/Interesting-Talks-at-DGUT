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
	FOLLOWER_ID = "follower_id"
	FOLLOWED_ID = "followed_id"
	IS_FRiEND   = "is_friend"
	TRUE        = 1
	FALSE       = 0
)

// @Title        friend.go
// @Description
// @Create       XdpCs 2025-03-19 下午3:21
// @Update       XdpCs 2025-03-19 下午3:21
type FriendRepo struct {
	DB *gorm.DB
}

func NewFriendRepo(db *gorm.DB) *FriendRepo {
	return &FriendRepo{DB: db}
}

// JudgeFriend
//
//	@Description: 判断是否是好友
//	@receiver r
//	@param UserId
//	@param FollowedId
//	@return bool
//	@return error
func (r *FriendRepo) JudgeFriend(UserId, FollowedId int64) (bool, error) {
	err := r.DB.Model(&model.Follow{}).
		Where(fmt.Sprintf("%s = ? AND %s = ? AND %s = ?", FOLLOWER_ID, FOLLOWED_ID, IS_FRiEND), UserId, FollowedId, TRUE).
		First(&model.Follow{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Errorf("用户不存在：%v", err)
			return false, nil
		}
		zlog.Errorf("查询用户好友失败：%v", err)
		return false, err
	}
	return true, nil
}

// GetFriendsId
//
//	@Description: 获取好友ID
//	@receiver r
//	@param UserId
//	@return FriendsId
//	@return err
func (r *FriendRepo) GetFriendsId(UserId int64) (FriendsId []int64, err error) {
	err = r.DB.Model(&model.Follow{}).
		Where(fmt.Sprintf("%s = ? AND %s = ?", FOLLOWER_ID, IS_FRiEND), UserId, TRUE).
		Pluck(FOLLOWED_ID, &FriendsId).
		Error
	if err != nil {
		zlog.Errorf("查询用户好友失败：%v", err)
		return FriendsId, err
	}
	return
}

// GetFriendList
//
//	@Description: 获取好友信息
//	@receiver r
//	@param FriendsID
//	@return Friends
//	@return err
func (r *FriendRepo) GetFriendList(FriendsID []int64) (Friends []types.FriendInfo, err error) {
	var Friend types.FriendInfo
	for _, FriendId := range FriendsID {
		err = r.DB.Model(&model.UserDisplay{}).
			Select(ID, AVATAR, NICKNAME, TAG).
			Where(fmt.Sprintf("%s = ?", ID), FriendId).
			First(&Friend).
			Error
		if err != nil {
			zlog.Errorf("查询用户好友信息失败：%v", err)
			return Friends, err
		}
		Friends = append(Friends, Friend)
	}
	return
}
