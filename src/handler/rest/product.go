package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary Get List Product
// @Description Get List All Product
// @Security BearerAuth
// @Tags Product
// @Produce json
// @Param category_id query int false "category id param"
// @Success 200 {object} entity.Response{data=[]entity.Product{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/product [GET]
func (r *rest) GetListProduct(ctx *gin.Context) {
	var productParam entity.ProductParam
	if err := ctx.ShouldBindWith(&productParam, binding.Query); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	products, err := r.uc.Product.GetList(productParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get list all product", products)
}

// @Summary Get Product
// @Description Get a Product
// @Security BearerAuth
// @Tags Product
// @Produce json
// @Param product_id path int true "product id param"
// @Success 200 {object} entity.Response{data=[]entity.Product{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/product/{product_id} [GET]
func (r *rest) GetProduct(ctx *gin.Context) {
	var productParam entity.ProductParam
	if err := ctx.ShouldBindUri(&productParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	product, err := r.uc.Product.Get(productParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get a product", product)
}
