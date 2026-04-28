package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ClassTransaction(app *gin.Engine, classTransactionController controller.ClassTransactionController, middleware middleware.Middleware) {
	// Public route for Midtrans callback
	app.POST("/api/transactions/midtrans-callback", classTransactionController.MidtransCallback)

	routes := app.Group("/api/transactions/classes")
	{
		routes.Use(middleware.Authenticate())
		routes.POST("", classTransactionController.Checkout)
		routes.GET("", classTransactionController.GetAll)
		routes.POST("/:transaction_id/complete", classTransactionController.Complete)
	}
}
