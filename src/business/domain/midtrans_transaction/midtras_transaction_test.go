package midtranstransaction

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

func Test_midtransTransaction_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "INSERT INTO"
	query := regexp.QuoteMeta(querySql)

	mockMidtransTransaction := entity.MidtransTransaction{
		MidtransID: "a",
	}

	type args struct {
		midtransTransaction entity.MidtransTransaction
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.MidtransTransaction
		wantErr     bool
	}{
		{
			name: "failed to create midtransTransaction",
			args: args{
				midtransTransaction: mockMidtransTransaction,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    mockMidtransTransaction,
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				midtransTransaction: mockMidtransTransaction,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
				sqlMock.ExpectationsWereMet()
				return sqlServer, err
			},
			want:    mockMidtransTransaction,
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
			_, err = u.Create(tt.args.midtransTransaction)
			if (err != nil) != tt.wantErr {
				t.Errorf("midtransTransaction.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_midtransTransaction_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `midtrans_transactions` WHERE `midtrans_transactions`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	mockParam := entity.MidtransTransactionParam{}

	type args struct {
		param entity.MidtransTransactionParam
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.MidtransTransaction
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
			want:    entity.MidtransTransaction{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"order_id"})
				row.AddRow("cl-1-1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			want: entity.MidtransTransaction{
				OrderID: "cl-1-1",
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
				t.Errorf("midtransTransaction.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_midtransTransaction_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "UPDATE"
	query := regexp.QuoteMeta(querySql)

	selectParam := entity.MidtransTransactionParam{
		ID: 1,
	}

	updateParam := entity.UpdateMidtransTransactionParam{
		Status: "inactive",
	}

	type args struct {
		selectParam entity.MidtransTransactionParam
		updateParam entity.UpdateMidtransTransactionParam
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
				selectParam: selectParam,
				updateParam: updateParam,
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
				selectParam: selectParam,
				updateParam: updateParam,
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

			u := Init(sqlClient)
			err = u.Update(tt.args.selectParam, tt.args.updateParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("midtransTransaction.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
