package category_test

import (
	"context"
	mock_category "go-clean/src/business/domain/mock/category"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/category"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_category_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	categoryMock := mock_category.NewMockInterface(ctrl)

	categoryOkMock := []entity.Category{
		{
			Name: "category 1",
		},
	}

	c := category.Init(categoryMock)

	type mockFields struct {
		category *mock_category.MockInterface
	}
	mocks := mockFields{
		category: categoryMock,
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields)
		want     []entity.Category
		wantErr  bool
	}{
		{
			name: "failed to get all menu",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields) {
				mock.category.EXPECT().GetList(context.Background()).Return([]entity.Category{}, assert.AnError)
			},
			want:    []entity.Category{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields) {
				mock.category.EXPECT().GetList(context.Background()).Return(categoryOkMock, nil)
			},
			want:    categoryOkMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks)
			got, err := c.GetList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("category.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
