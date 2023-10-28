package middleware

import (
	"github.com/alexfaker/Pantasy/middleware/log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				connId, _ := c.Get("Conn-ID")
				log.Errorf("[%v] error: %s, stack: %s", connId, err, string(debug.Stack()))
				Response(c, StatusInternalServerError, nil)
			}
		}()
		c.Next()
	}
}
