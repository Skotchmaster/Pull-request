package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type testEnv struct {
	t       *testing.T
	client  *http.Client
	baseURL string
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	return &testEnv{
		t: t,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: baseURL,
	}
}


func (env *testEnv) doJSON(method, path string, reqBody any) (*http.Response, []byte) {
	env.t.Helper()

	var body io.Reader
	if reqBody != nil {
		raw, err := json.Marshal(reqBody)
		if err != nil {
			env.t.Fatalf("marshal %s %s body: %v", method, path, err)
		}
		body = bytes.NewReader(raw)
	}

	url := env.baseURL + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		env.t.Fatalf("new request %s %s: %v", method, url, err)
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := env.client.Do(req)
	if err != nil {
		env.t.Fatalf("do request %s %s: %v", method, url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		env.t.Fatalf("read response body: %v", err)
	}

	return resp, respBody
}