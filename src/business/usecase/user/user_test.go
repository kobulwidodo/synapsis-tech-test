package user_test

import (
	"errors"
	mock_user "go-clean/src/business/domain/mock/user"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/user"
	mock_auth "go-clean/src/lib/tests/mock/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Test_user_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mock_user.NewMockInterface(ctrl)
	hashPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)

	mockParams := entity.CreateUserParam{
		Username: "mail",
		Password: "password",
	}

	mockUserResult := entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "mail",
		Password: string(hashPass),
	}

	u := user.Init(userMock, nil)

	type mockfields struct {
		user *mock_user.MockInterface
	}

	mocks := mockfields{
		user: userMock,
	}

	type args struct {
		params entity.CreateUserParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     entity.User
		wantErr  bool
	}{
		{
			name: "failed to create user",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Create(gomock.Any()).Return(mockUserResult, assert.AnError)
			},
			args: args{
				params: mockParams,
			},
			want: entity.User{
				Username: "mail",
			},
			wantErr: true,
		},
		{
			name: "all ok",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Create(gomock.Any()).Return(mockUserResult, nil)
			},
			args: args{
				params: mockParams,
			},
			want: entity.User{
				Username: "mail",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Create(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Username, got.Username)
		})
	}
}

func Test_user_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mock_user.NewMockInterface(ctrl)

	userParamMock := entity.UserParam{
		ID: 1,
	}

	userOkResult := entity.User{
		Username: "mail",
	}

	u := user.Init(userMock, nil)

	type mockFields struct {
		product *mock_user.MockInterface
	}
	mocks := mockFields{
		product: userMock,
	}

	type args struct {
		id uint
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     entity.User
		wantErr  bool
	}{
		{
			name: "failed to get products",
			args: args{
				id: 1,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.product.EXPECT().Get(userParamMock).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				id: 1,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.product.EXPECT().Get(userParamMock).Return(userOkResult, nil)
			},
			want:    userOkResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.GetById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_user_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mock_user.NewMockInterface(ctrl)
	authMock := mock_auth.NewMockInterface(ctrl)

	mockParams := entity.LoginUserParam{
		Username: "mail",
		Password: "password",
	}

	mockGetUserParam := entity.UserParam{
		Username: "mail",
	}

	hashPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)

	mockUserResult := entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "mail",
		Password: string(hashPass),
	}

	mockToken := "mockToken"

	u := user.Init(userMock, authMock)

	type mockfields struct {
		user *mock_user.MockInterface
		auth *mock_auth.MockInterface
	}

	mocks := mockfields{
		user: userMock,
		auth: authMock,
	}

	type args struct {
		params entity.LoginUserParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "failed to find user",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(mockGetUserParam).Return(entity.User{}, errors.New("user not found"))
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "user not found",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(mockGetUserParam).Return(entity.User{}, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "password incorrect",
			mockFunc: func(mock mockfields, arg args) {
				mockUserResultWithWrongPassword := mockUserResult
				mockUserResultWithWrongPassword.Password = "wrongPassword"
				mock.user.EXPECT().Get(mockGetUserParam).Return(mockUserResultWithWrongPassword, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "failed to generate token",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(mockGetUserParam).Return(mockUserResult, nil)
				mock.auth.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("failed to generate token"))
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(mockGetUserParam).Return(mockUserResult, nil)
				mock.auth.EXPECT().GenerateToken(gomock.Any()).Return(mockToken, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    mockToken,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Login(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
