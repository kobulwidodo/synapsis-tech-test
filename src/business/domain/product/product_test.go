package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-clean/src/business/entity"
	"go-clean/src/lib/redis"
	mock_redis "go-clean/src/lib/tests/mock/redis"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_product_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `products` WHERE `products`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	mockRedis := mock_redis.NewMockInterface(ctrl)

	mockParam := entity.ProductParam{}
	marshalledParam, _ := json.Marshal(mockParam)

	mockProductResult := []entity.Product{
		{
			Name: "product 1",
		},
	}
	marshalledResult, _ := json.Marshal(mockProductResult)

	type mockFields struct {
		redis *mock_redis.MockInterface
	}

	mocks := mockFields{
		redis: mockRedis,
	}

	type args struct {
		ctx   context.Context
		param entity.ProductParam
	}

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		mockFunc    func(mock mockFields)
		want        []entity.Product
		wantErr     bool
	}{
		{
			name: "success get from cache",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return(string(marshalledResult), nil)
			},
			want:    mockProductResult,
			wantErr: false,
		},
		{
			name: "failed to get cache list",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return("", assert.AnError)
			},
			want:    []entity.Product{},
			wantErr: true,
		},
		{
			name: "failed to exec query",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return("", redis.Nil)
			},
			want:    []entity.Product{},
			wantErr: true,
		},
		{
			name: "all ok but failed to set cache",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("product 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return("", redis.Nil)
				mock.redis.EXPECT().SetEX(context.Background(), fmt.Sprintf(getProductList, marshalledParam), string(marshalledResult), time.Minute).Return(assert.AnError)
			},
			want: []entity.Product{
				{
					Name: "product 1",
				},
			},
			wantErr: false,
		},
		{
			name: "all ok",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("product 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return("", redis.Nil)
				mock.redis.EXPECT().SetEX(context.Background(), fmt.Sprintf(getProductList, marshalledParam), string(marshalledResult), time.Minute).Return(nil)
			},
			want: []entity.Product{
				{
					Name: "product 1",
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
			got, err := u.GetList(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_product_GetListByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `products` WHERE `products`.`id` IN (?,?) AND `products`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	mockParam := []uint{1, 2}
	marshalledParam, _ := json.Marshal(mockParam)

	mockRedis := mock_redis.NewMockInterface(ctrl)

	mockResult := []entity.Product{
		{
			Name: "product 1",
		},
	}

	marshalledResult, _ := json.Marshal(mockResult)

	type mockFields struct {
		redis *mock_redis.MockInterface
	}

	mocks := mockFields{
		redis: mockRedis,
	}

	type args struct {
		ctx        context.Context
		productIDs []uint
	}

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		mockFunc    func(mock mockFields)
		want        []entity.Product
		wantErr     bool
	}{
		{
			name: "success get from cache",
			args: args{
				ctx:        context.Background(),
				productIDs: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return(string(marshalledResult), nil)
			},
			want:    mockResult,
			wantErr: false,
		},
		{
			name: "failed to get cache list",
			args: args{
				ctx:        context.Background(),
				productIDs: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return("", assert.AnError)
			},
			want:    []entity.Product{},
			wantErr: true,
		},
		{
			name: "failed to exec query",
			args: args{
				ctx:        context.Background(),
				productIDs: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return("", redis.Nil)
			},
			want:    []entity.Product{},
			wantErr: true,
		},
		{
			name: "all ok but failed to set cache",
			args: args{
				ctx:        context.Background(),
				productIDs: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("product 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return("", redis.Nil)
				mock.redis.EXPECT().SetEX(context.Background(), fmt.Sprintf(getProductList, string(marshalledParam)), string(marshalledResult), time.Minute).Return(assert.AnError)
			},
			want: []entity.Product{
				{
					Name: "product 1",
				},
			},
			wantErr: false,
		},
		{
			name: "all ok",
			args: args{
				ctx:        context.Background(),
				productIDs: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("product 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductList, marshalledParam)).Return("", redis.Nil)
				mock.redis.EXPECT().SetEX(context.Background(), fmt.Sprintf(getProductList, string(marshalledParam)), string(marshalledResult), time.Minute).Return(nil)
			},
			want: []entity.Product{
				{
					Name: "product 1",
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
			got, err := u.GetListByID(tt.args.ctx, tt.args.productIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.GetListByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_product_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `products` WHERE `products`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	querySqlOK := "SELECT * FROM `products` WHERE `products`.`id` = ? AND `products`.`deleted_at` IS NULL ORDER BY `products`.`id` LIMIT 1"
	queryOK := regexp.QuoteMeta(querySqlOK)

	mockParamEmpty := entity.ProductParam{}

	marshalledParamEmpty, _ := json.Marshal(mockParamEmpty)

	mockParam := entity.ProductParam{
		ID: 1,
	}

	marshalledParam, _ := json.Marshal(mockParam)

	mockRedis := mock_redis.NewMockInterface(ctrl)

	resultMock := entity.Product{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "product 1",
	}

	resultMarshalled, _ := json.Marshal(resultMock)

	type mockFields struct {
		redis *mock_redis.MockInterface
	}

	mocks := mockFields{
		redis: mockRedis,
	}

	type args struct {
		ctx   context.Context
		param entity.ProductParam
	}

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		mockFunc    func(mock mockFields)
		want        entity.Product
		wantErr     bool
	}{
		{
			name: "success get from cache",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductByIdKey, marshalledParam)).Return(string(resultMarshalled), nil)
			},
			want:    resultMock,
			wantErr: false,
		},
		{
			name: "failed to get from cache",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductByIdKey, marshalledParam)).Return("", assert.AnError)
			},
			want:    entity.Product{},
			wantErr: true,
		},
		{
			name: "failed to exec query",
			args: args{
				ctx:   context.Background(),
				param: mockParamEmpty,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductByIdKey, marshalledParamEmpty)).Return("", redis.Nil)
			},
			want:    entity.Product{},
			wantErr: true,
		},
		{
			name: "all ok but failed to set cache",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"id", "name"})
				row.AddRow(1, "product 1")
				sqlMock.ExpectQuery(queryOK).WillReturnRows(row)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductByIdKey, marshalledParam)).Return("", redis.Nil)
				mock.redis.EXPECT().SetEX(context.Background(), fmt.Sprintf(getProductByIdKey, marshalledParam), string(resultMarshalled), time.Minute).Return(assert.AnError)
			},
			want:    resultMock,
			wantErr: false,
		},
		{
			name: "all ok",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"id", "name"})
				row.AddRow(1, "product 1")
				sqlMock.ExpectQuery(queryOK).WillReturnRows(row)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields) {
				mock.redis.EXPECT().Get(context.Background(), fmt.Sprintf(getProductByIdKey, marshalledParam)).Return("", redis.Nil)
				mock.redis.EXPECT().SetEX(context.Background(), fmt.Sprintf(getProductByIdKey, marshalledParam), string(resultMarshalled), time.Minute).Return(nil)
			},
			want:    resultMock,
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
			got, err := u.Get(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
