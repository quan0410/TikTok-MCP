package tiktok

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// GenerateSignature generates the HMAC-SHA256 signature for TikTok Shop Open API requests.
// 1. Collect all query parameters except "sign" and "access_token".
// 2. Sort the parameter keys in alphabetical order (ASCII sort).
// 3. Concatenate each key with its value without delimiters.
// 4. Construct string to sign: appSecret + path + sortedParams + body + appSecret.
// 5. Sign with HMAC-SHA256 using appSecret, resulting in lowercase hex.
func GenerateSignature(path string, query url.Values, body []byte, appSecret string) string {
	if appSecret == "" {
		return ""
	}

	var keys []string
	for k := range query {
		if k == "sign" || k == "access_token" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		vals := query[k]
		for _, v := range vals {
			sb.WriteString(k)
			sb.WriteString(v)
		}
	}
	sortedParams := sb.String()

	bodyStr := ""
	if len(body) > 0 {
		bodyStr = string(body)
	}

	stringToSign := fmt.Sprintf("%s%s%s%s%s", appSecret, path, sortedParams, bodyStr, appSecret)

	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(stringToSign))
	return hex.EncodeToString(h.Sum(nil))
}
