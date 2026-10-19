package tiktok

import (
	"os"
	"testing"
)

func TestParseCredentials_JSON(t *testing.T) {
	jsonBlob := `{
		"type": "oauth",
		"app_key": "my_app_key",
		"access_token": "token_123",
		"shop_cipher": "cipher_456",
		"shop_id": "shop_789"
	}`

	creds, err := ParseCredentials(jsonBlob, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.AccessToken != "token_123" {
		t.Errorf("expected access_token token_123, got %s", creds.AccessToken)
	}
	if creds.ShopCipher != "cipher_456" {
		t.Errorf("expected shop_cipher cipher_456, got %s", creds.ShopCipher)
	}
	if creds.AppKey != "my_app_key" {
		t.Errorf("expected app_key my_app_key, got %s", creds.AppKey)
	}
}

func TestParseCredentials_EnvFallback(t *testing.T) {
	os.Setenv("TIKTOK_SHOP_ACCESS_TOKEN", "env_token_abc")
	os.Setenv("TIKTOK_SHOP_CIPHER", "env_cipher_def")
	os.Setenv("TIKTOK_SHOP_APP_KEY", "env_key_ghi")
	defer func() {
		os.Unsetenv("TIKTOK_SHOP_ACCESS_TOKEN")
		os.Unsetenv("TIKTOK_SHOP_CIPHER")
		os.Unsetenv("TIKTOK_SHOP_APP_KEY")
	}()

	creds, err := ParseCredentials("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.AccessToken != "env_token_abc" {
		t.Errorf("expected token env_token_abc, got %s", creds.AccessToken)
	}
	if creds.ShopCipher != "env_cipher_def" {
		t.Errorf("expected cipher env_cipher_def, got %s", creds.ShopCipher)
	}
}

func TestParseCredentials_MissingToken(t *testing.T) {
	os.Unsetenv("TIKTOK_SHOP_ACCESS_TOKEN")
	os.Unsetenv("TIKTOK_ACCESS_TOKEN")

	_, err := ParseCredentials("", "")
	if err == nil {
		t.Fatalf("expected error on missing token, got nil")
	}
}
