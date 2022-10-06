package zohobooks

import (
	"encoding/json"
)

// TransactionTypeTransfer constant is a bank transaction type of "transfer_fund"
const TransactionTypeTransfer = "transfer_fund"

// BankTransaction struct will contain all the information of bank
type BankTransaction struct {
	ID          string  `json:"transaction_id"`
	FromAccID   string  `json:"from_account_id"`
	FromAccName string  `json:"from_account_name"`
	ToAccID     string  `json:"to_account_id"`
	ToAccName   string  `json:"to_account_name"`
	Type        string  `json:"transaction_type"`
	Amount      float64 `json:"amount"`
	PaymentMode string  `json:"payment_mode"`
	Date        string  `json:"date"`
	RefNO       string  `json:"reference_number"`
	Description string  `json:"description"`
}

// BankTransactionParams struct contains the params used to create a
// bank transaction on zohobooks
type BankTransactionParams struct {
	FromAccID   string  `json:"from_account_id"`
	ToAccID     string  `json:"to_account_id"`
	Type        string  `json:"transaction_type"`
	Amount      float64 `json:"amount"`
	PaymentMode string  `json:"payment_mode"`
	Date        string  `json:"date"`
	RefNO       string  `json:"reference_number"`
	Description string  `json:"description"`
}

// New method will create an object and return a pointer to it
func (bt *BankTransaction) New() Resource {
	var obj = &BankTransaction{}
	return obj
}

// Endpoint method returns the endpoint of the resource
func (bt *BankTransaction) Endpoint() string {
	return "/banktransactions"
}

// Create method will try to create a bank transaction on zohobooks
func (bt *BankTransaction) Create(params *BankTransactionParams, client *Client) (*BankTransaction, error) {
	var body, _ = json.Marshal(params)
	resp, err := client.Post(bt.Endpoint(), string(body))

	respData, err := SendResp(resp, err, bt)
	if err != nil {
		return bt, err
	}
	return &respData.BankTransaction, err
}
