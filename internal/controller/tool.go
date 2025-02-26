package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	// 定义用户未登录的错误
	ErrorUserNotLogin = errors.New("用户未登录")
)

// getCurrentUserID 获取当前登录用户ID
func getCurrentUserID(c *gin.Context) (userID uint64, err error) {
	_userID, ok := c.Get("user_id")
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = _userID.(uint64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

// 分页
func paginate(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.Query("page"))         // 从查询参数中获取页码，默认为1
	pageSize, _ := strconv.Atoi(c.Query("pageSize")) // 从查询参数中获取每页大小，默认为10

	if page <= 0 { // 如果页码小于等于0
		page = 1 // 设置默认页码为1
	}
	if pageSize <= 0 { // 如果每页大小小于等于0
		pageSize = 10 // 设置默认每页大小为10
	}
	return page, pageSize // 返回分页参数
}
