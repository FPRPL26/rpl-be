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
	ClassRequestTransactionController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Complete(ctx *gin.Context)
		MidtransCallback(ctx *gin.Context)
	}

	classRequestTransactionController struct {
		service service.ClassRequestTransactionService
	}
)

func NewClassRequestTransactionController(service service.ClassRequestTransactionService) ClassRequestTransactionController {
	return &classRequestTransactionController{service}
}

func (c *classRequestTransactionController) Create(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	var req dto.CreateClassRequestTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.CreateClassRequestTransactionRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), userID, req)
	if err != nil {
		response.NewFailed("failed create class request transaction", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request transaction created successfully", res).Send(ctx)
}

func (c *classRequestTransactionController) GetAll(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	res, metaRes, err := c.service.GetAllByUserId(ctx.Request.Context(), userID, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed to get class request transactions", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request transactions retrieved successfully", res, metaRes).Send(ctx)
}

func (c *classRequestTransactionController) GetById(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	id := ctx.Param("id")
	res, err := c.service.GetById(ctx.Request.Context(), userID, id)
	if err != nil {
		response.NewFailed("failed to get class request transaction", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request transaction retrieved successfully", res).Send(ctx)
}

func (c *classRequestTransactionController) Update(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	id := ctx.Param("id")
	var req dto.UpdateClassRequestTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.UpdateClassRequestTransactionRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.service.UpdateStatus(ctx.Request.Context(), userID, id, req)
	if err != nil {
		response.NewFailed("failed update class request transaction", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request transaction updated successfully", res).Send(ctx)
}

func (c *classRequestTransactionController) Complete(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("unauthorized", err).Send(ctx)
		return
	}

	transactionID := ctx.Param("id")
	if err := c.service.Complete(ctx.Request.Context(), userID, transactionID); err != nil {
		response.NewFailed("failed to complete class request transaction", err).Send(ctx)
		return
	}

	response.NewSuccess("Class request transaction completed successfully", nil).Send(ctx)
}

func (c *classRequestTransactionController) MidtransCallback(ctx *gin.Context) {
	var payload map[string]interface{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.NewFailed("failed to bind request", err).Send(ctx)
		return
	}

	err := c.service.HandleMidtransCallback(ctx.Request.Context(), payload)
	if err != nil {
		response.NewFailed("failed to handle midtrans callback", err).Send(ctx)
		return
	}

	response.NewSuccess("Midtrans callback handled successfully", nil).Send(ctx)
}
