package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
