package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Portofolio(app *gin.Engine, c controller.PortofolioController, middleware middleware.Middleware) {
	routes := app.Group("/api/portofolios")
	{
		routes.GET("", c.GetAll)
		routes.GET("/me", middleware.Authenticate(), c.GetMyPortofolios)
		routes.GET("/tutor/:tutor_profile_id", c.GetByTutorProfile)
		routes.GET("/:id", c.GetById)

		protected := routes.Group("")
		protected.Use(middleware.Authenticate())
		{
			protected.POST("", c.Create)
			protected.PUT("/:id", c.Update)
			protected.DELETE("/:id", c.Delete)
		}
	}
}
