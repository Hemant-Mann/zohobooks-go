package zohobooks

import (
	"encoding/json"
	"errors"
)

type InvoiceInfo struct {
	InvoiceID         string  `json:"invoice_id"`
	Number            string  `json:"invoice_number,omitempty"`
	Date              string  `json:"date,omitempty"`
	Amount            float64 `json:"invoice_amount,omitempty"`
	AmountApplied     float64 `json:"amount_applied"`
	BalanceAmount     float64 `json:"balance_amount,omitempty"`
	TaxAmountWithheld float64 `json:"tax_amount_withheld,omitempty"`
}

type Payment struct {
	ID             string        `json:"payment_id"`
	Mode           string        `json:"payment_mode"`
	Amount         float64       `json:"amount"`
	AmountRefunded float64       `json:"amount_refunded"`
	BankCharges    float64       `json:"bank_charges"`
	Date           string        `json:"date"`
	Status         string        `json:"status"`
	RefNo          string        `json:"reference_number"`
	Description    string        `json:"description"`
	CustomerID     string        `json:"customer_id"`
	CustomerName   string        `json:"customer_name"`
	Email          string        `json:"email"`
	Invoices       []InvoiceInfo `json:"invoices"`
	CurrencyCode   string        `json:"currency_code"`
	CurrencySymbol string        `json:"currency_symbol"`
}

type PaymentFindOptions struct {
	CustomerID string
}

type PaymentParams struct {
	CustomerID  string  `json:"customer_id"`
	Mode        string  `json:"payment_mode"` // This can be check, cash, creditcard, banktransfer, bankremittance, autotransaction or others
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	RefNo       string  `json:"reference_number,omitempty"`
	Description string  `json:"description,omitempty"`

	Invoices       []InvoiceInfo `json:"invoices"`
	BankCharges    float64       `json:"bank_charges"`
	AccountID      string        `json:"account_id,omitempty"`
	ContactPersons []string      `json:"contact_persons,omitempty"`
}

// New method will create a payment object and return a pointer to it
func (p *Payment) New() Resource {
	var obj = &Payment{}
	return obj
}

// Endpoint method returns the endpoint of the resource
func (p *Payment) Endpoint() string {
	return "/customerpayments"
}

// Create method will try to create a contact on razorpay
func (p *Payment) Create(params *PaymentParams, client *Client) (*Payment, error) {
	var body, _ = json.Marshal(params)
	resp, err := client.Post(p.Endpoint(), string(body))

	respData, err := sendResp(resp, err, p)
	if err != nil {
		return p, err
	}
	return &respData.Payment, err
}

// FindOne tries to find the contact with given id
func (p *Payment) FindOne(id string, client *Client) (*Payment, error) {
	resp, err := client.Get(p.Endpoint() + "/" + id)
	respData, err := sendResp(resp, err, p)
	if err != nil {
		return p, err
	}
	return &respData.Payment, err
}

// Delete tries to delete the payment with given id
func (p *Payment) Delete(id string, client *Client) error {
	resp, err := client.Delete(p.Endpoint() + "/" + id)
	respData, err := sendResp(resp, err, p)
	if err != nil {
		return err
	}
	if respData.Code == 0 {
		return nil
	}
	return errors.New(respData.Message)
}

// FindAll tries to find the payment with given options
func (p *Payment) FindAll(opts *PaymentFindOptions, client *Client) ([]Payment, error) {
	resp, err := client.Get(p.Endpoint())
	respData, err := sendResp(resp, err, p)

	var results []Payment
	if err != nil {
		return results, err
	}
	for _, pmt := range respData.Payments {
		results = append(results, pmt)
	}
	return results, err
}
