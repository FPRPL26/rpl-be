package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ClassRequestTutorApplication(app *gin.Engine, controller controller.ClassRequestTutorApplicationController, middleware middleware.Middleware) {
	routes := app.Group("/api/class-request-tutor-applications")
	routes.Use(middleware.Authenticate())
	{
		routes.GET("", controller.GetAll)
		routes.GET(":id", controller.GetById)
		routes.POST("", controller.Create)
		routes.PATCH(":id/status", controller.UpdateStatus)
	}
}
