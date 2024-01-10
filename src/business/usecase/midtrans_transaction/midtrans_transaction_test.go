package midtranstransaction_test

import (
	"encoding/json"
	mock_cart "go-clean/src/business/domain/mock/cart"
	mock_midtrans "go-clean/src/business/domain/mock/midtrans"
	mock_midtranstransaction "go-clean/src/business/domain/mock/midtrans_transaction"
	"go-clean/src/business/entity"
	"testing"

	midtranstransaction "go-clean/src/business/usecase/midtrans_transaction"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_midtransTransaction_GetPaymentDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	midtransTransactionMock := mock_midtranstransaction.NewMockInterface(ctrl)

	midtransTransactionParamMock := entity.MidtransTransactionParam{
		ID: 1,
	}

	paymentDataMock := entity.PaymentData{
		Key: "key",
		Qr:  "qr",
	}

	paymentDataMarshalledMock, _ := json.Marshal(paymentDataMock)

	midtransTransactionResultMock := entity.MidtransTransaction{
		PaymentData: string(paymentDataMarshalledMock),
		OrderID:     "1",
		Status:      entity.StatusSuccess,
	}

	resultMock := entity.MidtransTransactionPaymentDetail{
		Status:      entity.StatusSuccess,
		PaymentData: paymentDataMock,
		MidtransID:  "1",
	}

	mt := midtranstransaction.Init(midtransTransactionMock, nil, nil)

	type mockFields struct {
		midtrans_transaction *mock_midtranstransaction.MockInterface
	}

	mocks := mockFields{
		midtrans_transaction: midtransTransactionMock,
	}

	type args struct {
		param entity.MidtransTransactionParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     entity.MidtransTransactionPaymentDetail
		wantErr  bool
	}{
		{
			name: "failed to get midtrans transaction",
			args: args{
				param: midtransTransactionParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(entity.MidtransTransaction{}, assert.AnError)
			},
			want:    entity.MidtransTransactionPaymentDetail{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				param: midtransTransactionParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
			},
			want:    resultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := mt.GetPaymentDetail(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("midtransTransaction.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_midtransTransaction_HandleNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	midtransMock := mock_midtrans.NewMockInterface(ctrl)
	midtransTransactionMock := mock_midtranstransaction.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)

	payloadMock := map[string]interface{}{
		"order_id": "1",
	}

	transactionResponseMock := &coreapi.TransactionStatusResponse{
		TransactionStatus: "capture",
		FraudStatus:       "accept",
	}

	transactionResponseChallengeMock := &coreapi.TransactionStatusResponse{
		TransactionStatus: "capture",
		FraudStatus:       "challenge",
	}

	transactionResponseSettlementMock := &coreapi.TransactionStatusResponse{
		TransactionStatus: "settlement",
	}

	transactionResponseDenyMock := &coreapi.TransactionStatusResponse{
		TransactionStatus: "deny",
	}

	transactionResponseCancelMock := &coreapi.TransactionStatusResponse{
		TransactionStatus: "cancel",
	}

	transactionResponsePendingMock := &coreapi.TransactionStatusResponse{
		TransactionStatus: "pending",
	}

	midtransTransactionParamMock := entity.MidtransTransactionParam{
		OrderID: "1",
	}

	midtransTransactionResultMock := entity.MidtransTransaction{
		Model: gorm.Model{
			ID: 1,
		},
		TransactionID: 1,
	}

	midtransTransactionUpdateParamMock := entity.MidtransTransactionParam{
		ID: 1,
	}

	midtransTransactionUpdateMock := entity.UpdateMidtransTransactionParam{
		Status: entity.StatusSuccess,
	}

	midtransTransactionUpdateChallangeMock := entity.UpdateMidtransTransactionParam{
		Status: entity.StatusChallange,
	}

	midtransTransactionUpdateDenyMock := entity.UpdateMidtransTransactionParam{
		Status: entity.StatusDeny,
	}

	midtransTransactionUpdateFailureMock := entity.UpdateMidtransTransactionParam{
		Status: entity.StatusFailure,
	}

	midtransTransactionUpdatePendingMock := entity.UpdateMidtransTransactionParam{
		Status: entity.StatusPending,
	}

	cartUpdateParamMock := entity.CartParam{
		Status:        entity.StatusUnpaid,
		TransactionID: 1,
	}

	cartUpdateMock := entity.UpdateCartParam{
		Status: entity.StatusPaid,
	}

	mt := midtranstransaction.Init(midtransTransactionMock, midtransMock, cartMock)

	type mockFields struct {
		midtrans             *mock_midtrans.MockInterface
		midtrans_transaction *mock_midtranstransaction.MockInterface
		cart                 *mock_cart.MockInterface
	}

	mocks := mockFields{
		midtrans:             midtransMock,
		midtrans_transaction: midtransTransactionMock,
		cart:                 cartMock,
	}

	type args struct {
		payload map[string]interface{}
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		wantErr  bool
	}{
		{
			name: "failed to get order id",
			args: args{
				payload: map[string]interface{}{},
			},
			mockFunc: func(mock mockFields, arg args) {
			},
			wantErr: true,
		},
		{
			name: "failed to handle midtrans",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed get midtrans transaction",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponseMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(entity.MidtransTransaction{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed get update midtrans transaction",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponseMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Update(midtransTransactionUpdateParamMock, midtransTransactionUpdateMock).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed update cart",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponseMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Update(midtransTransactionUpdateParamMock, midtransTransactionUpdateMock).Return(nil)
				mock.cart.EXPECT().Update(cartUpdateParamMock, cartUpdateMock).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponseMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Update(midtransTransactionUpdateParamMock, midtransTransactionUpdateMock).Return(nil)
				mock.cart.EXPECT().Update(cartUpdateParamMock, cartUpdateMock).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "all success settlement",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponseSettlementMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Update(midtransTransactionUpdateParamMock, midtransTransactionUpdateMock).Return(nil)
				mock.cart.EXPECT().Update(cartUpdateParamMock, cartUpdateMock).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "all success challenge",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponseChallengeMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Update(midtransTransactionUpdateParamMock, midtransTransactionUpdateChallangeMock).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "all success deny",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponseDenyMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Update(midtransTransactionUpdateParamMock, midtransTransactionUpdateDenyMock).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "all success cancel",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponseCancelMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Update(midtransTransactionUpdateParamMock, midtransTransactionUpdateFailureMock).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "all success pending",
			args: args{
				payload: payloadMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.midtrans.EXPECT().HandleNotification("1").Return(transactionResponsePendingMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Update(midtransTransactionUpdateParamMock, midtransTransactionUpdatePendingMock).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := mt.HandleNotification(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("midtransTransaction.HandleNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
