package response

// @Title        error.go
// @Description
// @Create       XdpCs 2025-02-16 上午1:29
// @Update       XdpCs 2025-02-16 上午1:29
type RespError struct {
	JsonMsgResult
}

// NewRespError 包装响应错误类型，简化返回信息流程。
func ErrResp(err error, result MsgCode) error {
	respError := &RespError{}
	respError.Code = result.Code
	respError.Message = result.Msg
	respError.Data = err
	return respError
}

func (r RespError) Error() string {
	return r.Message
}
