package response

// @Title        error.go
// @Description
// @Create       XdpCs 2025-02-16 上午1:29
// @Update       XdpCs 2025-02-16 上午1:29
type RespError struct {
	Code    int         `json:"code"`    // 确保标签与JsonMsgResult一致
	Message string      `json:"message"` // 必须与响应结构字段名匹配
	Data    interface{} `json:"data"`
}

// NewRespError 包装响应错误类型，简化返回信息流程。
func ErrResp(err error, result MsgCode) error {
	respError := &RespError{}
	respError.Code = result.Code
	respError.Message = result.Msg
	respError.Data = err
	return respError
}

func (e *RespError) Error() string {
	return e.Message
}
