package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"fleetglance/internal/protocol"
)

const telemetryPath = "/api/telemetry"

type Client struct {
	httpClient *http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Get(ctx context.Context, shipUrl string) (*protocol.Telemetry, error) {
	url := appendPath(shipUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create telemetry request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get telemetry: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get telemetry: unexpected status %d", res.StatusCode)
	}

	var response protocol.Response[protocol.Telemetry]
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode telemetry response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("telemetry response error: %s", response.Error.Message)
	}

	if response.Data == nil {
		return nil, fmt.Errorf("telemetry response missing data")
	}

	return response.Data, nil
}

func appendPath(raw string) string {
	if strings.HasSuffix(raw, telemetryPath) {
		return raw
	}

	return strings.TrimRight(raw, "/") + telemetryPath
}
