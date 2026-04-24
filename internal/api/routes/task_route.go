package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Task(app *gin.Engine, taskcontroller controller.TaskController, middleware middleware.Middleware) {
	routes := app.Group("/api/task")
	{
		routes.POST("", taskcontroller.Create)
		routes.GET("", taskcontroller.GetAll)
		routes.GET("/:id", taskcontroller.GetById)
		routes.PUT("/:id", taskcontroller.Update)
		routes.DELETE("/:id", taskcontroller.Delete)
	}
}
