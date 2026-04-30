package controller

import (
	"strconv"

	"github.com/FPRPL26/rpl-be/internal/api/service"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"github.com/FPRPL26/rpl-be/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type (
	SkillController interface {
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
	}

	skillController struct {
		skillService service.SkillService
	}
)

func NewSkillController(skillService service.SkillService) SkillController {
	return &skillController{skillService}
}

func (c *skillController) GetAll(ctx *gin.Context) {
	metaReq := meta.New(ctx)
	res, metaRes, err := c.skillService.GetAll(ctx.Request.Context(), metaReq)
	if err != nil {
		response.NewFailed("failed to get skills", err).Send(ctx)
		return
	}

	response.NewSuccess("Skills retrieved successfully", res, metaRes).Send(ctx)
}

func (c *skillController) GetById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.NewFailed("invalid skill id", err).Send(ctx)
		return
	}

	res, err := c.skillService.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed to get skill", err).Send(ctx)
		return
	}

	response.NewSuccess("Skill retrieved successfully", res).Send(ctx)
}
