package routes

import (
	"net/http"

	"github.com/doyeon0307/tickit-backend/config"
	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/handler"

	"github.com/gin-gonic/gin"
)

type HandlerContainer struct {
	TicketUsecase   domain.TicketUsecase
	ScheduleUsecase domain.ScheduleUsecase
	S3Config        config.S3Config
}

func SetupRouter(handlers HandlerContainer) *gin.Engine {
	router := gin.Default()

	config.SetUpSwagger(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	v1 := router.Group("/api")
	{
		v1.GET("/health", healthCheck)
		handler.NewTicketHandler(v1, handlers.TicketUsecase, &handlers.S3Config)
		handler.NewScheduleHandler(v1, handlers.ScheduleUsecase)
	}

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
