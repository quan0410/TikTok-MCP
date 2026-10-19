package tools

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
)

func registerAnalyticsTools(add toolAdder) {
	// 1. tiktok_get_shop_performance
	add("tiktok_get_shop_performance",
		"Retrieve store business performance analytics (GMV, order count, page views, conversion rate).",
		baseProps(map[string]interface{}{
			"start_date": map[string]interface{}{
				"type":        "string",
				"description": "Start date in YYYY-MM-DD format",
			},
			"end_date": map[string]interface{}{
				"type":        "string",
				"description": "End date in YYYY-MM-DD format",
			},
		}),
		nil,
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			q := url.Values{}
			if s := strArg(args, "start_date"); s != "" {
				q.Set("start_date", s)
			}
			if e := strArg(args, "end_date"); e != "" {
				q.Set("end_date", e)
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodGet, "/analytics/202405/shop/performance", q, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 2. tiktok_get_bestsellers
	add("tiktok_get_bestsellers",
		"Retrieve trending and bestselling products ranking data on TikTok Shop (scope: seller.data.bestselling.public.read).",
		baseProps(map[string]interface{}{
			"category_id": map[string]interface{}{
				"type":        "string",
				"description": "Optional category ID filter",
			},
			"time_range": map[string]interface{}{
				"type":        "string",
				"description": "Time range: LAST_7_DAYS, LAST_28_DAYS (default: LAST_7_DAYS)",
			},
		}),
		nil,
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			q := url.Values{}
			if catID := strArg(args, "category_id"); catID != "" {
				q.Set("category_id", catID)
			}
			tr := strArg(args, "time_range")
			if tr == "" {
				tr = "LAST_7_DAYS"
			}
			q.Set("time_range", tr)
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodGet, "/analytics/202511/products/bestselling", q, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)
}
