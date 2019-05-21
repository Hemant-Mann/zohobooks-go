package zohobooks

import (
	"encoding/json"
)

// TaxIGST name of tax type
const TaxIGST = "IGST"

// TaxIGST18 name of the tax
const TaxIGST18 = "IGST18"

type taxInfo struct {
	TaxName   string  `json:"tax_name"`
	TaxAmount float64 `json:"tax_amount"`
}

// Invoice struct represents the information of the invoice
type Invoice struct {
	ID             string   `json:"invoice_id"`
	CustomerID     string   `json:"customer_id"`
	CustomerName   string   `json:"customer_name"`
	ContactPersons []string `json:"contact_persons"`
	InvoiceNumber  string   `json:"invoice_number"`
	PlaceOfSupply  string   `json:"place_of_supply"`

	// possible values ---> vat_registered,vat_not_registered,gcc_vat_not_registered,gcc_vat_registered,non_gcc,dz_vat_registered and dz_vat_not_registered.
	TaxTreatment string `json:"tax_treatment"`
	GstNO        string `json:"gst_no"`        // 15 digit
	GstTreatment string `json:"gst_treatment"` // Allowed values are business_gst , business_none , overseas , consumer

	Status            string     `json:"status"`
	Date              string     `json:"date"`
	PaymentTerms      int        `json:"payment_terms"`
	PaymentTermsLabel string     `json:"payment_terms_label"`
	DueDate           string     `json:"due_date"`
	CurrencyCode      string     `json:"currency_code"`
	CurrencyID        string     `json:"currency_id"`
	Discount          float64    `json:"discount"`
	TaxID             string     `json:"tax_id"`
	RefNo             string     `json:"reference_number"`
	LineItems         []LineItem `json:"line_items"`
	Notes             string     `json:"notes"`
	Terms             string     `json:"terms"`

	SubTotal float64   `json:"sub_total"`
	TaxTotal float64   `json:"tax_total"`
	Total    float64   `json:"total"`
	Taxes    []taxInfo `json:"taxes"`

	PaymentReminder   bool    `json:"payment_reminder_enabled"`
	PaymentMade       float64 `json:"payment_made"`
	CreditsApplied    float64 `json:"credits_applied"`
	TaxAmountWithheld float64 `json:"tax_amount_withheld"`
	Balance           float64 `json:"balance"`
	WriteOffAmount    float64 `json:"write_off_amount"`
	CreatedTime       string  `json:"created_time"`
	LastModifiedTime  string  `json:"last_modified_time"`
	InvoiceURL        string  `json:"invoice_url"`
}

type LineItem struct {
	ItemID      string  `json:"item_id,omitempty"`
	ProjectID   string  `json:"project_id,omitempty"`
	ProductType string  `json:"product_type,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rate        float64 `json:"rate"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit,omitempty"`
	TaxID       string  `json:"tax_id,omitempty"`
	TaxName     string  `json:"tax_name,omitempty"`
	TaxType     string  `json:"tax_type,omitempty"`
	TaxPercent  float64 `json:"tax_percentage,omitempty"`
}

// InvoiceParams struct represents the information to create a invoice
type InvoiceParams struct {
	CustomerID     string   `json:"customer_id"`
	ContactPersons []string `json:"contact_persons,omitempty"`
	InvoiceNumber  string   `json:"invoice_number,omitempty"`
	PlaceOfSupply  string   `json:"place_of_supply,omitempty"`

	// possible values ---> vat_registered,vat_not_registered,gcc_vat_not_registered,gcc_vat_registered,non_gcc,dz_vat_registered and dz_vat_not_registered.
	TaxTreatment string `json:"tax_treatment,omitempty"`
	GstNO        string `json:"gst_no,omitempty"`        // 15 digit
	GstTreatment string `json:"gst_treatment,omitempty"` // Allowed values are business_gst , business_none , overseas , consumer

	Date              string     `json:"date,omitempty"`
	PaymentTerms      int        `json:"payment_terms,omitempty"`
	PaymentTermsLabel string     `json:"payment_terms_label,omitempty"`
	DueDate           string     `json:"due_date,omitempty"`
	Discount          float64    `json:"discount,omitempty"`
	TaxID             string     `json:"tax_id,omitempty"`
	RefNo             string     `json:"reference_number,omitempty"`
	LineItems         []LineItem `json:"line_items"`
	Notes             string     `json:"notes,omitempty"`
	Terms             string     `json:"terms,omitempty"`
}

type InvoiceEmailParams struct {
	SendFromOrgEmail bool     `json:"send_from_org_email_id"`
	ToMailIDs        []string `json:"to_mail_ids"`
	CCMailIDs        []string `json:"cc_mail_ids,omitempty"`
	Subject          string   `json:"subject,omitempty"`
	Body             string   `json:"body,omitempty"`
}

// New method will create a invoice object and return a pointer to it
func (i *Invoice) New() Resource {
	var obj = &Invoice{}
	return obj
}

// Endpoint method returns the endpoint of the resource
func (i *Invoice) Endpoint() string {
	return "/invoices"
}

// Create method will try to create a invoice on razorpay
func (i *Invoice) Create(params *InvoiceParams, client *Client) (*Invoice, error) {
	var body, _ = json.Marshal(params)
	resp, err := client.Post(i.Endpoint(), string(body))

	respData, err := sendResp(resp, err, i)
	if err != nil {
		return i, err
	}
	return &respData.Invoice, err
}

// Update method will try to update a invoice on razorpay
func (i *Invoice) Update(id string, params *InvoiceParams, client *Client) (*Invoice, error) {
	var body, _ = json.Marshal(params)
	resp, err := client.Put(i.Endpoint()+"/"+id, string(body))

	respData, err := sendResp(resp, err, i)
	if err != nil {
		return i, err
	}
	return &respData.Invoice, err
}

// FindOne tries to find the invoice with given id
func (i *Invoice) FindOne(id string, client *Client) (*Invoice, error) {
	resp, err := client.Get(i.Endpoint() + "/" + id)
	respData, err := sendResp(resp, err, i)
	if err != nil {
		return i, err
	}
	return &respData.Invoice, err
}

func (i *Invoice) Email(id string, params *InvoiceEmailParams, client *Client) {
	var body, _ = json.Marshal(params)
	resp, err := client.Post(i.Endpoint()+"/"+id+"/email?send_attachment=true", string(body))

	sendResp(resp, err, i)
}
