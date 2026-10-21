package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
)

func registerOrderTools(add toolAdder) {
	// 1. tiktok_search_orders
	add("tiktok_search_orders",
		"Search and list orders with filters (order_status: UNPAID, ON_HOLD, AWAITING_SHIPMENT, AWAITING_COLLECTION, IN_TRANSIT, DELIVERED, COMPLETED, CANCELLED; page_size, page_token).",
		baseProps(map[string]interface{}{
			"order_status": map[string]interface{}{
				"type":        "string",
				"description": "Filter by order status: UNPAID, ON_HOLD, AWAITING_SHIPMENT, AWAITING_COLLECTION, IN_TRANSIT, DELIVERED, COMPLETED, CANCELLED",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Page size (1-100, default 20)",
			},
			"page_token": map[string]interface{}{
				"type":        "string",
				"description": "Pagination token from previous search",
			},
			"sort_by": map[string]interface{}{
				"type":        "string",
				"description": "Sorting attribute: CREATE_TIME, UPDATE_TIME, PAID_TIME (default: CREATE_TIME)",
			},
			"sort_type": map[string]interface{}{
				"type":        "string",
				"description": "Sort direction: ASC, DESC (default: DESC)",
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
			if token := strArg(args, "page_token"); token != "" {
				q.Set("page_token", token)
			}
			reqBody := map[string]interface{}{}
			if status := strArg(args, "order_status"); status != "" {
				reqBody["order_status"] = status
			}
			if sortBy := strArg(args, "sort_by"); sortBy != "" {
				reqBody["sort_by"] = sortBy
			}
			if sortType := strArg(args, "sort_type"); sortType != "" {
				reqBody["sort_type"] = sortType
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodPost, "/order/202309/orders/search", q, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 2. tiktok_get_order
	add("tiktok_get_order",
		"Retrieve full details of an order including recipient address, payment info, package list, and line items.",
		baseProps(map[string]interface{}{
			"order_id": map[string]interface{}{
				"type":        "string",
				"description": "The unique TikTok Shop order ID",
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
			q := url.Values{}
			q.Set("ids", orderID)
			data, _, err := client.Request(c, http.MethodGet, "/order/202309/orders", q, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 3. tiktok_cancel_order
	add("tiktok_cancel_order",
		"Cancel an order with a seller reason (e.g., OUT_OF_STOCK, WRONG_PRICE).",
		baseProps(map[string]interface{}{
			"order_id": map[string]interface{}{
				"type":        "string",
				"description": "Order ID to cancel",
			},
			"cancel_reason": map[string]interface{}{
				"type":        "string",
				"description": "Cancellation reason code: OUT_OF_STOCK, WRONG_PRICE, BUYER_REQUEST, etc.",
			},
		}),
		[]string{"order_id", "cancel_reason"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			orderID := strArg(args, "order_id")
			reason := strArg(args, "cancel_reason")
			if orderID == "" || reason == "" {
				return mcp.ToolResultError("order_id and cancel_reason are required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/order/202309/orders/%s/cancel", url.PathEscape(orderID))
			reqBody := map[string]interface{}{
				"cancel_reason": reason,
			}
			data, _, err := client.Request(c, http.MethodPost, path, url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)
}
