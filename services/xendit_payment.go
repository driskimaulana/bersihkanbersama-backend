package services

import (
	"github.com/xendit/xendit-go/invoice"
)

func CreateInvoice(amount *float64, payerEmail *string, userId *string) (string, error) {
	// create new invoice
	inv := invoice.CreateParams{
		ExternalID:  *userId,
		Amount:      *amount,
		PayerEmail:  *payerEmail,
		Description: "Donation Invoice",
		Currency:    "IDR",
	}

	resp, err := invoice.Create(&inv)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func GetPaymentDetails(paymentId *string) (string, error) {

	getParams := invoice.GetParams{
		ID: *paymentId,
	}

	resp, err := invoice.Get(&getParams)
	if err != nil {
		return "", err
	}

	return resp.InvoiceURL, nil

}
