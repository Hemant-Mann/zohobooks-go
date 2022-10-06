package zohobooks

// Currency struct represents the information of the currency
type Currency struct {
	ID             string  `json:"currency_id"`
	Code           string  `json:"currency_code"`
	Name           string  `json:"currency_name"`
	Symbol         string  `json:"currency_symbol"`
	PricePrecision int     `json:"price_precision"`
	IsBaseCurrency bool    `json:"is_base_currency"`
	ExchangeRate   float64 `json:"exchange_rate"`
	EffectiveDate  string  `json:"effective_date"`
}

// New method will create a contact object and return a pointer to it
func (c *Currency) New() Resource {
	var obj = &Currency{}
	return obj
}

// Endpoint method returns the endpoint of the resource
func (c *Currency) Endpoint() string {
	return "/settings/currencies"
}

// FindAll will return the list of currencies present in zohobooks org
func (c *Currency) FindAll(client *Client) ([]Currency, error) {
	var results []Currency
	resp, err := client.Get(c.Endpoint())
	respData, err := SendResp(resp, err, c)
	if err != nil {
		return results, err
	}
	for _, cr := range respData.Currencies {
		results = append(results, cr)
	}
	return results, nil
}
