package controller

import (
	"github.com/FPRPL26/rpl-be/internal/api/service"
	"github.com/FPRPL26/rpl-be/internal/dto"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/FPRPL26/rpl-be/internal/pkg/response"
	"github.com/FPRPL26/rpl-be/internal/utils"
	"github.com/gin-gonic/gin"
)

type (
	ClassRequestController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	classRequestController struct {
		svc service.ClassRequestService
	}
)

func NewClassRequest(svc service.ClassRequestService) ClassRequestController {
	return &classRequestController{svc}
}

func (c *classRequestController) Create(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	var req dto.CreateClassRequestRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.CreateClassRequestRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.svc.Create(ctx.Request.Context(), userId, req)
	if err != nil {
		response.NewFailed("failed create class request", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request created successfully", res).Send(ctx)
}

func (c *classRequestController) GetAll(ctx *gin.Context) {
	items, err := c.svc.GetAll(ctx.Request.Context())
	if err != nil {
		response.NewFailed("failed to get class requests", err).Send(ctx)
		return
	}

	response.NewSuccess("Class requests retrieved successfully", items).Send(ctx)
}

func (c *classRequestController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.svc.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed to get class request", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request retrieved successfully", res).Send(ctx)
}

func (c *classRequestController) Update(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	id := ctx.Param("id")
	var req dto.UpdateClassRequestRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.UpdateClassRequestRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.svc.Update(ctx.Request.Context(), userId, id, req)
	if err != nil {
		response.NewFailed("failed update class request", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request updated successfully", res).Send(ctx)
}

func (c *classRequestController) Delete(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	id := ctx.Param("id")
	if err := c.svc.Delete(ctx.Request.Context(), userId, id); err != nil {
		response.NewFailed("failed delete class request", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request deleted successfully", nil).Send(ctx)
}
