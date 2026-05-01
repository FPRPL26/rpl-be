package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ClassTransaction(app *gin.Engine, classTransactionController controller.ClassTransactionController, reviewController controller.ReviewController, middleware middleware.Middleware) {
	// Public route for Midtrans callback
	app.POST("/api/transactions/midtrans-callback", classTransactionController.MidtransCallback)

	routes := app.Group("/api/transactions")
	{
		routes.Use(middleware.Authenticate())

		// Class specific transactions
		classes := routes.Group("/classes")
		{
			classes.POST("", classTransactionController.Checkout)
			classes.GET("", classTransactionController.GetAll)
			classes.POST("/:transaction_id/complete", classTransactionController.Complete)
		}
	}
}
