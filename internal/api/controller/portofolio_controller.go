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
	PortofolioController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		GetMyPortofolios(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	portofolioController struct {
		service service.PortofolioService
	}
)

func NewPortofolioController(service service.PortofolioService) PortofolioController {
	return &portofolioController{service: service}
}

func (c *portofolioController) Create(ctx *gin.Context) {
	var req dto.CreatePortofolioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err := myerror.GetErrBodyRequest(err, dto.CreatePortofolioRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	tutorProfileID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), tutorProfileID, req)
	if err != nil {
		response.NewFailed("failed create portofolio", err).Send(ctx)
		return
	}

	response.NewSuccess("Portofolio created successfully", res).Send(ctx)
}

func (c *portofolioController) GetAll(ctx *gin.Context) {
	tutorProfileID := ctx.Query("tutor_profile_id")

	res, err := c.service.GetAll(ctx.Request.Context(), tutorProfileID)
	if err != nil {
		response.NewFailed("failed get portofolios", err).Send(ctx)
		return
	}

	response.NewSuccess("Portofolios retrieved successfully", dto.PortofolioListResponse{Data: res}).Send(ctx)
}

func (c *portofolioController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.service.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed get portofolio", err).Send(ctx)
		return
	}

	response.NewSuccess("Portofolio retrieved successfully", res).Send(ctx)
}

func (c *portofolioController) GetMyPortofolios(ctx *gin.Context) {
	tutorProfileID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	res, err := c.service.GetAll(ctx.Request.Context(), tutorProfileID)
	if err != nil {
		response.NewFailed("failed get my portofolios", err).Send(ctx)
		return
	}

	response.NewSuccess("Portofolios retrieved successfully", dto.PortofolioListResponse{Data: res}).Send(ctx)
}

func (c *portofolioController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdatePortofolioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err := myerror.GetErrBodyRequest(err, dto.UpdatePortofolioRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	tutorProfileID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	res, err := c.service.Update(ctx.Request.Context(), tutorProfileID, id, req)
	if err != nil {
		response.NewFailed("failed update portofolio", err).Send(ctx)
		return
	}

	response.NewSuccess("Portofolio updated successfully", res).Send(ctx)
}

func (c *portofolioController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	tutorProfileID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), tutorProfileID, id); err != nil {
		response.NewFailed("failed delete portofolio", err).Send(ctx)
		return
	}

	response.NewSuccess("Portofolio deleted successfully", nil).Send(ctx)
}
