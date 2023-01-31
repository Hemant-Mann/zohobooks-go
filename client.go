package zohobooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BaseURL stores the API base URL
const BaseURL = "https://books.zoho.com/api/v3"
const OAuthURL = "https://accounts.zoho.com/oauth/v2/token"
const OAuthURLIn = "https://accounts.zoho.in/oauth/v2/token"

// Client struct
type Client struct {
	Key          string
	OAuthToken   string
	clientID     string
	clientSecret string
	refreshToken string
	redirectURI  string
	OrgID        string
	Datacenter   string
	httpClient   *http.Client
}

type ClientConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	OAuthToken   string
	RefreshToken string
	Datacenter   string
	OrgID        string
	Timeout      int
}

type OAuthResponse struct {
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
	APIDomain   string `json:"api_domain"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Response struct to handle zohobooks response
type Response struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Contact Contact `json:"contact"`
	Invoice Invoice `json:"invoice"`
	Payment Payment `json:"payment"`

	BankTransaction BankTransaction `json:"banktransaction"`

	Contacts     []Contact     `json:"contacts"`
	Payments     []Payment     `json:"customerpayments"`
	Currencies   []Currency    `json:"currencies"`
	BankAccounts []BankAccount `json:"bankaccounts"`
	Data         zohoRespError `json:"data"`
}

type zohoRespError struct {
	Errors []response2 `json:"errors"`
}

type response2 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Resource interface is to be used for generic decoding of object
type Resource interface {
	New() Resource
	Endpoint() string
}

// NewClient returns a pointer to the zohobooks client
func NewClient(key, orgID, datacenter string) *Client {
	var c = &Client{
		Key:   key,
		OrgID: orgID,

		Datacenter: datacenter,
	}
	c.httpClient = getHTTPClient(10)
	return c
}

// NewClientWithTimeout returns a pointer to the zohobooks client
func NewClientWithTimeout(key, orgID, datacenter string, timeout int) *Client {
	var c = &Client{
		Key:   key,
		OrgID: orgID,

		Datacenter: datacenter,
	}
	c.httpClient = getHTTPClient(timeout)
	return c
}

// NewClientWithTimeout returns a pointer to the zohobooks client
func NewClientWithConfig(conf *ClientConfig) *Client {
	var c = &Client{
		OrgID:        conf.OrgID,
		OAuthToken:   conf.OAuthToken,
		Datacenter:   conf.Datacenter,
		clientID:     conf.ClientID,
		clientSecret: conf.ClientSecret,
		redirectURI:  conf.RedirectURI,
		refreshToken: conf.RefreshToken,
	}
	c.httpClient = getHTTPClient(conf.Timeout)
	return c
}

// GetBaseURL will return the base URL for zohobooks based on the specified
// datacenter while initializing the client
func (c *Client) GetBaseURL() string {
	if c.Datacenter == "in" || c.Datacenter == "IN" {
		return "https://books.zoho.in/api/v3"
	}
	if c.Datacenter == "eu" {
		return "https://books.zoho.eu/api/v3"
	}
	if c.Datacenter == "au" {
		return "https://books.zoho.com.au/api/v3"
	}
	return "https://books.zoho.com/api/v3"
}

func getHTTPClient(timeout int) *http.Client {
	var httpClient = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return httpClient
}

func SendResp(resp *http.Response, err error, rs Resource) (*Response, error) {
	var newResp = &Response{}
	if err != nil {
		return newResp, err
	}
	body, readErr := readBody(resp)
	if readErr != nil {
		return newResp, readErr
	}
	parseError := json.Unmarshal(body, newResp)
	if parseError == nil && newResp.Code > 0 {
		return newResp, errors.New(newResp.Message)
	}
	return newResp, parseError
}

func readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func (c *Client) getURL(path string) string {
	var baseURL = c.GetBaseURL()
	if strings.Contains(path, "?") {
		return baseURL + path + "&organization_id=" + c.OrgID
	}
	return baseURL + path + "?organization_id=" + c.OrgID
}

func (c *Client) makeRequest(method, path string, body *bytes.Buffer, headers map[string]string) (*http.Response, error) {
	if len(c.OAuthToken) == 0 || len(c.OrgID) == 0 {
		return nil, errors.New("missing oauthtoken or org id")
	}
	req, _ := http.NewRequest(method, c.getURL(path), body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if len(c.OAuthToken) > 0 {
		req.Header.Set("Authorization", "Zoho-oauthtoken "+c.OAuthToken)
	} else if len(c.Key) > 0 { // Zoho authtoken are deprecated use oauthtokens
		req.Header.Set("Authorization", "Zoho-authtoken "+c.Key)
	}
	resp, err := c.httpClient.Do(req)
	return resp, err
}

// Get method makes a GET request to the resource
func (c *Client) Get(path string) (*http.Response, error) {
	return c.makeRequest("GET", path, bytes.NewBuffer([]byte("")), nil)
}

// Post method makes a POST and sends data in json format
func (c *Client) Post(path string, body string) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
	}
	data := url.Values{}
	data.Set("JSONString", body)
	byteBody := []byte(data.Encode())
	return c.makeRequest("POST", path, bytes.NewBuffer(byteBody), headers)
}

// Put method makes a PUT and sends data in json format
func (c *Client) Put(path string, body string) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
	}
	data := url.Values{}
	data.Set("JSONString", body)
	byteBody := []byte(data.Encode())
	return c.makeRequest("PUT", path, bytes.NewBuffer(byteBody), headers)
}

// Delete method makes a DELETE request to the resource
func (c *Client) Delete(path string) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
	}
	return c.makeRequest("DELETE", path, bytes.NewBuffer([]byte("")), headers)
}

func (c *Client) GetOauthURL() string {
	if c.Datacenter == "in" || c.Datacenter == "IN" {
		return OAuthURLIn
	}
	return OAuthURL
}

func (c *Client) GenAccessToken() (string, error) {
	var query = "refresh_token=" + c.refreshToken + "&client_id=" + c.clientID + "&client_secret=" + c.clientSecret + "&redirect_uri=" + c.redirectURI + "&grant_type=refresh_token"
	req, err := http.NewRequest("POST", c.GetOauthURL()+"?"+query, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	body, readErr := readBody(resp)
	if readErr != nil {
		return "", readErr
	}
	oauthResp := &OAuthResponse{}
	parseError := json.Unmarshal(body, oauthResp)
	if parseError != nil {
		return "", parseError
	}
	if oauthResp.Error != "" {
		return "", errors.New(oauthResp.Error)
	}
	return oauthResp.AccessToken, nil
}
