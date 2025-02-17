package jitorpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type sendRequestResponse struct {
	Resp     json.RawMessage `json:"response"`
	BundleID string
}

func (c *JitoJsonRpcClient) sendRequest(endpoint, method string, params interface{}) (sendRequestResponse, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	request := JsonRpcRequest{
		JsonRpc: "2.0",
		ID:      1,
		Method:  method,
		Params:  params,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return sendRequestResponse{}, fmt.Errorf("error marshaling request: %w", err)
	}

	if c.isDebugEnabled() {
		fmt.Printf("Sending request to: %s\n", url)
		fmt.Printf("Request body: %s\n", string(requestBody))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return sendRequestResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.UUID != "" {
		req.Header.Set("x-jito-auth", c.UUID)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return sendRequestResponse{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if c.isDebugEnabled() {
		fmt.Printf("Response status: %s\n", resp.Status)
	}

	var jsonResp JsonRpcResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return sendRequestResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	if jsonResp.Error != nil {
		return sendRequestResponse{}, fmt.Errorf("RPC error: %s", jsonResp.Error.Message)
	}

	if c.isDebugEnabled() {
		fmt.Printf("Response body: %s\n", string(jsonResp.Result))
	}
	return sendRequestResponse{
		Resp:     jsonResp.Result,
		BundleID: resp.Header.Get("x-bundle-id"),
	}, nil
}
