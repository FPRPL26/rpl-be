package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Review(app *gin.Engine, reviewController controller.ReviewController, middleware middleware.Middleware) {
	routes := app.Group("/api/reviews")
	{
		routes.Use(middleware.Authenticate())
		routes.POST("", reviewController.Submit)
	}
}
