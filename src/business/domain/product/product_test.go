package product

import (
	"database/sql"
	"go-clean/src/business/entity"
	"regexp"
	"testing"

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

	mockParam := entity.ProductParam{}

	type args struct {
		param entity.ProductParam
	}

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        []entity.Product
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    []entity.Product{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("product 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
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

			u := Init(sqlClient)
			got, err := u.GetList(tt.args.param)
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

	type args struct {
		productIDs []uint
	}

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        []entity.Product
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			args: args{
				productIDs: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    []entity.Product{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				productIDs: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("product 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
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

			u := Init(sqlClient)
			got, err := u.GetListByID(tt.args.productIDs)
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
	mockParam := entity.ProductParam{
		ID: 1,
	}

	type args struct {
		param entity.ProductParam
	}

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.Product
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			args: args{
				param: mockParamEmpty,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    entity.Product{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"id", "name"})
				row.AddRow(1, "product 1")
				sqlMock.ExpectQuery(queryOK).WillReturnRows(row)
				return sqlServer, err
			},
			want: entity.Product{
				Model: gorm.Model{
					ID: 1,
				},
				Name: "product 1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			u := Init(sqlClient)
			got, err := u.Get(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
