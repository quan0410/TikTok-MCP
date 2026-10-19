package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
)

func registerFinanceTools(add toolAdder) {
	// 1. tiktok_search_settlements
	add("tiktok_search_settlements",
		"Search and list financial settlements and payouts to the seller account (page_size, page_token).",
		baseProps(map[string]interface{}{
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
			q := url.Values{}
			q.Set("page_size", fmt.Sprint(intArg(args, "page_size", 20)))
			q.Set("sort_field", "create_time")
			if token := strArg(args, "page_token"); token != "" {
				q.Set("page_token", token)
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodGet, "/finance/202309/payments", q, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 2. tiktok_search_statements
	add("tiktok_search_statements",
		"Search financial statements and billing records for the shop.",
		baseProps(map[string]interface{}{
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
			q := url.Values{}
			q.Set("page_size", fmt.Sprint(intArg(args, "page_size", 20)))
			q.Set("sort_field", "statement_time")
			if token := strArg(args, "page_token"); token != "" {
				q.Set("page_token", token)
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodGet, "/finance/202309/statements", q, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 3. tiktok_get_order_statement_transactions
	add("tiktok_get_order_statement_transactions",
		"Retrieve financial transaction breakdowns (product sales, commission fee, shipping fee, taxes) for a specific order.",
		baseProps(map[string]interface{}{
			"order_id": map[string]interface{}{
				"type":        "string",
				"description": "Order ID",
			},
		}),
		[]string{"order_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			orderID := strArg(args, "order_id")
			if orderID == "" {
				return mcp.ToolResultError("order_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/finance/202309/orders/%s/statement_transactions", url.PathEscape(orderID))
			data, _, err := client.Request(c, http.MethodGet, path, url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)
}
