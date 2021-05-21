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
	"net/http"
	"strings"
)

const subscriptionsPath = "/subscriptions"

type Subscription struct {
	ID           string `json:"_id"`
	Created      string `json:"_created"`
	Installation `json:"installation"`
	Boards       []Board `json:"boards"`
}

type SubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

func (c *Client) Subscriptions(ctx context.Context) (*SubscriptionsResponse, error) {
	url := c.baseURL + subscriptionsPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set(APIKeyHeader, c.apiKey)
	req.Header.Set(APIKeySecret, c.apiSecret)

	var response SubscriptionsResponse
	_, err = c.do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type TextMessage struct {
	Text string `json:"text"`
}

type Message struct {
	ID      string `json:"id"`
	Created string `json:"created"`
	Text    string `json:"text,omitempty"`
}

type MessageResponse struct {
	Message `json:"message"`
}

func (c *Client) SendText(ctx context.Context, subscriptionID string, text string) (*MessageResponse, error) {
	text = strings.ToUpper(text)
	if err := ValidText(text); err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	var b bytes.Buffer
	body := &TextMessage{
		Text: text,
	}
	if err := json.NewEncoder(&b).Encode(body); err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", err)
	}

	url := fmt.Sprintf("%s%s/%s/message", c.baseURL, subscriptionsPath, subscriptionID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set(APIKeyHeader, c.apiKey)
	req.Header.Set(APIKeySecret, c.apiSecret)

	var response MessageResponse
	_, err = c.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
