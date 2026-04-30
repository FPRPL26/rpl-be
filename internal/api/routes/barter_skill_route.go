package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func BarterSkill(app *gin.Engine, barterController controller.BarterSkillController, middleware middleware.Middleware) {
	routes := app.Group("/api/barters")
	{
		routes.GET("", barterController.GetAllOffers)
		routes.GET("/:barter_id", barterController.GetOfferById)

		routes.Use(middleware.Authenticate())
		routes.GET("/me/offers", barterController.GetMyOffers)
		routes.GET("/me/requests", barterController.GetMyRequests)
		routes.GET("/me/incoming-requests", barterController.GetIncomingRequests)
		routes.POST("", barterController.CreateOffer)
		routes.POST("/:barter_id/request", barterController.RequestOffer)
		routes.POST("/:barter_id/approve", barterController.ApproveRequest)
	}
}
