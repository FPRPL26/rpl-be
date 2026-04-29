package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ClassRequest(app *gin.Engine, classRequestController controller.ClassRequestController, middleware middleware.Middleware) {
	routes := app.Group("/api/class-requests")
	routes.GET("", classRequestController.GetAll)
	routes.GET("/:id", classRequestController.GetById)

	protected := routes.Group("")
	protected.Use(middleware.Authenticate())
	{
		protected.POST("", classRequestController.Create)
		protected.PATCH("/:id", classRequestController.Update)
		protected.DELETE("/:id", classRequestController.Delete)
	}
}
