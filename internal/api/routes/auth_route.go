package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Auth(app *gin.Engine, authcontroller controller.AuthController, middleware middleware.Middleware) {
	routes := app.Group("/api/auth")
	{
		routes.POST("/login", authcontroller.Login)
		routes.POST("/register", authcontroller.Register)
		routes.POST("/forget", authcontroller.ForgetPassword)
		routes.POST("/change", authcontroller.ChangePassword)
		routes.GET("/verify-email", authcontroller.VerifyEmail)
		routes.POST("/send-email-verification", authcontroller.SendEmailVerification)
		routes.GET("/refresh-token", authcontroller.RefreshToken)
		routes.POST("/logout", authcontroller.Logout)
		routes.GET("/me", middleware.Authenticate(), authcontroller.Me)
	}
}
