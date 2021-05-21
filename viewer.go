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
	"context"
	"net/http"
)

const viewerPath = "/viewer"

type Installable struct {
	ID string `json:"_id"`
}

type Installation struct {
	ID          string `json:"_id"`
	Installable `json:"installable"`
}

type ViewerResponse struct {
	Type         string `json:"type"`
	ID           string `json:"_id"`
	Created      string `json:"_created"`
	Installation `json:"installation"`
}

func (c *Client) Viewer(ctx context.Context) (*ViewerResponse, error) {
	url := c.baseURL + viewerPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set(APIKeyHeader, c.apiKey)
	req.Header.Set(APIKeySecret, c.apiSecret)

	var response ViewerResponse
	_, err = c.do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
