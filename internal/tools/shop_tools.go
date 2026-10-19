package tools

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
)

func registerShopTools(add toolAdder) {
	// 1. tiktok_get_authorized_shops
	add("tiktok_get_authorized_shops",
		"Retrieve the list of authorized shops and shop ciphers associated with the access token.",
		baseProps(map[string]interface{}{}),
		nil,
		func(raw json.RawMessage) map[string]interface{} {
			client, _, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodGet, "/authorization/202309/shops", url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 2. tiktok_get_shop_profile
	add("tiktok_get_shop_profile",
		"Retrieve the active shop profile, status, region, and associated warehouses.",
		baseProps(map[string]interface{}{}),
		nil,
		func(raw json.RawMessage) map[string]interface{} {
			client, _, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodGet, "/seller/202309/shops", url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)
}
