package rest

import (
	"encoding/json"
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Payment Detail
// @Description Get Payment Detail by Transaction ID
// @Security BearerAuth
// @Tags Midtrans Transaction
// @Produce json
// @Param transaction_id path integer true "transaction id"
// @Success 200 {object} entity.Response{data=entity.MidtransTransactionPaymentDetail{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction/{transaction_id}/payment-detail [GET]
func (r *rest) GetPaymentDetail(ctx *gin.Context) {
	var param entity.MidtransTransactionParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.MidtransTransaction.GetPaymentDetail(param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get payment detail", result)
}

func (r *rest) HandleNotification(ctx *gin.Context) {
	var notifPayload map[string]interface{}
	err := json.NewDecoder(ctx.Request.Body).Decode(&notifPayload)
	if err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := r.uc.MidtransTransaction.HandleNotification(notifPayload); err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully handle transaction", nil)
}
