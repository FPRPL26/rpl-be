package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Tutor(r *gin.Engine, c controller.TutorController, middleware middleware.Middleware) {
	routes := r.Group("/api/tutors")
	{
		routes.GET("", c.List)
		routes.GET("/:id", c.GetByID)

		protected := routes.Group("")
		protected.Use(middleware.Authenticate())
		{
			protected.POST("/upgrade", c.Create)
			protected.PATCH("/:id", c.Update)
			protected.DELETE("/:id", c.Delete)
		}
	}
}
