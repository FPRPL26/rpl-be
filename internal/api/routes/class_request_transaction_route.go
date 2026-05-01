package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ClassRequestTransaction(app *gin.Engine, classRequestTransactionController controller.ClassRequestTransactionController, middleware middleware.Middleware) {
	app.POST("/api/class-request-transactions/midtrans-callback", classRequestTransactionController.MidtransCallback)

	routes := app.Group("/api/class-request-transactions")
	routes.Use(middleware.Authenticate())
	{
		routes.POST("", classRequestTransactionController.Create)
		routes.GET("", classRequestTransactionController.GetAll)
		routes.GET("/:id", classRequestTransactionController.GetById)
		routes.POST("/:id/complete", classRequestTransactionController.Complete)
		routes.PATCH("/:id", classRequestTransactionController.Update)
	}
}
