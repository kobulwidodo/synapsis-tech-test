package category

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

func Test_category_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `categories` WHERE `categories`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	tests := []struct {
		name        string
		prepSqlMock func() (*sql.DB, error)
		want        []entity.Category
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    []entity.Category{},
			wantErr: true,
		},
		{
			name: "all ok",
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("category 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
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
			got, err := u.GetList()
			if (err != nil) != tt.wantErr {
				t.Errorf("category.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
