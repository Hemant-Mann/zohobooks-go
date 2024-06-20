package zohobooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// TaxIGST name of tax type
const TaxIGST = "IGST"

// TaxIGST0 name of tax type
const TaxIGST0 = "IGST0"

const InvStatusPushed = "pushed"

const InvStatusCancelled = "cancelled"

const InvStatusYTP = "yet_to_be_pushed"

// TaxIGST18 name of the tax
const TaxIGST18 = "IGST18"

type taxInfo struct {
	TaxName   string  `json:"tax_name"`
	TaxAmount float64 `json:"tax_amount"`
}

type LineItemTaxes struct {
	TaxId   string  `json:"tax_id"`
	TaxName string  `json:"tax_name"`
	TaxAmt  float32 `json:"tax_amount"`
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
	BranchID          string     `json:"branch_id"`
	BranchName        string     `json:"branch_name"`

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

	Country     string      `json:"country"`
	EInvDetails EInvDetails `json:"einvoice_details"`

	BillingAddress  BillingAddress `json:"billing_address"`
	ShippingAddress BillingAddress `json:"shipping_address"`
}

type EInvDetails struct {
	InvRefNo     string `json:"inv_ref_num"`
	AckNo        string `json:"ack_number"`
	Status       string `json:"status"`
	FormatStatus string `json:"formatted_status"`
	AckDate      string `json:"ack_date"`
}

// LineItem struct contains info about the line items of the invoice
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
	HsnOrSac    string  `json:"hsn_or_sac,omitempty"`
	TaxName     string  `json:"tax_name,omitempty"`
	TaxType     string  `json:"tax_type,omitempty"`
	TaxPercent  float64 `json:"tax_percentage,omitempty"`

	LineItemTaxes []LineItemTaxes `json:"line_item_taxes,omitempty"`
}

// InvoiceParams struct represents the information to create a invoice
type InvoiceParams struct {
	CustomerID     string   `json:"customer_id"`
	ContactPersons []string `json:"contact_persons,omitempty"`
	InvoiceNumber  string   `json:"invoice_number,omitempty"`
	ReferenceNo    string   `json:"reference_number,omitempty"`
	PlaceOfSupply  string   `json:"place_of_supply,omitempty"`

	// possible values ---> vat_registered,vat_not_registered,gcc_vat_not_registered,gcc_vat_registered,non_gcc,dz_vat_registered and dz_vat_not_registered.
	TaxTreatment string `json:"tax_treatment,omitempty"`
	GstNO        string `json:"gst_no,omitempty"`        // 15 digit
	GstTreatment string `json:"gst_treatment,omitempty"` // Allowed values are business_gst , business_none , overseas , consumer

	Date              string     `json:"date,omitempty"`
	PaymentTerms      int        `json:"payment_terms,omitempty"`
	PaymentTermsLabel string     `json:"payment_terms_label,omitempty"`
	DueDate           string     `json:"due_date,omitempty"`
	IsInclusiveTax    bool       `json:"is_inclusive_tax"`
	Discount          float64    `json:"discount,omitempty"`
	TaxID             string     `json:"tax_id,omitempty"`
	Reason            string     `json:"reason,omitempty"` // required when updating sent invoice
	LineItems         []LineItem `json:"line_items"`
	Notes             string     `json:"notes,omitempty"`
	Terms             string     `json:"terms,omitempty"`
	BranchID          string     `json:"branch_id,omitempty"`

	Country string `json:"country"`
}

// InvoiceEmailParams struct contains the parameters to be used while sending invoices
type InvoiceEmailParams struct {
	SendFromOrgEmail bool     `json:"send_from_org_email_id"`
	ToMailIDs        []string `json:"to_mail_ids"`
	CCMailIDs        []string `json:"cc_mail_ids,omitempty"`
	Subject          string   `json:"subject,omitempty"`
	Body             string   `json:"body,omitempty"`
}

// Billing Address Params for Invoice
type BAddrInvoiceParams struct {
	Address string `json:"address,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Country string `json:"country,omitempty"`
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

	respData, err := SendResp(resp, err, i)
	if err != nil {
		return i, err
	}
	return &respData.Invoice, err
}

// Update method will try to update a invoice on razorpay
func (i *Invoice) Update(id string, params *InvoiceParams, client *Client) (*Invoice, error) {
	var body, _ = json.Marshal(params)
	resp, err := client.Put(i.Endpoint()+"/"+id, string(body))

	respData, err := SendResp(resp, err, i)
	if err != nil {
		return i, err
	}
	return &respData.Invoice, err
}

// update invoice billing address
func (i *Invoice) UpdateInvBillingAddress(id string, billingAddress *BAddrInvoiceParams, client *Client) (*Invoice, error) {
	url := fmt.Sprintf("%s/%s/address/billing", i.Endpoint(), id)
	var body, _ = json.Marshal(billingAddress)
	resp, err := client.Put(url, string(body))
	respData, err := SendResp(resp, err, i)
	if err != nil {
		return i, err
	}
	return &respData.Invoice, err
}

// push Invoice to IRP portal and returns the Invoice with IRP Ack Num and Ref Num
func (i *Invoice) PushInvoiceToIRP(id string, client *Client) (*Invoice, error) {
	url := fmt.Sprintf("%s/%s/einvoice/push", i.Endpoint(), id)
	headers := map[string]string{}
	var emptyBody []byte

	resp, err := client.makeRequest("POST", url, bytes.NewBuffer(emptyBody), headers)
	respData, err := SendResp(resp, err, i)
	if err != nil {
		return nil, err
	}
	if len(respData.Data.Errors) > 0 {
		return nil, errors.New(strings.ToLower(respData.Data.Errors[0].Message))
	}
	time.Sleep(1 * time.Second)
	pushedInvoice, err := i.FindOne(id, client)
	if err != nil {
		return nil, err
	}
	return pushedInvoice, err
}

// FindOne tries to find the invoice with given id
func (i *Invoice) FindOne(id string, client *Client) (*Invoice, error) {
	resp, err := client.Get(i.Endpoint() + "/" + id)
	respData, err := SendResp(resp, err, i)
	if err != nil {
		return i, err
	}
	return &respData.Invoice, err
}

// Email method will send the invoice to the customer
func (i *Invoice) Email(id string, params *InvoiceEmailParams, client *Client) {
	var body, _ = json.Marshal(params)
	resp, err := client.Post(i.Endpoint()+"/"+id+"/email?send_attachment=true", string(body))

	SendResp(resp, err, i)
}

//This methods will update the invoice status
func (i *Invoice) UpdateInvStatus(id, status string, client *Client) {
	resp, err := client.Post(i.Endpoint()+"/"+id+"/status/"+status, "")
	SendResp(resp, err, i)
}

// DownloadPDF method will download the pdf to the given filepath
func (i *Invoice) DownloadPDF(id, filepath string, client *Client) error {
	// Create the file
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get the data
	resp, err := client.Get(i.Endpoint() + "/pdf?invoice_ids=" + id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
