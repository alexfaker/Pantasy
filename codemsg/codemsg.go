package codemsg

/**
 * 网关专用状态码
 */

const (
	// StatusOK ok
	StatusOK = 200
	// StatusBadRequest bad request
	StatusBadRequest = 400
	// StatusUnauthorized unauthorized
	StatusUnauthorized = 401
	// StatusNotFound not found
	StatusNotFound = 404
	// StatusInternalServerError internal server error
	StatusInternalServerError = 500
)

var codeMsg = map[int]string{
	StatusOK:                  "success",
	StatusBadRequest:          "param error",
	StatusUnauthorized:        "verify error",
	StatusNotFound:            "not found",
	StatusInternalServerError: "service error",
}

// Msg msg
func Msg(code int) string {
	msg, _ := codeMsg[code]
	return msg
}
