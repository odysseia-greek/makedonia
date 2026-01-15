package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type response struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// Execute sends a GraphQL POST request to url and unmarshals the "data" object into v.
func Execute(ctx context.Context, url, query string, variables map[string]any, v any) error {
	body, _ := json.Marshal(&request{Query: query, Variables: variables})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("graphql http status %d: %s", resp.StatusCode, string(b))
	}
	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if len(r.Errors) > 0 {
		return fmt.Errorf("graphql errors: %v", r.Errors[0].Message)
	}
	if v == nil {
		return nil
	}
	return json.Unmarshal(r.Data, v)
}
