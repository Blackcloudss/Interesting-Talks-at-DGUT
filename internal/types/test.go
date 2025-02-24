package types

// @Title        test.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:22
// @Update       XdpCs 2025-02-24 下午11:22
// 接收参数与返回参数层

type TestO1Req struct {
	UserID string `json:"user_id"`
}
type Test01Resp struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
