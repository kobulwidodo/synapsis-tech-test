package midtrans

import (
	midtransSdk "github.com/midtrans/midtrans-go"
)

const (
	GopayPayment = 1
)

type CreateOrderParam struct {
	PaymentID       int
	OrderID         uint
	GrossAmount     int64
	ItemsDetails    []ItemsDetails
	CustomerDetails CustomerDetails
}

type ItemsDetails struct {
	ID    string
	Price int64
	Qty   int
	Name  string
}

type CustomerDetails struct {
	Name  string
	Email string
}

func (cop *CreateOrderParam) convertToItemDetails() *[]midtransSdk.ItemDetails {
	itemsDetails := []midtransSdk.ItemDetails{}
	for _, i := range cop.ItemsDetails {
		itemDetail := midtransSdk.ItemDetails{
			ID:    i.ID,
			Price: i.Price,
			Qty:   int32(i.Qty),
			Name:  i.Name,
		}
		itemsDetails = append(itemsDetails, itemDetail)
	}

	return &itemsDetails
}
