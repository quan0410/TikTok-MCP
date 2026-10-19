package tiktok

import (
	"net/url"
	"testing"
)

func TestGenerateSignature(t *testing.T) {
	appSecret := "test_secret_123"
	path := "/product/202309/products/search"
	query := url.Values{}
	query.Set("app_key", "app_key_abc")
	query.Set("timestamp", "1700000000")
	query.Set("shop_cipher", "cipher_xyz")
	query.Set("sign", "should_be_ignored")
	query.Set("access_token", "should_also_be_ignored")

	body := []byte(`{"page_size":10}`)

	sig := GenerateSignature(path, query, body, appSecret)
	if sig == "" {
		t.Fatalf("expected non-empty signature")
	}

	// Deterministic verification
	sig2 := GenerateSignature(path, query, body, appSecret)
	if sig != sig2 {
		t.Fatalf("signature should be deterministic: %s != %s", sig, sig2)
	}
}

func TestGenerateSignature_EmptySecret(t *testing.T) {
	sig := GenerateSignature("/path", url.Values{}, nil, "")
	if sig != "" {
		t.Fatalf("expected empty signature when appSecret is empty")
	}
}
