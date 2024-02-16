package zohobooks

import (
	"encoding/json"
	"errors"
)

type ContactPerson struct {
	ID        string `json:"contact_person_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Skype     string `json:"skype,omitempty"`
	IsPrimary bool   `json:"is_primary"`
}

// Contact struct represents the information of the contact
type Contact struct {
	ID           string `json:"contact_id"`
	Name         string `json:"contact_name"`
	Company      string `json:"company_name"`
	Website      string `json:"website"`
	LanguageCode string `json:"language_code"`
	ContactType  string `json:"contact_type"`
	Notes        string `json:"notes"`
	CurrencyID   string `json:"currency_id"`

	ContactPersons  []ContactPerson `json:"contact_persons"`
	BillingAddress  BillingAddress  `json:"billing_address"`
	ShippingAddress BillingAddress  `json:"shipping_address"`

	// possible values ---> vat_registered,vat_not_registered,gcc_vat_not_registered,gcc_vat_registered,non_gcc,dz_vat_registered and dz_vat_not_registered.
	TaxTreatment string `json:"tax_treatment"`
	GstNO        string `json:"gst_no"`        // 15 digit
	GstTreatment string `json:"gst_treatment"` // Allowed values are business_gst , business_none , overseas , consumer
	TaxID        string `json:"tax_id"`
	CreatedTime  string `json:"created_time"`

	LastModifiedTime string `json:"last_modified_time"`

	Status    string `json:"status"`
	CC        string `json:"country_code"`
	POC       string `json:"place_of_contact"` // should be same as billing state
	TaxName   string `json:"tax_name"`         // IGST0
	LegalName string `json:"legal_name"`
	Country   string `json:"country"`
}

type BillingAddress struct {
	Attention string `json:"attention"`
	Address   string `json:"address"`
	Street2   string `json:"street2"`
	StateCode string `json:"state_code"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       string `json:"zip"`
	Country   string `json:"country"`
	Fax       string `json:"fax"`
	Phone     string `json:"phone"`
}

type ContactFindOptions struct {
	EmailContains string
}

// ContactParams struct represents the information to create a contact
type ContactParams struct {
	Name         string `json:"contact_name"`
	Company      string `json:"company_name,omitempty"`
	Website      string `json:"website,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
	ContactType  string `json:"contact_type,omitempty"`
	Notes        string `json:"notes,omitempty"`
	CurrencyID   string `json:"currency_id"`
	OwnerID      string `json:"owner_id"`

	ContactPersons  []ContactPerson `json:"contact_persons"`
	BillingAddress  BillingAddress  `json:"billing_address,omitempty"`
	ShippingAddress BillingAddress  `json:"shipping_address,omitempty"`

	// possible values ---> vat_registered,vat_not_registered,gcc_vat_not_registered,gcc_vat_registered,non_gcc,dz_vat_registered and dz_vat_not_registered.
	TaxTreatment string `json:"tax_treatment,omitempty"`
	GstNO        string `json:"gst_no,omitempty"`        // 15 digit
	GstTreatment string `json:"gst_treatment,omitempty"` // Allowed values are business_gst , business_none , overseas , consumer
	TaxID        string `json:"tax_id,omitempty"`

	Status        string  `json:"status,omitempty"`
	TaxPercentage float64 `json:"tax_percentage,omitempty"`
	CC            string  `json:"country_code,omitempty"`
	POC           string  `json:"place_of_contact,omitempty"` // should be same as billing state
	TaxName       string  `json:"tax_name,omitempty"`         // IGST0
	LegalName     string  `json:"legal_name,omitempty"`
	Country       string  `json:"country,omitempty"`
}

// New method will create a contact object and return a pointer to it
func (c *Contact) New() Resource {
	var obj = &Contact{}
	return obj
}

// Endpoint method returns the endpoint of the resource
func (c *Contact) Endpoint() string {
	return "/contacts"
}

// Create method will try to create a contact on razorpay
func (c *Contact) Create(params *ContactParams, client *Client) (*Contact, error) {
	var body, _ = json.Marshal(params)
	resp, err := client.Post(c.Endpoint(), string(body))

	respData, err := SendResp(resp, err, c)
	if err != nil {
		return c, err
	}
	return &respData.Contact, err
}

// FindOne tries to find the contact with given id
func (c *Contact) FindOne(id string, client *Client) (*Contact, error) {
	resp, err := client.Get(c.Endpoint() + "/" + id)
	respData, err := SendResp(resp, err, c)
	if err != nil {
		return c, err
	}
	return &respData.Contact, err
}

// FindAll tries to find the contacts with given options
func (c *Contact) FindAll(opts *ContactFindOptions, client *Client) ([]Contact, error) {
	resp, err := client.Get(c.Endpoint() + "?email_contains=" + opts.EmailContains)
	respData, err := SendResp(resp, err, c)

	var results []Contact
	if err != nil {
		return results, err
	}
	for _, ct := range respData.Contacts {
		results = append(results, ct)
	}
	return results, err
}

// Update method will try to update a invoice on razorpay
func (c *Contact) Update(id string, params *ContactParams, client *Client) (*Contact, error) {
	var body, _ = json.Marshal(params)
	resp, err := client.Put(c.Endpoint()+"/"+id, string(body))

	respData, err := SendResp(resp, err, c)
	if err != nil {
		return c, err
	}
	return &respData.Contact, err
}

// Delete tries to delete the contact with given id
func (c *Contact) Delete(id string, client *Client) error {
	resp, err := client.Delete(c.Endpoint() + "/" + id)
	respData, err := SendResp(resp, err, c)
	if err != nil {
		return err
	}
	if respData.Code == 0 {
		return nil
	}
	return errors.New(respData.Message)
}
