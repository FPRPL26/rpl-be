package controller

import (
	"net/http"

	"github.com/FPRPL26/rpl-be/internal/api/service"
	"github.com/FPRPL26/rpl-be/internal/dto"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/FPRPL26/rpl-be/internal/pkg/response"
	"github.com/FPRPL26/rpl-be/internal/utils"
	"github.com/gin-gonic/gin"
)

type (
	ReviewController interface {
		Submit(ctx *gin.Context)
	}

	reviewController struct {
		reviewService service.ReviewService
	}
)

func NewReviewController(reviewService service.ReviewService) ReviewController {
	return &reviewController{reviewService}
}

func (c *reviewController) Submit(ctx *gin.Context) {
	var req dto.SubmitReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.LoginRequest{})
		response.NewFailed("failed to bind request", err).Send(ctx)
		return
	}

	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	res, err := c.reviewService.SubmitReview(ctx.Request.Context(), userID, req)
	if err != nil {
		response.NewFailed("failed to submit review", err).Send(ctx)
		return
	}

	response.NewSuccess("Review submitted successfully", res).ChangeStatusCode(http.StatusCreated).Send(ctx)
}
