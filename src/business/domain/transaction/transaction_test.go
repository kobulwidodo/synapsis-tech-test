package transaction

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

func Test_transaction_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "INSERT INTO"
	query := regexp.QuoteMeta(querySql)

	mockTransaction := entity.Transaction{
		UserID: 1,
	}

	type args struct {
		transaction entity.Transaction
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.Transaction
		wantErr     bool
	}{
		{
			name: "failed to create transaction",
			args: args{
				transaction: mockTransaction,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    mockTransaction,
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				transaction: mockTransaction,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
				sqlMock.ExpectationsWereMet()
				return sqlServer, err
			},
			want:    mockTransaction,
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
			_, err = u.Create(tt.args.transaction)
			if (err != nil) != tt.wantErr {
				t.Errorf("transaction.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
