package routes

import (
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/gin-gonic/gin"
)

func Skill(app *gin.Engine, skillController controller.SkillController) {
	routes := app.Group("/api/skills")
	{
		routes.GET("", skillController.GetAll)
		routes.GET("/:id", skillController.GetById)
	}
}
