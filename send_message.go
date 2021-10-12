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
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrMessageTruncated  = errors.New("message truncated")
	ErrInvalidCoordinate = errors.New("invalid coordinate")
)

type Layout [6][22]int

func NewLayout() Layout {
	var layout Layout
	return layout
}

func (l *Layout) ValidCoordinate(x, y int) error {
	if x < 0 || y < 0 || x > 5 || y > 21 {
		return ErrInvalidCoordinate
	}
	return nil
}

func (l *Layout) Print(sx, sy int, s string) error {
	if err := l.ValidCoordinate(sx, sy); err != nil {
		return err
	}

	s = strings.ToUpper(s)
	// Preflight the string before an invalid set.
	if err := ValidText(s, false); err != nil {
		return err
	}

	x, y := sx, sy
	for _, c := range s {
		if x > 5 {
			return ErrMessageTruncated
		}
		l[x][y], _ = CharToCode(string(c))
		y++
		if y == 22 {
			x++
			y = 0
		}
	}

	return nil
}

func (l *Layout) SetColor(x, y int, c Color) error {
	if err := l.ValidCoordinate(x, y); err != nil {
		return err
	}
	if c < PoppyRed || c > White {
		return ErrInvalidColor
	}
	l[x][y] = int(c)
	return nil
}

type TextMessage struct {
	Text string `json:"text"`
}

type LayoutMessage struct {
	Layout Layout `json:"characters"`
}

type Message struct {
	ID      string `json:"id"`
	Created string `json:"created"`
	Text    string `json:"text,omitempty"`
}

type MessageResponse struct {
	Message `json:"message"`
}

func (c *Client) SendMessage(ctx context.Context, subscriptionID string, l Layout) (*MessageResponse, error) {
	var b bytes.Buffer
	body := &LayoutMessage{
		Layout: l,
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
	resp, err := c.do(req, &response)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return &response, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &response, nil
}

func (c *Client) SendText(ctx context.Context, subscriptionID string, text string) (*MessageResponse, error) {
	text = strings.ToUpper(text)
	if err := ValidText(text, true); err != nil {
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
	resp, err := c.do(req, &response)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return &response, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &response, nil
}
