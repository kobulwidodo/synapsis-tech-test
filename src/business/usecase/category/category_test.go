package category_test

import (
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

	tests := []struct {
		name     string
		mockFunc func(mock mockFields)
		want     []entity.Category
		wantErr  bool
	}{
		{
			name: "failed to get all menu",
			mockFunc: func(mock mockFields) {
				mock.category.EXPECT().GetList().Return([]entity.Category{}, assert.AnError)
			},
			want:    []entity.Category{},
			wantErr: true,
		},
		{
			name: "all ok",
			mockFunc: func(mock mockFields) {
				mock.category.EXPECT().GetList().Return(categoryOkMock, nil)
			},
			want:    categoryOkMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks)
			got, err := c.GetList()
			if (err != nil) != tt.wantErr {
				t.Errorf("category.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
