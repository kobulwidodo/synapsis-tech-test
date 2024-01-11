package category

import (
	"context"
	"database/sql"
	"encoding/json"
	"go-clean/src/business/entity"
	"regexp"
	"testing"
	"time"

	"go-clean/src/lib/redis"
	mock_redis "go-clean/src/lib/tests/mock/redis"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_category_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `categories` WHERE `categories`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	mockRedis := mock_redis.NewMockInterface(ctrl)

	categoriesMock := []entity.Category{
		{
			Name: "category 1",
		},
	}
	marshalledCategories, _ := json.Marshal(categoriesMock)
	stringCategoriesMock := string(marshalledCategories)

	type mockFields struct {
		redis *mock_redis.MockInterface
	}

	mocks := mockFields{
		redis: mockRedis,
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		mockFunc    func(mock mockFields)
		want        []entity.Category
		wantErr     bool
	}{
		{
			name: "success get with cache",
			args: args{
				ctx: context.Background(),
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), getCategoryList).Return(stringCategoriesMock, nil)
			},
			want:    categoriesMock,
			wantErr: false,
		},
		{
			name: "failed to get cache list",
			args: args{
				ctx: context.Background(),
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), getCategoryList).Return("", assert.AnError)
			},
			want:    []entity.Category{},
			wantErr: true,
		},
		{
			name: "failed to exec query",
			args: args{
				ctx: context.Background(),
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), getCategoryList).Return("", redis.Nil)
			},
			want:    []entity.Category{},
			wantErr: true,
		},
		{
			name: "all ok but failed to set cache",
			args: args{
				ctx: context.Background(),
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("category 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), getCategoryList).Return("", redis.Nil)
				mock.redis.EXPECT().SetEX(context.Background(), getCategoryList, string(marshalledCategories), time.Minute).Return(assert.AnError)
			},
			want: []entity.Category{
				{
					Name: "category 1",
				},
			},
			wantErr: false,
		},
		{
			name: "all ok",
			args: args{
				ctx: context.Background(),
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("category 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), getCategoryList).Return("", redis.Nil)
				mock.redis.EXPECT().SetEX(context.Background(), getCategoryList, string(marshalledCategories), time.Minute).Return(nil)
			},
			want: []entity.Category{
				{
					Name: "category 1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks)
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient, mockRedis)
			got, err := u.GetList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("category.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
