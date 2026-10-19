package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
)

func registerReturnTools(add toolAdder) {
	// 1. tiktok_search_returns
	add("tiktok_search_returns",
		"Search return and refund requests (return_type: REFUND_ONLY, RETURN_AND_REFUND; return_status: AWAITING_SELLER_REVIEW, AWAITING_BUYER_RETURN, etc.).",
		baseProps(map[string]interface{}{
			"return_status": map[string]interface{}{
				"type":        "string",
				"description": "Status filter: AWAITING_SELLER_REVIEW, AWAITING_BUYER_RETURN, AWAITING_SELLER_RECEIVE, COMPLETE, CLOSED",
			},
			"return_type": map[string]interface{}{
				"type":        "string",
				"description": "Filter by return type: REFUND_ONLY, RETURN_AND_REFUND",
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
			q := url.Values{}
			q.Set("page_size", fmt.Sprint(intArg(args, "page_size", 20)))
			if token := strArg(args, "page_token"); token != "" {
				q.Set("page_token", token)
			}
			reqBody := map[string]interface{}{}
			if status := strArg(args, "return_status"); status != "" {
				reqBody["return_status"] = status
			}
			if rType := strArg(args, "return_type"); rType != "" {
				reqBody["return_type"] = rType
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodPost, "/return_refund/202309/returns/search", q, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 2. tiktok_get_return
	add("tiktok_get_return",
		"Retrieve full details of a return or refund request including buyer reason, amount, proof images, and timeline.",
		baseProps(map[string]interface{}{
			"return_id": map[string]interface{}{
				"type":        "string",
				"description": "The return order ID",
			},
		}),
		[]string{"return_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			returnID := strArg(args, "return_id")
			if returnID == "" {
				return mcp.ToolResultError("return_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/return_refund/202309/returns/%s", url.PathEscape(returnID))
			data, _, err := client.Request(c, http.MethodGet, path, url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 3. tiktok_approve_return
	add("tiktok_approve_return",
		"Approve a return request or refund request.",
		baseProps(map[string]interface{}{
			"return_id": map[string]interface{}{
				"type":        "string",
				"description": "The return order ID to approve",
			},
			"decision": map[string]interface{}{
				"type":        "string",
				"description": "Decision type: AGREE_REFUND, AGREE_RETURN (default: AGREE_REFUND)",
			},
		}),
		[]string{"return_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			returnID := strArg(args, "return_id")
			if returnID == "" {
				return mcp.ToolResultError("return_id is required")
			}
			decision := strArg(args, "decision")
			if decision == "" {
				decision = "AGREE_REFUND"
			}
			reqBody := map[string]interface{}{
				"decision": decision,
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/return_refund/202309/returns/%s/approve", url.PathEscape(returnID))
			data, _, err := client.Request(c, http.MethodPost, path, url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 4. tiktok_reject_return
	add("tiktok_reject_return",
		"Reject a return or refund request with a specified rejection reason code and comments.",
		baseProps(map[string]interface{}{
			"return_id": map[string]interface{}{
				"type":        "string",
				"description": "The return order ID to reject",
			},
			"reject_reason": map[string]interface{}{
				"type":        "string",
				"description": "Rejection reason code",
			},
			"reject_comments": map[string]interface{}{
				"type":        "string",
				"description": "Comments or explanations for the buyer",
			},
		}),
		[]string{"return_id", "reject_reason"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			returnID := strArg(args, "return_id")
			reason := strArg(args, "reject_reason")
			if returnID == "" || reason == "" {
				return mcp.ToolResultError("return_id and reject_reason are required")
			}
			reqBody := map[string]interface{}{
				"reject_reason": reason,
			}
			if comments := strArg(args, "reject_comments"); comments != "" {
				reqBody["reject_comments"] = comments
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/return_refund/202309/returns/%s/reject", url.PathEscape(returnID))
			data, _, err := client.Request(c, http.MethodPost, path, url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)
}
