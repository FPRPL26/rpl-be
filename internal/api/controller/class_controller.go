package controller

import (
	"github.com/FPRPL26/rpl-be/internal/api/service"
	"github.com/FPRPL26/rpl-be/internal/dto"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"github.com/FPRPL26/rpl-be/internal/pkg/response"
	"github.com/FPRPL26/rpl-be/internal/utils"
	"github.com/gin-gonic/gin"
)

type (
	ClassController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		GetSchedules(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		AddSchedules(ctx *gin.Context)
		UpdateSchedule(ctx *gin.Context)
		DeleteSchedule(ctx *gin.Context)
	}

	classController struct {
		classService service.ClassService
	}
)

func NewClass(classService service.ClassService) ClassController {
	return &classController{
		classService: classService,
	}
}

func (c *classController) Create(ctx *gin.Context) {
	tutorId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	var req dto.CreateClassRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.CreateClassRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	result, err := c.classService.Create(ctx.Request.Context(), tutorId, req)
	if err != nil {
		response.NewFailed("failed create class", err).Send(ctx)
		return
	}

	resData := dto.CreateClassResponse{
		ID: result.ID,
	}

	response.NewSuccess("Class created successfully", resData).Send(ctx)
}

func (c *classController) GetAll(ctx *gin.Context) {
	classes, metaRes, err := c.classService.GetAll(ctx.Request.Context(), meta.New(ctx))
	if err != nil {
		response.NewFailed("failed to get classes", err).Send(ctx)
		return
	}

	response.NewSuccess("Classes retrieved successfully", classes, metaRes).Send(ctx)
}

func (c *classController) GetById(ctx *gin.Context) {
	classId := ctx.Param("class_id")
	result, err := c.classService.GetById(ctx.Request.Context(), classId)
	if err != nil {
		response.NewFailed("failed to get class", err).Send(ctx)
		return
	}

	response.NewSuccess("Class retrieved successfully", result).Send(ctx)
}

func (c *classController) GetSchedules(ctx *gin.Context) {
	classId := ctx.Param("class_id")
	schedules, metaRes, err := c.classService.GetSchedules(ctx.Request.Context(), meta.New(ctx), classId)
	if err != nil {
		response.NewFailed("failed to get schedules", err).Send(ctx)
		return
	}

	response.NewSuccess("Schedules retrieved successfully", schedules, metaRes).Send(ctx)
}

func (c *classController) AddSchedules(ctx *gin.Context) {
	classId := ctx.Param("class_id")

	var req dto.AddSchedulesRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.AddSchedulesRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	err := c.classService.AddSchedules(ctx.Request.Context(), classId, req)
	if err != nil {
		response.NewFailed("failed to add schedules", err).Send(ctx)
		return
	}

	response.NewSuccess("Schedules added successfully", nil).Send(ctx)
}

func (c *classController) Update(ctx *gin.Context) {
	tutorId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	classId := ctx.Param("class_id")
	var req dto.UpdateClassRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.UpdateClassRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	result, err := c.classService.Update(ctx.Request.Context(), tutorId, classId, req)
	if err != nil {
		response.NewFailed("failed update class", err).Send(ctx)
		return
	}

	response.NewSuccess("Class updated successfully", result).Send(ctx)
}

func (c *classController) Delete(ctx *gin.Context) {
	tutorId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	classId := ctx.Param("class_id")
	err = c.classService.Delete(ctx.Request.Context(), tutorId, classId)
	if err != nil {
		response.NewFailed("failed to delete class", err).Send(ctx)
		return
	}

	response.NewSuccess("Class deleted successfully", nil).Send(ctx)
}

func (c *classController) UpdateSchedule(ctx *gin.Context) {
	tutorId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	scheduleId := ctx.Param("schedule_id")
	var req dto.UpdateScheduleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.UpdateScheduleRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	result, err := c.classService.UpdateSchedule(ctx.Request.Context(), tutorId, scheduleId, req)
	if err != nil {
		response.NewFailed("failed update schedule", err).Send(ctx)
		return
	}

	response.NewSuccess("Schedule updated successfully", result).Send(ctx)
}

func (c *classController) DeleteSchedule(ctx *gin.Context) {
	tutorId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	scheduleId := ctx.Param("schedule_id")
	err = c.classService.DeleteSchedule(ctx.Request.Context(), tutorId, scheduleId)
	if err != nil {
		response.NewFailed("failed to delete schedule", err).Send(ctx)
		return
	}

	response.NewSuccess("Schedule deleted successfully", nil).Send(ctx)
}
