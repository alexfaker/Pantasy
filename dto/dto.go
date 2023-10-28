package dto

type Response struct {
	Code     int         `json:"code" default:"0"`          // 状态 0: 成功, 其它失败
	Msg      string      `json:"msg" default:"succeed"`     // 错误信息
	ToastMsg string      `json:"toast_msg" default:""`      // 弹窗信息
	Data     interface{} `json:"data" swaggerignore:"true"` // 返回的数据
}
