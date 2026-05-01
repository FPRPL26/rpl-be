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
	ClassRequestTutorApplicationController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		UpdateStatus(ctx *gin.Context)
	}

	classRequestTutorApplicationController struct {
		service service.ClassRequestTutorApplicationService
	}
)

func NewClassRequestTutorApplicationController(service service.ClassRequestTutorApplicationService) ClassRequestTutorApplicationController {
	return &classRequestTutorApplicationController{service}
}

func (c *classRequestTutorApplicationController) Create(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	var req dto.CreateClassRequestTutorApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.CreateClassRequestTutorApplicationRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.service.Apply(ctx.Request.Context(), userID, req)
	if err != nil {
		response.NewFailed("failed create class request tutor application", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request tutor application created successfully", res).Send(ctx)
}

func (c *classRequestTutorApplicationController) GetAll(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	requestID := ctx.Query("request_id")
	tutorProfileID := ctx.Query("tutor_profile_id")

	if requestID != "" {
		res, err := c.service.GetAllByRequest(ctx.Request.Context(), userID, requestID)
		if err != nil {
			response.NewFailed("failed to get applications", err).Send(ctx)
			return
		}
		response.NewSuccess("Class request applications retrieved successfully", res).Send(ctx)
		return
	}

	if tutorProfileID != "" {
		if tutorProfileID != userID {
			response.NewFailed("unauthorized", myerror.New("unauthorized tutor profile", 403)).Send(ctx)
			return
		}
		res, err := c.service.GetAllByTutorProfile(ctx.Request.Context(), tutorProfileID)
		if err != nil {
			response.NewFailed("failed to get applications", err).Send(ctx)
			return
		}
		response.NewSuccess("Class request applications retrieved successfully", res).Send(ctx)
		return
	}

	response.NewFailed("missing query parameters", myerror.New("request_id or tutor_profile_id required", 400)).Send(ctx)
}

func (c *classRequestTutorApplicationController) GetById(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	id := ctx.Param("id")
	res, err := c.service.GetById(ctx.Request.Context(), userID, id)
	if err != nil {
		response.NewFailed("failed to get class request tutor application", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request tutor application retrieved successfully", res).Send(ctx)
}

func (c *classRequestTutorApplicationController) UpdateStatus(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	id := ctx.Param("id")
	var req dto.UpdateClassRequestTutorApplicationStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.UpdateClassRequestTutorApplicationStatusRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.service.UpdateStatus(ctx.Request.Context(), userID, id, req)
	if err != nil {
		response.NewFailed("failed update class request tutor application status", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request tutor application status updated successfully", res).Send(ctx)
}
