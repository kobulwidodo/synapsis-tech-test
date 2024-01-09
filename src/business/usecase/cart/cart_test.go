package cart_test

import (
	"context"
	mock_cart "go-clean/src/business/domain/mock/cart"
	mock_product "go-clean/src/business/domain/mock/product"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/cart"
	"go-clean/src/lib/auth"
	"testing"

	mock_auth "go-clean/src/lib/tests/mock/auth"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_cart_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	productMock := mock_product.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)

	createCartParamMock := entity.CreateCartParam{
		ProductID: 1,
		Qty:       1,
	}

	userAuthMock := auth.UserAuthInfo{
		User: auth.User{
			ID: 1,
		},
	}

	productParamMock := entity.ProductParam{
		ID: 1,
	}

	productResultMock := entity.Product{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "product 1",
	}

	cartParamMock := entity.CartParam{
		UserID:    1,
		ProductID: 1,
		Status:    entity.StatusInCart,
	}

	cartResultMock := entity.Cart{
		Model: gorm.Model{
			ID: 1,
		},
		ProductID: 1,
		UserID:    1,
		Qty:       1,
	}

	cartUpdateParamMock := entity.CartParam{
		UserID:    1,
		ProductID: 1,
		Status:    entity.StatusInCart,
	}

	cartUpdateMock := entity.UpdateCartParam{
		Qty: 2,
	}

	createCartMock := entity.Cart{
		UserID:    1,
		ProductID: 1,
		Qty:       1,
		Status:    entity.StatusInCart,
	}

	c := cart.Init(cartMock, authMock, productMock)

	type mockFields struct {
		auth    *mock_auth.MockInterface
		product *mock_product.MockInterface
		cart    *mock_cart.MockInterface
	}

	mocks := mockFields{
		auth:    authMock,
		product: productMock,
		cart:    cartMock,
	}

	type args struct {
		ctx    context.Context
		params entity.CreateCartParam
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     entity.Cart
		wantErr  bool
	}{
		{
			name: "failed to get user auth info",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.UserAuthInfo{}, assert.AnError)
			},
			want:    entity.Cart{},
			wantErr: true,
		},
		{
			name: "failed to get product",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.product.EXPECT().Get(productParamMock).Return(entity.Product{}, assert.AnError)
			},
			want:    entity.Cart{},
			wantErr: true,
		},
		{
			name: "failed to update cart",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.product.EXPECT().Get(productParamMock).Return(productResultMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
				mock.cart.EXPECT().Update(cartUpdateParamMock, cartUpdateMock).Return(assert.AnError)
			},
			want:    cartResultMock,
			wantErr: true,
		},
		{
			name: "all ok update cart",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.product.EXPECT().Get(productParamMock).Return(productResultMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
				mock.cart.EXPECT().Update(cartUpdateParamMock, cartUpdateMock).Return(nil)
			},
			want:    cartResultMock,
			wantErr: false,
		},
		{
			name: "failed to create new cart",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.product.EXPECT().Get(productParamMock).Return(productResultMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(entity.Cart{}, nil)
				mock.cart.EXPECT().Create(createCartMock).Return(entity.Cart{}, assert.AnError)
			},
			want:    entity.Cart{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.product.EXPECT().Get(productParamMock).Return(productResultMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(entity.Cart{}, nil)
				mock.cart.EXPECT().Create(createCartMock).Return(createCartMock, nil)
			},
			want:    createCartMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := c.Create(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
