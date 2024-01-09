package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Add Product to Cart
// @Description Add product to cart
// @Security BearerAuth
// @Tags Cart
// @Param user body entity.CreateCartParam true "user info"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.Cart{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/cart [POST]
func (r *rest) CreateCart(ctx *gin.Context) {
	var param entity.CreateCartParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	cart, err := r.uc.Cart.Create(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfullt add item to cart", cart)
}

// @Summary Get List Cart
// @Description Get List Cart
// @Security BearerAuth
// @Tags Cart
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Cart{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/cart [GET]
func (r *rest) GetListCart(ctx *gin.Context) {
	carts, err := r.uc.Cart.GetList(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfullt get all product from cart", carts)
}

// @Summary Delete a Product
// @Description Delete a Product from Cart
// @Security BearerAuth
// @Tags Cart
// @Produce json
// @Param cart_id path int true "cart id"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/cart/{cart_id} [DELETE]
func (r *rest) DeleteCart(ctx *gin.Context) {
	var param entity.CartParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := r.uc.Cart.Delete(ctx.Request.Context(), param); err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfullt delete the product from cart", nil)
}
