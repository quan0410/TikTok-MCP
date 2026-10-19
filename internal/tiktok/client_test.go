package tiktok

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClient_Request(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-tts-access-token")
		if token != "valid_token" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"code": 105001, "message": "unauthorized"}`))
			return
		}
		q := r.URL.Query()
		if q.Get("shop_cipher") != "cipher_123" {
			t.Errorf("expected shop_cipher to be cipher_123, got %s", q.Get("shop_cipher"))
		}
		if q.Get("app_key") != "app_key_123" {
			t.Errorf("expected app_key to be app_key_123, got %s", q.Get("app_key"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code": 0, "message": "Success", "data": {"products": []}}`))
	}))
	defer ts.Close()

	creds := &Credentials{
		AccessToken: "valid_token",
		AppKey:      "app_key_123",
		AppSecret:   "app_secret_123",
		ShopCipher:  "cipher_123",
	}

	client := NewClient(creds)
	client.apiBase = ts.URL

	data, status, err := client.Request(context.Background(), http.MethodGet, "/product/202309/products/search", url.Values{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("expected status 200, got %d", status)
	}
	if len(data) == 0 {
		t.Errorf("expected non-empty response body")
	}
}

func TestClient_APIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code": 36001001, "message": "product not found"}`))
	}))
	defer ts.Close()

	creds := &Credentials{
		AccessToken: "valid_token",
		ShopCipher:  "cipher_123",
	}
	client := NewClient(creds)
	client.apiBase = ts.URL

	_, _, err := client.Request(context.Background(), http.MethodGet, "/product/202309/products/123", url.Values{}, nil)
	if err == nil {
		t.Fatalf("expected error from API response code != 0, got nil")
	}
}
