package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create Order
// @Description Create New Order
// @Security BearerAuth
// @Tags Transaction
// @Param transaction body entity.CreateTransactionParam true "transaction info"
// @Produce json
// @Success 200 {object} entity.Response{data=int}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction [POST]
func (r *rest) CreateOrder(ctx *gin.Context) {
	var inputParam entity.CreateTransactionParam
	if err := ctx.ShouldBindJSON(&inputParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	id, err := r.uc.Transaction.Create(ctx.Request.Context(), inputParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully created new order", gin.H{"id": id})
}
