package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Class(app *gin.Engine, classController controller.ClassController, middleware middleware.Middleware) {
	routes := app.Group("/api/classes")
	{
		routes.GET("", classController.GetAll)
		routes.GET("/:class_id", classController.GetById)
		routes.GET("/:class_id/schedules", classController.GetSchedules)
		
		// Protected routes
		protected := routes.Group("")
		protected.Use(middleware.Authenticate())
		{
			protected.POST("", classController.Create)
			protected.PUT("/:class_id", classController.Update)
			protected.DELETE("/:class_id", classController.Delete)
			
			protected.POST("/:class_id/schedules", classController.AddSchedules)
			protected.PUT("/schedules/:schedule_id", classController.UpdateSchedule)
			protected.DELETE("/schedules/:schedule_id", classController.DeleteSchedule)
		}
	}
}
