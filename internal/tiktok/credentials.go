package tiktok

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Credentials matches datumbridge-integrations' TikTokShopTokenBundle.
type Credentials struct {
	Type         string    `json:"type"`
	AppKey       string    `json:"app_key"`
	AppSecret    string    `json:"app_secret,omitempty"`
	OpenID       string    `json:"open_id,omitempty"`
	SellerName   string    `json:"seller_name,omitempty"`
	Region       string    `json:"region,omitempty"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	ShopID       string    `json:"shop_id,omitempty"`
	ShopCipher   string    `json:"shop_cipher,omitempty"`
	ShopName     string    `json:"shop_name,omitempty"`
}

// ParseCredentials loads Credentials from vault args (credentials_json /
// credentials_path), then falls back to process env for local/dev runs.
// Precedence: injected bundle first, TIKTOK_SHOP_* / TIKTOK_* env fills blanks.
func ParseCredentials(credentialsJSON, credentialsPath string) (*Credentials, error) {
	c := &Credentials{}
	raw := strings.TrimSpace(credentialsJSON)

	if raw == "" && strings.TrimSpace(credentialsPath) != "" {
		path := filepath.Clean(credentialsPath)
		if strings.Contains(path, "..") {
			return nil, fmt.Errorf("credentials_path must not contain ..")
		}
		allowedRoot := strings.TrimSpace(os.Getenv("TIKTOK_CREDENTIALS_DIR"))
		if allowedRoot == "" {
			allowedRoot = "/credentials"
		}
		allowedRoot = filepath.Clean(allowedRoot)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("credentials_path: %w", err)
		}
		absRoot, err := filepath.Abs(allowedRoot)
		if err != nil {
			return nil, fmt.Errorf("TIKTOK_CREDENTIALS_DIR: %w", err)
		}
		resolvedRoot, err := filepath.EvalSymlinks(absRoot)
		if err != nil {
			resolvedRoot = absRoot
		}
		resolvedPath, err := filepath.EvalSymlinks(absPath)
		if err != nil {
			return nil, fmt.Errorf("credentials_path: %w", err)
		}
		if resolvedPath != resolvedRoot && !strings.HasPrefix(resolvedPath, resolvedRoot+string(os.PathSeparator)) {
			return nil, fmt.Errorf("credentials_path must be under %s", absRoot)
		}
		b, err := os.ReadFile(resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("read credentials_path: %w", err)
		}
		raw = strings.TrimSpace(string(b))
	}

	if raw != "" {
		if err := json.Unmarshal([]byte(raw), c); err != nil {
			return nil, fmt.Errorf("invalid credentials_json: %w", err)
		}
	}

	if strings.TrimSpace(c.AccessToken) == "" {
		if envToken := os.Getenv("TIKTOK_SHOP_ACCESS_TOKEN"); envToken != "" {
			c.AccessToken = envToken
		} else if envToken := os.Getenv("TIKTOK_ACCESS_TOKEN"); envToken != "" {
			c.AccessToken = envToken
		}
	}
	if strings.TrimSpace(c.AppKey) == "" {
		if envKey := os.Getenv("TIKTOK_SHOP_APP_KEY"); envKey != "" {
			c.AppKey = envKey
		} else if envKey := os.Getenv("TIKTOK_APP_KEY"); envKey != "" {
			c.AppKey = envKey
		}
	}
	if strings.TrimSpace(c.AppSecret) == "" {
		if envSecret := os.Getenv("TIKTOK_SHOP_APP_SECRET"); envSecret != "" {
			c.AppSecret = envSecret
		} else if envSecret := os.Getenv("TIKTOK_APP_SECRET"); envSecret != "" {
			c.AppSecret = envSecret
		}
	}
	if strings.TrimSpace(c.ShopCipher) == "" {
		if envCipher := os.Getenv("TIKTOK_SHOP_CIPHER"); envCipher != "" {
			c.ShopCipher = envCipher
		} else if envCipher := os.Getenv("TIKTOK_CIPHER"); envCipher != "" {
			c.ShopCipher = envCipher
		}
	}
	if strings.TrimSpace(c.ShopID) == "" {
		if envShopID := os.Getenv("TIKTOK_SHOP_ID"); envShopID != "" {
			c.ShopID = envShopID
		}
	}

	if strings.TrimSpace(c.AccessToken) == "" {
		return nil, fmt.Errorf("missing TikTok Shop access_token in credentials_json or TIKTOK_SHOP_ACCESS_TOKEN")
	}

	return c, nil
}
