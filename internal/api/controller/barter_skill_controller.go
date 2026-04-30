package controller

import (
	"net/http"

	"github.com/FPRPL26/rpl-be/internal/api/service"
	"github.com/FPRPL26/rpl-be/internal/dto"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"github.com/FPRPL26/rpl-be/internal/pkg/response"
	"github.com/FPRPL26/rpl-be/internal/utils"
	"github.com/gin-gonic/gin"
)

type (
	BarterSkillController interface {
		CreateOffer(ctx *gin.Context)
		RequestOffer(ctx *gin.Context)
		ApproveRequest(ctx *gin.Context)
		GetAllOffers(ctx *gin.Context)
		GetOfferById(ctx *gin.Context)
		GetMyOffers(ctx *gin.Context)
		GetMyRequests(ctx *gin.Context)
		GetIncomingRequests(ctx *gin.Context)
	}

	barterSkillController struct {
		barterService service.BarterSkillService
	}
)

func NewBarterSkillController(barterService service.BarterSkillService) BarterSkillController {
	return &barterSkillController{barterService}
}

func (c *barterSkillController) CreateOffer(ctx *gin.Context) {
	var req dto.CreateBarterOfferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed to bind request", myerror.GetErrBodyRequest(err, req)).Send(ctx)
		return
	}

	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	res, err := c.barterService.CreateOffer(ctx.Request.Context(), userID, req)
	if err != nil {
		response.NewFailed("failed to create barter offer", err).Send(ctx)
		return
	}

	response.NewSuccess("Barter offer created successfully", res).ChangeStatusCode(http.StatusCreated).Send(ctx)
}

func (c *barterSkillController) RequestOffer(ctx *gin.Context) {
	barterID := ctx.Param("barter_id")
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	res, err := c.barterService.RequestOffer(ctx.Request.Context(), userID, barterID)
	if err != nil {
		response.NewFailed("failed to request barter offer", err).Send(ctx)
		return
	}

	response.NewSuccess("Barter offer requested", res).Send(ctx)
}

func (c *barterSkillController) ApproveRequest(ctx *gin.Context) {
	barterID := ctx.Param("barter_id")
	var req dto.ApproveBarterRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed to bind request", err).Send(ctx)
		return
	}

	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	err = c.barterService.ApproveRequest(ctx.Request.Context(), userID, barterID, req)
	if err != nil {
		response.NewFailed("failed to approve barter request", err).Send(ctx)
		return
	}

	response.NewSuccess("Barter request approved successfully", nil).Send(ctx)
}

func (c *barterSkillController) GetAllOffers(ctx *gin.Context) {
	metaReq := meta.New(ctx)
	res, metaRes, err := c.barterService.GetAllOffers(ctx.Request.Context(), metaReq)
	if err != nil {
		response.NewFailed("failed to get barter offers", err).Send(ctx)
		return
	}

	response.NewSuccess("Barter offers retrieved successfully", res, metaRes).Send(ctx)
}

func (c *barterSkillController) GetOfferById(ctx *gin.Context) {
	barterID := ctx.Param("barter_id")
	res, err := c.barterService.GetOfferById(ctx.Request.Context(), barterID)
	if err != nil {
		response.NewFailed("failed to get barter offer", err).Send(ctx)
		return
	}

	response.NewSuccess("Barter offer retrieved successfully", res).Send(ctx)
}

func (c *barterSkillController) GetMyOffers(ctx *gin.Context) {
	metaReq := meta.New(ctx)
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	res, metaRes, err := c.barterService.GetMyOffers(ctx.Request.Context(), userID, metaReq)
	if err != nil {
		response.NewFailed("failed to get my barter offers", err).Send(ctx)
		return
	}

	response.NewSuccess("My barter offers retrieved successfully", res, metaRes).Send(ctx)
}

func (c *barterSkillController) GetMyRequests(ctx *gin.Context) {
	metaReq := meta.New(ctx)
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	res, metaRes, err := c.barterService.GetMyRequests(ctx.Request.Context(), userID, metaReq)
	if err != nil {
		response.NewFailed("failed to get my barter requests", err).Send(ctx)
		return
	}

	response.NewSuccess("My barter requests retrieved successfully", res, metaRes).Send(ctx)
}

func (c *barterSkillController) GetIncomingRequests(ctx *gin.Context) {
	metaReq := meta.New(ctx)
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed to get user id from context", err).Send(ctx)
		return
	}

	res, metaRes, err := c.barterService.GetIncomingRequests(ctx.Request.Context(), userID, metaReq)
	if err != nil {
		response.NewFailed("failed to get incoming barter requests", err).Send(ctx)
		return
	}

	response.NewSuccess("Incoming barter requests retrieved successfully", res, metaRes).Send(ctx)
}
