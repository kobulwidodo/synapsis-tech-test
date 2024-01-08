package product_test

import (
	mock_product "go-clean/src/business/domain/mock/product"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/product"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_product_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productMock := mock_product.NewMockInterface(ctrl)

	productParamMock := entity.ProductParam{}

	productOkResult := []entity.Product{
		{
			Name: "product 1",
		},
	}

	p := product.Init(productMock)

	type mockFields struct {
		product *mock_product.MockInterface
	}
	mocks := mockFields{
		product: productMock,
	}

	type args struct {
		param entity.ProductParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     []entity.Product
		wantErr  bool
	}{
		{
			name: "failed to get products",
			args: args{
				param: productParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.product.EXPECT().GetList(productParamMock).Return([]entity.Product{}, assert.AnError)
			},
			want:    []entity.Product{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: productParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.product.EXPECT().GetList(productParamMock).Return(productOkResult, nil)
			},
			want:    productOkResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := p.GetList(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_product_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productMock := mock_product.NewMockInterface(ctrl)

	productParamMock := entity.ProductParam{}

	productOkResult := entity.Product{
		Name: "product 1",
	}

	p := product.Init(productMock)

	type mockFields struct {
		product *mock_product.MockInterface
	}
	mocks := mockFields{
		product: productMock,
	}

	type args struct {
		param entity.ProductParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     entity.Product
		wantErr  bool
	}{
		{
			name: "failed to get products",
			args: args{
				param: productParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.product.EXPECT().Get(productParamMock).Return(entity.Product{}, assert.AnError)
			},
			want:    entity.Product{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: productParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.product.EXPECT().Get(productParamMock).Return(productOkResult, nil)
			},
			want:    productOkResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := p.Get(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
