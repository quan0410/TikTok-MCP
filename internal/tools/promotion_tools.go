package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
)

func registerPromotionTools(add toolAdder) {
	// 1. tiktok_search_promotions
	add("tiktok_search_promotions",
		"Search and list marketing promotions / discount activities on TikTok Shop (activity_type: DIRECT_DISCOUNT, FLASHSALE; status: UPCOMING, ONGOING, EXPIRED, DEACTIVATED).",
		baseProps(map[string]interface{}{
			"activity_type": map[string]interface{}{
				"type":        "string",
				"description": "Promotion type: DIRECT_DISCOUNT, FLASHSALE",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "Activity status: UPCOMING, ONGOING, EXPIRED, DEACTIVATED",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Page size (1-100, default 20)",
			},
			"page_token": map[string]interface{}{
				"type":        "string",
				"description": "Pagination token",
			},
		}),
		nil,
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			reqBody := map[string]interface{}{
				"page_size": intArg(args, "page_size", 20),
			}
			if actType := strArg(args, "activity_type"); actType != "" {
				reqBody["activity_type"] = actType
			}
			if status := strArg(args, "status"); status != "" {
				reqBody["status"] = status
			}
			if token := strArg(args, "page_token"); token != "" {
				reqBody["page_token"] = token
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodPost, "/promotion/202309/activities/search", url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 2. tiktok_get_promotion
	add("tiktok_get_promotion",
		"Retrieve full details of a promotional activity including participation products, discounts, and duration.",
		baseProps(map[string]interface{}{
			"activity_id": map[string]interface{}{
				"type":        "string",
				"description": "Activity ID",
			},
		}),
		[]string{"activity_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			actID := strArg(args, "activity_id")
			if actID == "" {
				return mcp.ToolResultError("activity_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/promotion/202309/activities/%s", url.PathEscape(actID))
			data, _, err := client.Request(c, http.MethodGet, path, url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 3. tiktok_update_promotion_status
	add("tiktok_update_promotion_status",
		"Deactivate an active or upcoming promotional activity.",
		baseProps(map[string]interface{}{
			"activity_id": map[string]interface{}{
				"type":        "string",
				"description": "Activity ID to deactivate",
			},
		}),
		[]string{"activity_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			actID := strArg(args, "activity_id")
			if actID == "" {
				return mcp.ToolResultError("activity_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/promotion/202309/activities/%s/deactivate", url.PathEscape(actID))
			data, _, err := client.Request(c, http.MethodPost, path, url.Values{}, map[string]interface{}{})
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)
}
