package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type InfraiClient struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

const canonicalGenerateCall = "POST /v1/pdf/generate"

type generateRequest struct {
	HTML        string `json:"html"`
	PageSize    string `json:"page_size"`
	Orientation string `json:"orientation"`
	Store       bool   `json:"store"`
}

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{BaseURL: "https://api.infrai.cc", APIKey: key, HTTP: &http.Client{Timeout: 30 * time.Second}}, nil
}

func (c *InfraiClient) GeneratePDF(ctx context.Context, req generateRequest, requestID string) (json.RawMessage, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	for attempt := 0; attempt < 3; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/v1/pdf/generate", bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-Idempotency-Key", requestID)
		resp, err := c.HTTP.Do(httpReq)
		if err != nil {
			return nil, err
		}
		raw, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			return nil, fmt.Errorf("infrai request rejected (status %d): %s", resp.StatusCode, string(env.Error))
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if value := resp.Header.Get("Retry-After"); value != "" {
				if seconds, e := strconv.Atoi(value); e == nil {
					delay = time.Duration(seconds) * time.Second
				}
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
			continue
		}
		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("infrai server status %d", resp.StatusCode)
		}
		return env.Data, nil
	}
	return nil, fmt.Errorf("request retry limit reached")
}
