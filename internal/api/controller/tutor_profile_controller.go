package controller

import (
	"strconv"

	"github.com/FPRPL26/rpl-be/internal/api/service"
	"github.com/FPRPL26/rpl-be/internal/dto"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/FPRPL26/rpl-be/internal/pkg/response"
	"github.com/FPRPL26/rpl-be/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	TutorController interface {
		Create(ctx *gin.Context)
		Update(ctx *gin.Context)
		GetByID(ctx *gin.Context)
		Delete(ctx *gin.Context)
		List(ctx *gin.Context)
	}

	tutorController struct {
		service service.TutorService
	}
)

func NewTutorController(service service.TutorService) TutorController {
	return &tutorController{service: service}
}

func (c *tutorController) Create(ctx *gin.Context) {
	var req dto.TutorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.TutorRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	userIDStr, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.NewFailed("invalid user_id format", err).Send(ctx)
		return
	}

	res, err := c.service.CreateTutor(ctx.Request.Context(), userID, req)
	if err != nil {
		response.NewFailed("failed create tutor", err).Send(ctx)
		return
	}

	response.NewSuccess("Tutor created successfully", res).Send(ctx)
}

func (c *tutorController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.NewFailed("invalid id format", err).Send(ctx)
		return
	}
	var req dto.TutorUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.TutorUpdateRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.service.UpdateTutor(ctx.Request.Context(), id, req)
	if err != nil {
		response.NewFailed("failed update tutor", err).Send(ctx)
		return
	}

	response.NewSuccess("Tutor updated successfully", res).Send(ctx)
}

func (c *tutorController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.NewFailed("invalid id format", err).Send(ctx)
		return
	}
	res, err := c.service.GetTutorByID(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed get tutor", err).Send(ctx)
		return
	}

	response.NewSuccess("Tutor retrieved successfully", res).Send(ctx)
}

func (c *tutorController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.NewFailed("invalid id format", err).Send(ctx)
		return
	}
	if err := c.service.DeleteTutor(ctx.Request.Context(), id); err != nil {
		response.NewFailed("failed to delete tutor", err).Send(ctx)
		return
	}

	response.NewSuccess("Tutor deleted successfully", nil).Send(ctx)
}

func (c *tutorController) List(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	res, err := c.service.ListTutors(ctx.Request.Context(), limit, offset)
	if err != nil {
		response.NewFailed("failed to list tutors", err).Send(ctx)
		return
	}

	response.NewSuccess("Tutors retrieved successfully", res).Send(ctx)
}
