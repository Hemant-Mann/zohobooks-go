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

// Client struct
type Client struct {
	Key   string
	OrgID string

	httpClient *http.Client
}

// Response struct to handle zohobooks response
type Response struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Contact Contact `json:"contact"`
	Invoice Invoice `json:"invoice"`
	Payment Payment `json:"payment"`

	Contacts   []Contact  `json:"contacts"`
	Currencies []Currency `json:"currencies"`
}

// Resource interface is to be used for generic decoding of object
type Resource interface {
	New() Resource
	Endpoint() string
}

// NewClient returns a pointer to the zohobooks client
func NewClient(key, orgID string) *Client {
	var c = &Client{
		Key:   key,
		OrgID: orgID,
	}
	c.httpClient = getHTTPClient(10)
	return c
}

func getHTTPClient(timeout int) *http.Client {
	var httpClient = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return httpClient
}

func sendResp(resp *http.Response, err error, rs Resource) (*Response, error) {
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
	if strings.Contains(path, "?") {
		return BaseURL + path + "&organization_id=" + c.OrgID
	}
	return BaseURL + path + "?organization_id=" + c.OrgID
}

func (c *Client) makeRequest(method, path string, body *bytes.Buffer, headers map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest(method, c.getURL(path), body)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Authorization", "Zoho-authtoken "+c.Key)
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
