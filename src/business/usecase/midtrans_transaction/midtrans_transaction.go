package midtranstransaction

import (
	"encoding/json"
	"errors"
	cartDom "go-clean/src/business/domain/cart"
	midtransDom "go-clean/src/business/domain/midtrans"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	"go-clean/src/business/entity"
)

type Interface interface {
	GetPaymentDetail(param entity.MidtransTransactionParam) (entity.MidtransTransactionPaymentDetail, error)
	HandleNotification(payload map[string]interface{}) error
}

type midtransTransaction struct {
	midtrans            midtransDom.Interface
	midtransTransaction midtransTransactionDom.Interface
	cart                cartDom.Interface
}

func Init(mttd midtransTransactionDom.Interface, md midtransDom.Interface, cd cartDom.Interface) Interface {
	mtt := &midtransTransaction{
		midtrans:            md,
		midtransTransaction: mttd,
		cart:                cd,
	}

	return mtt
}

func (mtt *midtransTransaction) GetPaymentDetail(param entity.MidtransTransactionParam) (entity.MidtransTransactionPaymentDetail, error) {
	result := entity.MidtransTransactionPaymentDetail{}

	midtransTransaction, err := mtt.midtransTransaction.Get(param)
	if err != nil {
		return result, err
	}

	paymentData := entity.PaymentData{}
	if err := json.Unmarshal([]byte(midtransTransaction.PaymentData), &paymentData); err != nil {
		return result, err
	}

	result.Status = midtransTransaction.Status
	result.PaymentData = paymentData
	result.MidtransID = midtransTransaction.OrderID

	return result, nil
}

func (mtt *midtransTransaction) HandleNotification(payload map[string]interface{}) error {
	orderId, exist := payload["order_id"].(string)
	if !exist {
		return errors.New("order id not exist")
	}

	transactionResponse, err := mtt.midtrans.HandleNotification(orderId)
	if err != nil {
		return err
	}

	midtransTransaction, err := mtt.midtransTransaction.Get(entity.MidtransTransactionParam{
		OrderID: orderId,
	})
	if err != nil {
		return err
	}

	status := ""

	if transactionResponse != nil {
		// 5. Do set transaction status based on response from check transaction status
		if transactionResponse.TransactionStatus == "capture" {
			if transactionResponse.FraudStatus == "challenge" {
				// TODO set transaction status on your database to 'challenge'
				status = entity.StatusChallange
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			} else if transactionResponse.FraudStatus == "accept" {
				// TODO set transaction status on your database to 'success'
				status = entity.StatusSuccess
			}
		} else if transactionResponse.TransactionStatus == "settlement" {
			// TODO set transaction status on your databaase to 'success'
			status = entity.StatusSuccess
		} else if transactionResponse.TransactionStatus == "deny" {
			// TODO you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
			status = entity.StatusDeny
		} else if transactionResponse.TransactionStatus == "cancel" || transactionResponse.TransactionStatus == "expire" {
			// TODO set transaction status on your databaase to 'failure'
			status = entity.StatusFailure
		} else if transactionResponse.TransactionStatus == "pending" {
			// TODO set transaction status on your databaase to 'pending' / waiting payment
			status = entity.StatusPending
		}
	}

	if err := mtt.midtransTransaction.Update(entity.MidtransTransactionParam{
		ID: midtransTransaction.ID,
	}, entity.UpdateMidtransTransactionParam{
		Status: status,
	}); err != nil {
		return err
	}

	if status == entity.StatusSuccess {
		if err := mtt.cart.Update(entity.CartParam{
			Status:        entity.StatusUnpaid,
			TransactionID: midtransTransaction.TransactionID,
		}, entity.UpdateCartParam{
			Status: entity.StatusPaid,
		}); err != nil {
			return err
		}
	}

	return nil
}
