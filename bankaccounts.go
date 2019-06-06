package zohobooks

// BankAccount struct will contain all the information of bank
type BankAccount struct {
	ID       string `json:"account_id"`
	Name     string `json:"account_name"`
	Type     string `json:"account_type"`
	IsActive bool   `json:"is_active"`
	BankName string `json:"bank_name"`
}

type BankAccountFindOptions struct {
	FilterBy, SortColumn string
}

// New method will create a payment object and return a pointer to it
func (ba *BankAccount) New() Resource {
	var obj = &BankAccount{}
	return obj
}

// Endpoint method returns the endpoint of the resource
func (ba *BankAccount) Endpoint() string {
	return "/bankaccounts"
}

func (ba *BankAccount) findAllEndpoint(opts *BankAccountFindOptions) string {
	endpoint := ba.Endpoint()
	if opts != nil {
		var append string
		if len(opts.FilterBy) > 0 {
			append += "filter_by=" + opts.FilterBy + "&"
		}
		if len(opts.SortColumn) > 0 {
			append += "sort_column=" + opts.SortColumn + "&"
		}
		if len(append) > 0 {
			endpoint += "?" + append
		}
	}
	return endpoint
}

// FindAll tries to find the contacts with given options
func (ba *BankAccount) FindAll(opts *BankAccountFindOptions, client *Client) ([]BankAccount, error) {
	resp, err := client.Get(ba.findAllEndpoint(opts))
	respData, err := sendResp(resp, err, ba)

	var results []BankAccount
	if err != nil {
		return results, err
	}
	for _, ct := range respData.BankAccounts {
		results = append(results, ct)
	}
	return results, err
}
