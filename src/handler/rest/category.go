package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get List Category
// @Description Get List All Category
// @Security BearerAuth
// @Tags Category
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Category{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/category [GET]
func (r *rest) GetListCategory(ctx *gin.Context) {
	categories, err := r.uc.Category.GetList(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get list all category", categories)
}
