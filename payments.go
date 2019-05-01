package zohobooks

type invoiceInfo struct {
	ID            string  `json:"invoice_id"`
	Number        string  `json:"invoice_number"`
	Date          string  `json:"date"`
	Amount        float64 `json:"invoice_amount"`
	AmountApplied float64 `json:"amount_applied"`
	BalanceAmount float64 `json:"balance_amount"`
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
	Invoices       []invoiceInfo `json:"invoices"`
	CurrencyCode   string        `json:"currency_code"`
	CurrencySymbol string        `json:"currency_symbol"`
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
