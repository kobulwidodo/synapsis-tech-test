package transaction_test

import (
	"context"
	"encoding/json"
	mock_cart "go-clean/src/business/domain/mock/cart"
	mock_midtrans "go-clean/src/business/domain/mock/midtrans"
	mock_midtrans_transaction "go-clean/src/business/domain/mock/midtrans_transaction"
	mock_product "go-clean/src/business/domain/mock/product"
	mock_transaction "go-clean/src/business/domain/mock/transaction"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/transaction"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/midtrans"
	mock_auth "go-clean/src/lib/tests/mock/auth"
	"testing"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_transaction_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)
	productMock := mock_product.NewMockInterface(ctrl)
	midtransMock := mock_midtrans.NewMockInterface(ctrl)
	transactionMock := mock_transaction.NewMockInterface(ctrl)
	midtransTransactionMock := mock_midtrans_transaction.NewMockInterface(ctrl)

	tr := transaction.Init(authMock, transactionMock, cartMock, productMock, midtransMock, midtransTransactionMock)

	userAuthMock := auth.UserAuthInfo{
		User: auth.User{
			ID:   1,
			Name: "mail",
		},
	}

	paramsMock := entity.CreateTransactionParam{
		AddressShip: "purwakarta",
		PaymentID:   1,
	}

	paramsMockUndifinedPaymentMock := entity.CreateTransactionParam{
		AddressShip: "purwakarta",
		PaymentID:   999,
	}

	cartParamMock := entity.CartParam{
		UserID: 1,
		Status: entity.StatusInCart,
	}

	cartResultMock := []entity.Cart{
		{
			Model: gorm.Model{
				ID: 1,
			},
			UserID:    1,
			ProductID: 1,
			Qty:       1,
		},
	}

	productResultMock := []entity.Product{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:  "product 1",
			Price: 10000,
		},
	}

	newTransactionMock := entity.Transaction{
		UserID:      1,
		AddressShip: "purwakarta",
		TotalPrice:  10000,
	}

	transactionResultMock := entity.Transaction{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:      1,
		AddressShip: "purwakarta",
		TotalPrice:  10000,
	}

	midtransCreateParamMock := midtrans.CreateOrderParam{
		OrderID:     1,
		PaymentID:   1,
		GrossAmount: 10000,
		ItemsDetails: []midtrans.ItemsDetails{
			{
				ID:    "1",
				Price: 10000,
				Qty:   1,
				Name:  "product 1",
			},
		},
		CustomerDetails: midtrans.CustomerDetails{
			Name: "mail",
		},
	}

	midtransCreateParamUndifinedMock := midtrans.CreateOrderParam{
		OrderID:     1,
		PaymentID:   999,
		GrossAmount: 10000,
		ItemsDetails: []midtrans.ItemsDetails{
			{
				ID:    "1",
				Price: 10000,
				Qty:   1,
				Name:  "product 1",
			},
		},
		CustomerDetails: midtrans.CustomerDetails{
			Name: "mail",
		},
	}

	midtransResultMock := &coreapi.ChargeResponse{
		TransactionID: "1",
		OrderID:       "1",
		Actions: []coreapi.Action{
			{
				URL: "url 1",
			},
			{
				URL: "url 2",
			},
		},
	}

	paymentData, _ := json.Marshal(entity.PaymentData{
		Key: "url 2",
		Qr:  "url 1",
	})

	newMidtransTransactionMock := entity.MidtransTransaction{
		TransactionID: 1,
		MidtransID:    "1",
		OrderID:       "1",
		PaymentType:   1,
		Status:        entity.StatusPending,
		PaymentData:   string(paymentData),
	}

	selectParamCartMock := entity.CartParam{
		Status: entity.StatusInCart,
		UserID: 1,
	}

	updateParamCartMock := entity.UpdateCartParam{
		Status:        entity.StatusUnpaid,
		TransactionID: 1,
	}

	selectParamCartFinalPrice := entity.CartParam{
		ID: 1,
	}

	updateParamCartFinalPrice := entity.UpdateCartParam{
		FinalPricePerItem: 10000,
	}

	type mockfields struct {
		auth                 *mock_auth.MockInterface
		cart                 *mock_cart.MockInterface
		product              *mock_product.MockInterface
		midtrans             *mock_midtrans.MockInterface
		transaction          *mock_transaction.MockInterface
		midtrans_transaction *mock_midtrans_transaction.MockInterface
	}

	mocks := mockfields{
		auth:                 authMock,
		cart:                 cartMock,
		product:              productMock,
		midtrans:             midtransMock,
		transaction:          transactionMock,
		midtrans_transaction: midtransTransactionMock,
	}

	type args struct {
		ctx   context.Context
		param entity.CreateTransactionParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     entity.Transaction
		wantErr  bool
	}{
		{
			name: "failed to get auth user",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    entity.Transaction{},
			wantErr: true,
		},
		{
			name: "failed get cart list",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    entity.Transaction{},
			wantErr: true,
		},
		{
			name: "cart empty",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, nil)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    entity.Transaction{},
			wantErr: true,
		},
		{
			name: "failed to get products list",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.product.EXPECT().GetListByID(context.Background(), []uint{1}).Return([]entity.Product{}, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    entity.Transaction{},
			wantErr: true,
		},
		{
			name: "failed to create transaction",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.product.EXPECT().GetListByID(context.Background(), []uint{1}).Return(productResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    transactionResultMock,
			wantErr: true,
		},
		{
			name: "failed to create midtrans",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.product.EXPECT().GetListByID(context.Background(), []uint{1}).Return(productResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(nil, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    transactionResultMock,
			wantErr: true,
		},
		{
			name: "failed to update cart",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.product.EXPECT().GetListByID(context.Background(), []uint{1}).Return(productResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(midtransResultMock, nil)
				mock.cart.EXPECT().Update(selectParamCartMock, updateParamCartMock).Return(assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    transactionResultMock,
			wantErr: true,
		},
		{
			name: "failed to get payment data",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.product.EXPECT().GetListByID(context.Background(), []uint{1}).Return(productResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamUndifinedMock).Return(midtransResultMock, nil)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMockUndifinedPaymentMock,
			},
			want:    transactionResultMock,
			wantErr: true,
		},
		{
			name: "failed to create midtrans transaction",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.product.EXPECT().GetListByID(context.Background(), []uint{1}).Return(productResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(midtransResultMock, nil)
				mock.cart.EXPECT().Update(selectParamCartMock, updateParamCartMock).Return(nil)
				mock.midtrans_transaction.EXPECT().Create(newMidtransTransactionMock).Return(entity.MidtransTransaction{}, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    transactionResultMock,
			wantErr: true,
		},
		{
			name: "failed to update cart to set final price",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.product.EXPECT().GetListByID(context.Background(), []uint{1}).Return(productResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(midtransResultMock, nil)
				mock.cart.EXPECT().Update(selectParamCartMock, updateParamCartMock).Return(nil)
				mock.midtrans_transaction.EXPECT().Create(newMidtransTransactionMock).Return(entity.MidtransTransaction{}, nil)
				mock.cart.EXPECT().Update(selectParamCartFinalPrice, updateParamCartFinalPrice).Return(assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    transactionResultMock,
			wantErr: false,
		},
		{
			name: "all success",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.product.EXPECT().GetListByID(context.Background(), []uint{1}).Return(productResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(midtransResultMock, nil)
				mock.cart.EXPECT().Update(selectParamCartMock, updateParamCartMock).Return(nil)
				mock.midtrans_transaction.EXPECT().Create(newMidtransTransactionMock).Return(entity.MidtransTransaction{}, nil)
				mock.cart.EXPECT().Update(selectParamCartFinalPrice, updateParamCartFinalPrice).Return(nil)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    transactionResultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := tr.Create(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("transaction.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_transaction_ValidateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transactionMock := mock_transaction.NewMockInterface(ctrl)

	tr := transaction.Init(nil, transactionMock, nil, nil, nil, nil)

	authUserMock := auth.UserAuthInfo{
		User: auth.User{
			ID: 1,
		},
	}

	authUserFailedMock := auth.UserAuthInfo{
		User: auth.User{
			ID: 2,
		},
	}

	transactionParamMock := entity.TransactionParam{
		ID: 1,
	}

	transactionResultMock := entity.Transaction{
		UserID: 1,
	}

	type mockFields struct {
		transaction *mock_transaction.MockInterface
	}

	mocks := mockFields{
		transaction: transactionMock,
	}

	type args struct {
		ctx           context.Context
		transactionID uint
		user          auth.UserAuthInfo
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		wantErr  bool
	}{
		{
			name: "failed transaction id = 0",
			args: args{
				ctx:           context.Background(),
				transactionID: 0,
			},
			mockFunc: func(mock mockFields, arg args) {},
			wantErr:  true,
		},
		{
			name: "failed to get cart",
			args: args{
				ctx:           context.Background(),
				transactionID: 1,
				user:          authUserMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(entity.Transaction{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed unauthorized",
			args: args{
				ctx:           context.Background(),
				transactionID: 1,
				user:          authUserFailedMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(transactionResultMock, nil)
			},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx:           context.Background(),
				transactionID: 1,
				user:          authUserMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(transactionResultMock, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := tr.ValidateTransaction(tt.args.ctx, tt.args.transactionID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("menu.ValidateMenu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
