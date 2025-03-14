package jitorpc

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

type TipAccount struct {
	Address string `json:"address"`
}

type BundleStatusResponse struct {
	Context struct {
		Slot int64 `json:"slot"`
	} `json:"context"`
	Value []struct {
		BundleID           string   `json:"bundle_id"`
		Transactions       []string `json:"transactions"`
		Slot               int64    `json:"slot"`
		ConfirmationStatus string   `json:"confirmation_status"`
		Err                struct {
			Ok interface{} `json:"Ok"`
		} `json:"err"`
	} `json:"value"`
}

func (c *JitoJsonRpcClient) GetTipAccounts() (json.RawMessage, error) {
	endpoint := "/bundles"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	request, err := c.sendRequest(endpoint, "getTipAccounts", nil)
	if err != nil {
		return nil, err
	}
	return request.Resp, err
}

func (c *JitoJsonRpcClient) GetRandomTipAccount() (*TipAccount, error) {
	rawResponse, err := c.GetTipAccounts()
	if err != nil {
		return nil, err
	}

	var tipAddresses []string
	err = json.Unmarshal(rawResponse, &tipAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tip accounts: %w", err)
	}

	if len(tipAddresses) == 0 {
		return nil, fmt.Errorf("no tip accounts available")
	}

	randomIndex := rand.Intn(len(tipAddresses))
	return &TipAccount{Address: tipAddresses[randomIndex]}, nil
}

func (c *JitoJsonRpcClient) GetBundleStatuses(bundleIds []string) (*BundleStatusResponse, error) {
	endpoint := "/bundles"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	params := [][]string{bundleIds}
	responseBody, err := c.sendRequest(endpoint, "getBundleStatuses", params)
	if err != nil {
		return nil, err
	}

	var response BundleStatusResponse
	if err := json.Unmarshal(responseBody.Resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bundle statuses: %w", err)
	}

	return &response, nil
}

func (c *JitoJsonRpcClient) SendBundle(params interface{}) (json.RawMessage, error) {
	endpoint := "/bundles"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	request, err := c.sendRequest(endpoint, "sendBundle", params)
	if err != nil {
		return nil, err
	}
	return request.Resp, nil
}

type InflightBundleStatuses struct {
	Context InflightBundleStatusContext `json:"context"`
	Value   []InflightBundleStatusValue `json:"value"`
}

type InflightBundleStatusContext struct {
	Slot int `json:"slot"`
}

type InflightBundleStatusValue struct {
	BundleIdD string `json:"bundle_id"`
	Status    string `json:"status"`
}

func (c *JitoJsonRpcClient) GetInflightBundleStatuses(params []string) (InflightBundleStatuses, error) {
	endpoint := "/bundles"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	request, err := c.sendRequest(endpoint, "getInflightBundleStatuses", [][]string{params})
	if err != nil {
		return InflightBundleStatuses{}, err
	}
	var statues InflightBundleStatuses
	err = json.Unmarshal(request.Resp, &statues)
	if err != nil {
		return InflightBundleStatuses{}, fmt.Errorf("failed to unmarshal inflight bundle statuses: %w", err)
	}
	return statues, nil
}
