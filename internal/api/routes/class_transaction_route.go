package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ClassTransaction(app *gin.Engine, classTransactionController controller.ClassTransactionController, middleware middleware.Middleware) {
	routes := app.Group("/api/transactions/classes")
	{
		routes.POST("", middleware.Authenticate(), classTransactionController.Checkout)
		routes.POST("/:transaction_id/complete", middleware.Authenticate(), classTransactionController.Complete)
	}
}
