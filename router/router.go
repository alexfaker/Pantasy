package router

import (
	"github.com/alexfaker/Pantasy/config"
	"github.com/alexfaker/Pantasy/middleware"
	"github.com/alexfaker/Pantasy/middleware/log"
	"github.com/gin-gonic/gin"
	"time"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(
		middleware.TimeoutMiddleware(time.Second*config.HTTPRequestTimeout),
		middleware.RequestLog(),
		middleware.CrossDomain(),
		middleware.RecoveryMiddleware(),
		middleware.Validator(),
	)
	apiGroup := router.Group("/")
	apiGroup.Any("/", func(c *gin.Context) {
		c.JSON(200, nil)
	})

	app := router.Group("/next/app")
	{
		devGroup := app.Group("/dev")
		{
			devGroup.GET("/ping", func(c *gin.Context) {
				c.JSON(200, "ping success")
			})
		}

	}

	//curEnv := os.Getenv("env")
	//if curEnv == "test" || curEnv == "dev" || curEnv == "local" {
	//	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//}
	if err := router.Run(":9010"); err != nil {
		log.Errorf("%v", err)
	}
}
