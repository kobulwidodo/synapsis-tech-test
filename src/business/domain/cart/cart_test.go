package cart

import (
	"database/sql"
	"database/sql/driver"
	"go-clean/src/business/entity"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_cart_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "INSERT INTO"
	query := regexp.QuoteMeta(querySql)

	mockCart := entity.Cart{
		UserID:    1,
		ProductID: 1,
		Qty:       1,
	}

	type args struct {
		cart entity.Cart
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.Cart
		wantErr     bool
	}{
		{
			name: "failed to create cart",
			args: args{
				cart: mockCart,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    mockCart,
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				cart: mockCart,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
				sqlMock.ExpectationsWereMet()
				return sqlServer, err
			},
			want:    mockCart,
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

			c := Init(sqlClient)
			_, err = c.Create(tt.args.cart)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_cart_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `carts` WHERE `carts`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	mockParam := entity.CartParam{}

	type args struct {
		param entity.CartParam
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        []entity.Cart
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
			want:    []entity.Cart{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"user_id"})
				row.AddRow(1)
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			want: []entity.Cart{
				{
					UserID: 1,
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

			c := Init(sqlClient)
			got, err := c.GetList(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cart_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `carts` WHERE `carts`.`user_id` = ? AND `carts`.`deleted_at` IS NULL ORDER BY `carts`.`id` LIMIT 1"
	query := regexp.QuoteMeta(querySql)

	mockParam := entity.CartParam{
		UserID: 1,
	}

	type args struct {
		param entity.CartParam
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.Cart
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
			want:    entity.Cart{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"user_id"})
				row.AddRow(1)
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			want: entity.Cart{
				UserID: 1,
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

			c := Init(sqlClient)
			got, err := c.Get(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cart_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "UPDATE"
	query := regexp.QuoteMeta(querySql)

	mockSelectParam := entity.CartParam{
		ID: 1,
	}
	mockUpdateParam := entity.UpdateCartParam{
		Qty: 2,
	}

	type args struct {
		selectParam entity.CartParam
		updateParam entity.UpdateCartParam
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			args: args{
				selectParam: mockSelectParam,
				updateParam: mockUpdateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				selectParam: mockSelectParam,
				updateParam: mockUpdateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(1))
				sqlMock.ExpectCommit()
				sqlMock.ExpectationsWereMet()
				return sqlServer, err
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

			c := Init(sqlClient)
			err = c.Update(tt.args.selectParam, tt.args.updateParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
