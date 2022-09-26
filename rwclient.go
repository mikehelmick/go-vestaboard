// Copyright 2021 Mike Helmick
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vestaboard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	RWApiKeyHeader = "X-Vestaboard-Read-Write-Key"
)

type RWClient struct {
	rwKey      string
	httpClient *http.Client
	baseURL    string
}

func NewRWClient(rwKey string) *RWClient {
	return &RWClient{
		rwKey: rwKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: "https://rw.vestaboard.com",
	}
}

func (c *RWClient) do(req *http.Request, out interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	errPrefix := fmt.Sprintf("%s %s - %d", strings.ToUpper(req.Method), req.URL.String(), resp.StatusCode)

	r := io.LimitReader(resp.Body, MaxBodySize)
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read body: %w", errPrefix, err)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return nil, fmt.Errorf("%s: response content-type is not application/json (got %s): body: %s",
			errPrefix, ct, body)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return nil, fmt.Errorf("%s: failed to decode JSON response: %w: body: %s",
			errPrefix, err, body)
	}
	return resp, nil
}

type RWMessageResponse struct {
	Message `json:"status"`
}

func (c *RWClient) SendMessage(ctx context.Context, l Layout) (*MessageResponse, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(l); err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set(RWApiKeyHeader, c.rwKey)

	var response RWMessageResponse
	resp, err := c.do(req, &response)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	mr := MessageResponse(response)

	return &mr, nil
}
