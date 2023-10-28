package middleware

import (
	"bytes"
	"fmt"
	"github.com/alexfaker/Pantasy/dto"
	"github.com/alexfaker/Pantasy/middleware/log"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var IsEncrypted = false

var statusText = map[int]string{
	StatusOK:                  "success",
	StatusBadRequest:          "param error",
	StatusUnauthorized:        "verify error",
	StatusForbidden:           "forbidden",
	StatusNotFound:            "not found",
	StatusInternalServerError: "service error",
}

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// 返回不加密的结果
func PlainResponse(c *gin.Context, status int, data interface{}) {
	response(c, status, data, false, false)
}

func Response(c *gin.Context, status int, data interface{}) {
	response(c, status, data, false, IsEncrypted)
}

func ResponseEscapeHTML(c *gin.Context, status int, data interface{}) {
	responseEscapeHTML(c, status, data, false, IsEncrypted)
}

func ResponseBase64(c *gin.Context, status int, data interface{}) {
	responseBase64(c, status, data, false, IsEncrypted)
}

func ResponsePayload(c *gin.Context, status int, data interface{}) {
	response(c, status, data, true, IsEncrypted)
}

func response(c *gin.Context, status int, data interface{}, directly, isEncrypted bool) {
	var resp interface{}
	resp = dto.Response{
		Code:     status,
		Msg:      statusText[status],
		ToastMsg: "",
		Data:     data,
	}
	if directly {
		resp = data
	}
	c.Set("response", data)

	if isEncrypted == true {
		//var h public.AES_Handler
		//body, _ := jsoniter.Marshal(resp)
		//buf, _ := h.Encrypt(body)
		//
		//type responseEncrypted struct {
		//	Result string `json:"result"`
		//}
		//
		//resp = responseEncrypted{
		//	Result: buf,
		//}
	}

	startTime := time.Now().UnixNano() / 1000 // Golang 1.15不支持获取微妙时间戳
	buf, err := jsoniter.Marshal(resp)
	endTime := time.Now().UnixNano() / 1000

	if err != nil {
		log.Errorf(err.Error())
		return
	}

	log.Infof("uri:%v,json序列化耗时:%v微妙,序列化后的大小:%v字节", c.Request.RequestURI, endTime-startTime, len(buf))

	c.Writer.Header().Add("Content-Length", fmt.Sprintf("%v", len(buf)))

	nWriteIndex := 0
	for nWriteIndex < len(buf) {
		n, err := c.Writer.Write(buf[nWriteIndex:])
		if err != nil {
			log.Errorf("%v len:%v index:%v n:%v", err.Error(), len(buf), nWriteIndex, n)
			return
		}
		nWriteIndex += n
	}
	c.Writer.Flush()
}

func responseBase64(c *gin.Context, status int, data interface{}, directly, isEncrypted bool) {
	var resp interface{}
	resp = dto.Response{
		Code:     status,
		Msg:      statusText[status],
		ToastMsg: "",
		Data:     data,
	}
	if directly {
		resp = data
	}
	c.Set("response", data)

	//if isEncrypted == true {
	//	var h public.AES_Handler
	//	body, _ := jsoniter.Marshal(resp)
	//	body = []byte(base64.StdEncoding.EncodeToString(body))
	//	buf, _ := h.Encrypt(body)
	//
	//	type responseEncrypted struct {
	//		Result string `json:"result"`
	//	}
	//
	//	resp = responseEncrypted{
	//		Result: buf,
	//	}
	//}
	buf, err := jsoniter.Marshal(resp)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	c.Writer.Header().Add("Content-Length", fmt.Sprintf("%v", len(buf)))
	nWriteIndex := 0
	for nWriteIndex < len(buf) {
		n, err := c.Writer.Write(buf[nWriteIndex:])
		if err != nil {
			log.Errorf("%v len:%v index:%v n:%v", err.Error(), len(buf), nWriteIndex, n)
			return
		}
		nWriteIndex += n
	}
	c.Writer.Flush()
}

func responseEscapeHTML(c *gin.Context, status int, data interface{}, directly, isEncrypted bool) {
	var resp interface{}
	resp = dto.Response{
		Code:     status,
		Msg:      statusText[status],
		ToastMsg: "",
		Data:     data,
	}
	if directly {
		resp = data
	}
	c.Set("response", data)

	//if isEncrypted == true {
	//	var h public.AES_Handler
	//	buffer := &bytes.Buffer{}
	//	encoder := jsoniter.NewEncoder(buffer)
	//	encoder.SetEscapeHTML(false)
	//	encoder.Encode(resp)
	//	buf, _ := h.Encrypt(buffer.Bytes())
	//
	//	type responseEncrypted struct {
	//		Result string `json:"result"`
	//	}
	//
	//	resp = responseEncrypted{
	//		Result: buf,
	//	}
	//}
	buffer := &bytes.Buffer{}
	encoder := jsoniter.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(resp)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	buf := buffer.Bytes()
	c.Writer.Header().Add("Content-Length", fmt.Sprintf("%v", len(buf)))
	//traceId := trace.GetTraceIdFromCtx(c)
	//if traceId != "" {
	//	c.Writer.Header().Add("X-B3-Traceid", traceId)
	//}
	nWriteIndex := 0
	for nWriteIndex < len(buf) {
		n, err := c.Writer.Write(buf[nWriteIndex:])
		if err != nil {
			log.Errorf("%v len:%v index:%v n:%v", err.Error(), len(buf), nWriteIndex, n)
			return
		}
		nWriteIndex += n
	}
	c.Writer.Flush()
}
