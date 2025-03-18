package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/image"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/userinfo"
	"time"
)

var (
	REDIS_GET_FAULT        = response.MsgCode{50002, "SessionKey过期"}
	DECRYPT_DATA_FAILED    = response.MsgCode{50003, "解密数据失败"}
	ANALYSIS_FAILED        = response.MsgCode{50004, "解析数据失败"}
	SAVE_FAILED            = response.MsgCode{50005, "保存用户信息失败"}
	GET_COMMON_PROFILE     = response.MsgCode{50006, "获取用户基本信息失败"}
	UPDATE_COMMON_PROFILE  = response.MsgCode{50008, "更改用户公基本信息失败"}
	GET_PRIVATE_PROFILE    = response.MsgCode{50009, "获取用户隐私信息失败"}
	UPDATE_PRIVATE_PROFILE = response.MsgCode{50010, "更改用户隐私信息失败"}
)

// @Title        user.go
// @Description
// @Create       XdpCs 2025-03-14 上午11:47
// @Update       XdpCs 2025-03-14 上午11:47
type UserLogic struct {
}

func NewUserLogic() *UserLogic {
	return &UserLogic{}
}

// 获取微信头像，昵称
func (l *UserLogic) GetUserInfo(ctx context.Context, UserId int64, req types.UserInfoReq) (resp types.UserInfoResp, err error) {
	defer utils.RecordTime(time.Now())()
	OpenId, err := repo.NewUserRepo(global.DB).GetOpenId(UserId)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetUserInfo err: %v", err)
		return resp, response.ErrResp(err, response.OPENID_NOT_EXIST)
	}
	//把用户的SessionKey从redis中拿出来
	SessionKey, err := global.Rdb.Get(ctx, fmt.Sprintf(global.REDIS_SESSIONKEY, OpenId)).Result()
	if err != nil {
		zlog.CtxErrorf(ctx, "redis get session_key err: %v", err)
		return resp, response.ErrResp(err, REDIS_GET_FAULT)
	}
	//验证签名
	if !userinfo.VerifySignature(SessionKey, req.RawData, req.Signature) {
		zlog.CtxErrorf(ctx, "签名验证失败,SessionKey过期")
		return resp, response.ErrResp(err, REDIS_GET_FAULT)
	}
	// 解密数据
	encryptedBytes, _ := base64.StdEncoding.DecodeString(req.EncryptedData)
	ivBytes, _ := base64.StdEncoding.DecodeString(req.Iv)
	decryptedData, err := userinfo.DecryptData(SessionKey, encryptedBytes, ivBytes)
	if err != nil {
		zlog.CtxErrorf(ctx, "解密数据失败: %v", err)
		return resp, response.ErrResp(err, DECRYPT_DATA_FAILED)
	}
	// 解析用户信息
	var UserInfo types.UserInfo
	if err = json.Unmarshal(decryptedData, &UserInfo); err != nil {
		zlog.CtxErrorf(ctx, "解析用户信息失败: %v", err)
		return resp, response.ErrResp(err, ANALYSIS_FAILED)
	}
	//存储用户信息到数据库中
	err = repo.NewUserRepo(global.DB).SaveUserInfo(UserId, UserInfo)
	if err != nil {
		zlog.CtxErrorf(ctx, "存储用户信息失败: %v", err)
		return resp, response.ErrResp(err, SAVE_FAILED)
	}
	return
}

// GetCommonProfile
//
//	@Description: 获取用户基本信息
//	@receiver l
//	@param ctx
//	@param UserId
//	@return resp
//	@return err
func (l *UserLogic) GetCommonProfile(ctx context.Context, UserId int64) (resp types.GetCommonProfileResp, err error) {
	defer utils.RecordTime(time.Now())()
	resp, err = repo.NewUserRepo(global.DB).GetCommonProfile(UserId)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取用户基本信息失败: %v", err)
		return resp, response.ErrResp(err, GET_COMMON_PROFILE)
	}
	return
}

// UpdateCommonProfile
//
//	@Description:  更新用户基本信息
//	@receiver l
//	@param ctx
//	@param UserId
//	@param req
//	@return resp
//	@return err
func (l *UserLogic) UpdateCommonProfile(ctx context.Context, UserId int64, req types.UpdateCommonProfileReq) (resp types.UpdateCommonProfileResp, err error) {
	defer utils.RecordTime(time.Now())()
	AvatarUrl, err := image.UploadImage(req.Avatar)
	if err != nil {
		zlog.CtxErrorf(ctx, "Upload image error: %v", err)
		return resp, response.ErrResp(err, response.INTERANL_IMAGE_UPLOAD_ERROR)
	}
	err = repo.NewUserRepo(global.DB).UpdateCommonProfile(UserId, AvatarUrl, req)
	if err != nil {
		zlog.CtxErrorf(ctx, "更改用户基本信息失败: %v", err)
		return resp, response.ErrResp(err, UPDATE_COMMON_PROFILE)
	}
	return
}

// GetPrivateProfile
//
//	@Description: 获取用户隐私信息
//	@receiver l
//	@param ctx
//	@param UserId
//	@return resp
//	@return err
func (l *UserLogic) GetPrivateProfile(ctx context.Context, UserId int64) (resp types.GetPrivateProfileResp, err error) {
	defer utils.RecordTime(time.Now())()
	resp, err = repo.NewUserRepo(global.DB).GetPrivateProfile(UserId)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取用户隐私信息失败: %v", err)
		return resp, response.ErrResp(err, GET_PRIVATE_PROFILE)
	}
	return
}

// UpdatePrivateProfile
//
//	@Description: 更新用户隐私信息
//	@receiver l
//	@param ctx
//	@param UserId
//	@param req
//	@return resp
//	@return err
func (l *UserLogic) UpdatePrivateProfile(ctx context.Context, UserId int64, req types.UpdatePrivateProfileReq) (resp types.UpdatePrivateProfileResp, err error) {
	defer utils.RecordTime(time.Now())()
	role, err := repo.NewUserRepo(global.DB).UpdatePrivateProfile(UserId, req)
	if err != nil {
		zlog.CtxErrorf(ctx, "更改用户隐私信息失败: %v", err)
		return resp, response.ErrResp(err, UPDATE_PRIVATE_PROFILE)
	}
	resp.Role = role
	return
}
