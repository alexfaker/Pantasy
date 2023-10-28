package dto

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"codeup.aliyun.com/xhey/server/xappserver/public"
)

var (
	// IsEncrypted if encrypted
	IsEncrypted bool = false
	xlog             = public.NewLog()
)

type ResponseData struct {
	Status  int         `json:"status" default:"0"`           // 状态 0: 成功, 其它失败
	Payload interface{} `json:"payload" swaggerignore:"true"` // 返回的数据
}

type PageResult struct {
	TotalPage     uint64 `json:"totalPage" default:"10"` // 总页数
	TotalCount    uint64 `json:"total" default:"98"`     // 总条数
	PageNo        uint64 `json:"pageNo" default:"1"`     // 页码
	PageSize      uint64 `json:"pageSize" default:"10"`  // 每页数量
	PageStartTime int64  `json:"pageStartTime"`          // 第一页请求时间
}

type Response struct {
	Code     int         `json:"code" default:"0"`          // 状态 0: 成功, 其它失败
	Msg      string      `json:"msg" default:"succeed"`     // 错误信息
	ToastMsg string      `json:"toast_msg" default:""`      // 弹窗信息
	Data     interface{} `json:"data" swaggerignore:"true"` // 返回的数据
}

type H5Group struct {
	GroupID string `json:"groupID" ` //群唯一ID，长id
}

// ResponseStatus 响应状态
type ResponseStatus struct {
	Status int    `json:"status" default:"0"` // 状态 0: 成功, 其它失败
	Msg    string `json:"msg"`
}
type UserGroupSign struct {
	UserID  string `json:"userID" validate:"required"`            //用户的唯一id
	Sign    string `json:"sign" validate:"required,omitempty"`    //签名  userID+key求md5, 十六进制
	GroupID string `json:"groupID" validate:"required,omitempty"` //群唯一ID，长id
}

// UserSign 非团队请求
type UserSign struct {
	UserID string `json:"userID" validate:"required"`         //用户的唯一id
	Sign   string `json:"sign" validate:"required,omitempty"` //签名  userID+key求md5, 十六进制
}
type OptionalUserGroupSign struct {
	UserID  string `json:"userID"`  //用户的唯一id
	Sign    string `json:"sign"`    //签名  userID+key求md5, 十六进制
	GroupID string `json:"groupID"` //群唯一ID
}

func getReqParams(c *gin.Context, o interface{}) error {
	r := c.Request
	if IsEncrypted == true {
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		if len(strings.TrimSpace(string(body))) == 0 {
			return errors.New("the params is null")
		}

		var req struct {
			Params string `json:"params"`
		}

		if err := json.Unmarshal(body, &req); err != nil {
			return err
		}
		if len(req.Params) == 0 {
			return errors.New("the encrypted params is null")
		}
		var h public.AES_Handler
		body, err := h.Decrypt(req.Params)
		if err != nil {
			return err
		}
		json.Unmarshal(body, o)

		id, _ := c.Get("Conn-ID")
		xlog.Infof("[%v] %s parameter: %s, %s", id, r.RequestURI,
			(string)(body), public.HttpInfo(r))
		r.Header.Add("request-parameter", (string)(body))
	} else {
		if err := c.ShouldBind(o); err != nil {
			return err
		}
	}

	return nil

}

var ginValidate *validator.Validate

func BindingValidParams(c *gin.Context, o interface{}) error {
	if err := getReqParams(c, o); err != nil {
		return err
	}

	if ginValidate == nil {
		ginValidate = validator.New()
	}
	err := ginValidate.Struct(o)
	if err != nil {
		return err
	}

	return nil
}

// BindingParams bind json
func BindingParams(c *gin.Context, o interface{}) error {
	if err := getReqParams(c, o); err != nil {
		return err
	}
	return nil
}

func BindingParamsWithoutEncrypted(c *gin.Context, o interface{}) error {
	if err := c.ShouldBind(o); err != nil {
		return err
	}

	if ginValidate == nil {
		ginValidate = validator.New()
	}
	err := ginValidate.Struct(o)
	if err != nil {
		return err
	}
	return nil
}

// CheckGinValidate check
func CheckGinValidate(o interface{}) error {
	if ginValidate == nil {
		ginValidate = validator.New()
	}
	err := ginValidate.Struct(o)
	if err != nil {
		return err
	}

	return nil
}
