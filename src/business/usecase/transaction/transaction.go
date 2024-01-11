package transaction

import (
	"context"
	"encoding/json"
	"errors"
	cartDom "go-clean/src/business/domain/cart"
	midtransDom "go-clean/src/business/domain/midtrans"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	productDom "go-clean/src/business/domain/product"
	transactionDom "go-clean/src/business/domain/transaction"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/midtrans"
	"log"
	"strconv"

	"github.com/midtrans/midtrans-go/coreapi"
)

type Interface interface {
	Create(ctx context.Context, createParam entity.CreateTransactionParam) (entity.Transaction, error)
	ValidateTransaction(ctx context.Context, transactionID uint, user auth.UserAuthInfo) error
}

type transaction struct {
	auth                auth.Interface
	cart                cartDom.Interface
	product             productDom.Interface
	transaction         transactionDom.Interface
	midtrans            midtransDom.Interface
	midtransTransaction midtransTransactionDom.Interface
}

func Init(auth auth.Interface, td transactionDom.Interface, cd cartDom.Interface, pd productDom.Interface, md midtransDom.Interface, mtd midtransTransactionDom.Interface) Interface {
	t := &transaction{
		auth:                auth,
		cart:                cd,
		product:             pd,
		transaction:         td,
		midtrans:            md,
		midtransTransaction: mtd,
	}

	return t
}

func (t *transaction) Create(ctx context.Context, createParam entity.CreateTransactionParam) (entity.Transaction, error) {
	user, err := t.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.Transaction{}, err
	}

	carts, err := t.cart.GetList(entity.CartParam{
		UserID: user.User.ID,
		Status: entity.StatusInCart,
	})
	if err != nil {
		return entity.Transaction{}, err
	}

	if len(carts) == 0 {
		return entity.Transaction{}, errors.New("cart is empty")
	}

	productIDs := []uint{}
	for _, c := range carts {
		productIDs = append(productIDs, c.ProductID)
	}

	products, err := t.product.GetListByID(ctx, productIDs)
	if err != nil {
		return entity.Transaction{}, err
	}

	productMap := make(map[uint]entity.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	var totalPrice int64 = 0
	for _, c := range carts {
		totalPrice += int64(c.Qty * productMap[c.ProductID].Price)
	}

	transaction, err := t.transaction.Create(entity.Transaction{
		UserID:      user.User.ID,
		AddressShip: createParam.AddressShip,
		TotalPrice:  totalPrice,
	})
	if err != nil {
		return transaction, err
	}

	coreApiRes, err := t.midtrans.Create(midtrans.CreateOrderParam{
		OrderID:      transaction.ID,
		PaymentID:    createParam.PaymentID,
		GrossAmount:  totalPrice,
		ItemsDetails: t.convertToItemsDetails(carts, productMap),
		CustomerDetails: midtrans.CustomerDetails{
			Name: user.User.Name,
		},
	})
	if err != nil {
		return transaction, err
	}

	paymentData, err := t.getPaymentData(createParam.PaymentID, coreApiRes)
	if err != nil {
		return transaction, err
	}

	paymenDataMarshal, err := json.Marshal(paymentData)
	if err != nil {
		return transaction, err
	}

	if err := t.cart.Update(entity.CartParam{
		Status: entity.StatusInCart,
		UserID: user.User.ID,
	}, entity.UpdateCartParam{
		Status:        entity.StatusUnpaid,
		TransactionID: transaction.ID,
	}); err != nil {
		return transaction, err
	}

	_, err = t.midtransTransaction.Create(entity.MidtransTransaction{
		TransactionID: transaction.ID,
		MidtransID:    coreApiRes.TransactionID,
		OrderID:       coreApiRes.OrderID,
		PaymentType:   createParam.PaymentID,
		Status:        entity.StatusPending,
		PaymentData:   string(paymenDataMarshal),
	})
	if err != nil {
		return transaction, err
	}

	for _, c := range carts {
		if err := t.cart.Update(entity.CartParam{
			ID: c.ID,
		}, entity.UpdateCartParam{
			FinalPricePerItem: productMap[c.ProductID].Price,
		}); err != nil {
			log.Printf("failed to set final price id %d", c.ID)
		}
	}

	return transaction, nil
}

func (t *transaction) getPaymentData(paymentId int, coreApiRes *coreapi.ChargeResponse) (entity.PaymentData, error) {
	paymentData := entity.PaymentData{}
	if paymentId == midtrans.GopayPayment {
		paymentData.Key = coreApiRes.Actions[1].URL
		paymentData.Qr = coreApiRes.Actions[0].URL
	} else {
		return paymentData, errors.New("failed to get payment data")
	}

	return paymentData, nil
}

func (t *transaction) convertToItemsDetails(carts []entity.Cart, products map[uint]entity.Product) []midtrans.ItemsDetails {
	res := []midtrans.ItemsDetails{}
	for _, c := range carts {
		resTemp := midtrans.ItemsDetails{
			ID:    strconv.Itoa(int(c.ID)),
			Price: int64(products[c.ProductID].Price),
			Qty:   c.Qty,
			Name:  products[c.ProductID].Name,
		}
		res = append(res, resTemp)
	}

	return res
}

func (t *transaction) ValidateTransaction(ctx context.Context, transactionID uint, user auth.UserAuthInfo) error {
	if transactionID == 0 {
		return errors.New("please provide transaction id")
	}

	transaction, err := t.transaction.Get(entity.TransactionParam{
		ID: transactionID,
	})
	if err != nil {
		return err
	}

	if transaction.UserID != user.User.ID {
		return errors.New("unauthorized")
	}

	return nil
}
