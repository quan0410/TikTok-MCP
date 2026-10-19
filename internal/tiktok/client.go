package tiktok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultOpenAPIBase = "https://open-api.tiktokglobalshop.com"
)

// Client calls the TikTok Shop Open API, signed with HMAC-SHA256 and authenticated with
// the partner's access token and shop cipher.
type Client struct {
	creds      *Credentials
	httpClient *http.Client
	apiBase    string
}

// NewClient initializes a new TikTok Shop Open API client.
func NewClient(creds *Credentials) *Client {
	timeoutSec := 30
	if v := strings.TrimSpace(os.Getenv("TIKTOK_HTTP_TIMEOUT_SEC")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeoutSec = n
		}
	}
	base := strings.TrimSpace(os.Getenv("TIKTOK_API_BASE"))
	if base == "" {
		base = DefaultOpenAPIBase
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{
		Timeout:       5 * time.Second,
		KeepAlive:     30 * time.Second,
		FallbackDelay: 200 * time.Millisecond,
	}).DialContext
	transport.TLSHandshakeTimeout = 5 * time.Second
	transport.ResponseHeaderTimeout = 15 * time.Second

	return &Client{
		creds: creds,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   time.Duration(timeoutSec) * time.Second,
		},
		apiBase: strings.TrimRight(base, "/"),
	}
}

// Credentials returns the active credentials.
func (c *Client) Credentials() *Credentials {
	if c == nil {
		return nil
	}
	return c.creds
}

type apiResponseEnvelope struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	RequestID string          `json:"request_id"`
	Data      json.RawMessage `json:"data"`
}

// Request executes a signed HTTP request against the TikTok Shop Open API.
func (c *Client) Request(ctx context.Context, method, path string, query url.Values, body interface{}) ([]byte, int, error) {
	if c == nil || c.creds == nil {
		return nil, 0, fmt.Errorf("tiktok client not configured")
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if query == nil {
		query = url.Values{}
	}

	// Attach app_key if configured and not present
	if c.creds.AppKey != "" && query.Get("app_key") == "" {
		query.Set("app_key", c.creds.AppKey)
	}

	// Attach timestamp if not present
	if query.Get("timestamp") == "" {
		query.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	}

	// Pin shop_cipher from credentials for shop-scoped endpoints
	if c.creds.ShopCipher != "" && query.Get("shop_cipher") == "" && !isGlobalOnlyPath(path) {
		query.Set("shop_cipher", c.creds.ShopCipher)
	}

	// Marshal body if provided
	var payload []byte
	var reader io.Reader
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("encode request body: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	// Sign the request
	if c.creds.AppSecret != "" {
		sig := GenerateSignature(path, query, payload, c.creds.AppSecret)
		if sig != "" {
			query.Set("sign", sig)
		}
	}

	fullURL := c.apiBase + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reader)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("x-tts-access-token", c.creds.AccessToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	var resp *http.Response
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			var retryReader io.Reader
			if payload != nil {
				retryReader = bytes.NewReader(payload)
			}
			req, err = http.NewRequestWithContext(ctx, method, fullURL, retryReader)
			if err != nil {
				return nil, 0, err
			}
			req.Header.Set("x-tts-access-token", c.creds.AccessToken)
			req.Header.Set("Accept", "application/json")
			if body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
		}
		resp, err = c.httpClient.Do(req)
		if err == nil {
			break
		}
		if attempt == 2 {
			return nil, 0, fmt.Errorf("tiktok request failed: %w", err)
		}
		time.Sleep(300 * time.Millisecond)
	}
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, 10<<20) // 10 MiB cap
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := sanitizeError(string(data))
		if msg == "" || msg == "{}" {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				msg = "400 Bad Request"
			case http.StatusUnauthorized:
				msg = "401 Unauthorized (token expired or invalid)"
			case http.StatusForbidden:
				msg = "403 Forbidden (missing required scope or permission)"
			case http.StatusNotFound:
				msg = "404 Not Found"
			case http.StatusTooManyRequests:
				msg = "429 Too Many Requests (rate limit reached)"
			default:
				msg = fmt.Sprintf("HTTP %d error", resp.StatusCode)
			}
		}
		return data, resp.StatusCode, fmt.Errorf("tiktok HTTP %d: %s", resp.StatusCode, msg)
	}

	// Check TikTok response envelope code
	var env apiResponseEnvelope
	if err := json.Unmarshal(data, &env); err == nil {
		if env.Code != 0 {
			return data, resp.StatusCode, fmt.Errorf("tiktok API error (code %d): %s", env.Code, env.Message)
		}
	}

	return data, resp.StatusCode, nil
}

func isGlobalOnlyPath(path string) bool {
	return strings.HasPrefix(path, "/authorization/") || strings.HasPrefix(path, "/seller/")
}

func sanitizeError(raw string) string {
	s := strings.TrimSpace(raw)
	if len(s) > 300 {
		s = s[:300] + "..."
	}
	return s
}
