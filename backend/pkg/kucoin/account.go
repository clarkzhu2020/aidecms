package kucoin

import (
	"fmt"
	"net/url"
)

// Account 账户信息
type Account struct {
	Id             string `json:"id"`
	Currency       string `json:"currency"`
	Type           string `json:"type"`           // main, trade, margin
	Balance        string `json:"balance"`
	Available      string `json:"available"`
	Holds          string `json:"holds"`
	BaseCurrency   string `json:"baseCurrency"`   // 只适用于杠杆账户
	BaseCurrencyBalance string `json:"baseCurrencyBalance"` // 只适用于杠杆账户
	BaseCurrencyAvailable string `json:"baseCurrencyAvailable"` // 只适用于杠杆账户
	BaseCurrencyHolds string `json:"baseCurrencyHolds"` // 只适用于杠杆账户
	BaseCurrencyHold string `json:"baseCurrencyHold"` // 只适用于杠杆账户
}

// GetAccounts 获取账户列表
func (c *Client) GetAccounts(currency, accountType string) ([]Account, error) {
	endpoint := "/api/v1/accounts"

	params := url.Values{}
	if currency != "" {
		params.Add("currency", currency)
	}
	if accountType != "" {
		params.Add("type", accountType)
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	if err := c.httpClient.unmarshalJSON(respBody, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// GetAccountDetail 获取账户详情
func (c *Client) GetAccountDetail(accountId string) (*Account, error) {
	endpoint := fmt.Sprintf("/api/v1/accounts/%s", accountId)

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := c.httpClient.unmarshalJSON(respBody, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// AccountLedger 账户流水
type AccountLedger struct {
	Amount          string `json:"amount"`
	Fee             string `json:"fee"`
	Balance         string `json:"balance"`
	BillType        string `json:"billType"` // DEPOSIT, WITHDRAW, TRADE, TRANSFER, STAKING
	BizType         string `json:"bizType"`
		Context        string `json:"context"`
		CreatedAt      int64  `json:"createdAt"`
}

// GetAccountLedger 获取账户流水
func (c *Client) GetAccountLedger(accountId, currency, direction string, startAt, endAt, page, pageSize int64) ([]AccountLedger, error) {
	endpoint := fmt.Sprintf("/api/v1/accounts/%s/ledgers", accountId)

	params := url.Values{}
	if currency != "" {
		params.Add("currency", currency)
	}
	if direction != "" {
		params.Add("direction", direction)
	}
	if startAt > 0 {
		params.Add("startAt", fmt.Sprintf("%d", startAt))
	}
	if endAt > 0 {
		params.Add("endAt", fmt.Sprintf("%d", endAt))
	}
	if page > 0 {
		params.Add("page", fmt.Sprintf("%d", page))
	}
	if pageSize > 0 {
		params.Add("pageSize", fmt.Sprintf("%d", pageSize))
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var response struct {
		CurrentPage int              `json:"currentPage"`
		PageSize    int              `json:"pageSize"`
		TotalNum    int              `json:"totalNum"`
		TotalPage   int              `json:"totalPage"`
		Items       []AccountLedger  `json:"items"`
	}

	if err := c.httpClient.unmarshalJSON(respBody, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// AccountHold 账户冻结记录
type AccountHold struct {
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	OrderId        string `json:"orderId"`
	TransactionId  string `json:"transactionId"`
	Type           string `json:"type"` // OPEN_ORDER, FREEZE, FREEZE_OPEN_ORDER
	UpdatedAt      int64  `json:"updatedAt"`
	CreatedAt      int64  `json:"createdAt"`
}

// GetAccountHolds 获取账户冻结记录
func (c *Client) GetAccountHolds(accountId string) ([]AccountHold, error) {
	endpoint := fmt.Sprintf("/api/v1/accounts/%s/holds", accountId)

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var holds []AccountHold
	if err := c.httpClient.unmarshalJSON(respBody, &holds); err != nil {
		return nil, err
	}

	return holds, nil
}

// InnerTransferRequest 内部转账请求
type InnerTransferRequest struct {
	ClientOid   string `json:"clientOid"`
	Currency    string `json:"currency"`
	From        string `json:"from"`    // main, trade, margin
	To          string `json:"to"`      // main, trade, margin
	Amount      string `json:"amount"`
}

// InnerTransfer 内部转账
func (c *Client) InnerTransfer(req *InnerTransferRequest) error {
	endpoint := "/api/v2/accounts/inner-transfer"

	body, err := c.httpClient.marshalJSON(req)
	if err != nil {
		return err
	}

	_, err = c.httpClient.Post(endpoint, body, true)
	if err != nil {
		return err
	}

	return nil
}

// marshalJSON 序列化JSON
func (c *HTTPClient) marshalJSON(v interface{}) ([]byte, error) {
	// 简单的序列化方法，实际可能需要更复杂的实现
	return nil, nil
}
