package controller

import (
	"github.com/FPRPL26/rpl-be/internal/api/service"
	"github.com/FPRPL26/rpl-be/internal/dto"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"github.com/FPRPL26/rpl-be/internal/pkg/response"
	"github.com/FPRPL26/rpl-be/internal/utils"
	"github.com/gin-gonic/gin"
)

type (
	ClassTransactionController interface {
		Checkout(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		Complete(ctx *gin.Context)
		MidtransCallback(ctx *gin.Context)
	}

	classTransactionController struct {
		transactionService service.ClassTransactionService
	}
)

func NewClassTransactionController(transactionService service.ClassTransactionService) ClassTransactionController {
	return &classTransactionController{transactionService}
}

func (c *classTransactionController) Checkout(ctx *gin.Context) {
	var req dto.CheckoutClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed to bind request", err).Send(ctx)
		return
	}

	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	res, err := c.transactionService.Checkout(ctx.Request.Context(), userID, req)
	if err != nil {
		response.NewFailed("failed to checkout class", err).Send(ctx)
		return
	}

	response.NewSuccess("Transaction created successfully", res).Send(ctx)
}

func (c *classTransactionController) GetAll(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	res, metaRes, err := c.transactionService.GetAllByUserId(ctx.Request.Context(), userID, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed to get transactions", err).Send(ctx)
		return
	}

	response.NewSuccess("Transactions retrieved successfully", res, metaRes).Send(ctx)
}

func (c *classTransactionController) Complete(ctx *gin.Context) {
	transactionID := ctx.Param("transaction_id")

	err := c.transactionService.Complete(ctx.Request.Context(), transactionID)
	if err != nil {
		response.NewFailed("failed to complete class transaction", err).Send(ctx)
		return
	}

	response.NewSuccess("Class completed. Payment released to tutor.", nil).Send(ctx)
}

func (c *classTransactionController) MidtransCallback(ctx *gin.Context) {
	var payload map[string]interface{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.NewFailed("failed to bind request", err).Send(ctx)
		return
	}

	err := c.transactionService.HandleMidtransCallback(ctx.Request.Context(), payload)
	if err != nil {
		response.NewFailed("failed to handle midtrans callback", err).Send(ctx)
		return
	}

	response.NewSuccess("Midtrans callback handled successfully", nil).Send(ctx)
}
