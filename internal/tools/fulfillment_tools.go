package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
)

func registerFulfillmentTools(add toolAdder) {
	// 1. tiktok_get_package_shipping_services
	add("tiktok_get_package_shipping_services",
		"Look up available shipping services, providers, and handover methods for a package.",
		baseProps(map[string]interface{}{
			"package_id": map[string]interface{}{
				"type":        "string",
				"description": "Package ID",
			},
		}),
		[]string{"package_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			packageID := strArg(args, "package_id")
			if packageID == "" {
				return mcp.ToolResultError("package_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/fulfillment/202309/packages/%s/shipping_services", url.PathEscape(packageID))
			data, _, err := client.Request(c, http.MethodGet, path, url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 2. tiktok_ship_package
	add("tiktok_ship_package",
		"Confirm handover / dispatch for a package (handover_method: PICKUP, DROP_OFF; optional tracking_number and shipping_provider_id).",
		baseProps(map[string]interface{}{
			"package_id": map[string]interface{}{
				"type":        "string",
				"description": "Package ID to ship",
			},
			"handover_method": map[string]interface{}{
				"type":        "string",
				"description": "Handover method: PICKUP, DROP_OFF (default: DROP_OFF)",
			},
			"shipping_provider_id": map[string]interface{}{
				"type":        "string",
				"description": "Optional shipping provider / carrier ID",
			},
			"tracking_number": map[string]interface{}{
				"type":        "string",
				"description": "Optional tracking number if seller shipping",
			},
		}),
		[]string{"package_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			packageID := strArg(args, "package_id")
			if packageID == "" {
				return mcp.ToolResultError("package_id is required")
			}
			reqBody := map[string]interface{}{
				"handover_method": strArg(args, "handover_method"),
			}
			if reqBody["handover_method"] == "" {
				reqBody["handover_method"] = "DROP_OFF"
			}
			if spID := strArg(args, "shipping_provider_id"); spID != "" {
				reqBody["shipping_provider_id"] = spID
			}
			if trk := strArg(args, "tracking_number"); trk != "" {
				reqBody["tracking_number"] = trk
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/fulfillment/202309/packages/%s/ship", url.PathEscape(packageID))
			data, _, err := client.Request(c, http.MethodPost, path, url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 3. tiktok_get_shipping_document
	add("tiktok_get_shipping_document",
		"Obtain printable shipping documents (shipping label PDF or packing slip) for a package.",
		baseProps(map[string]interface{}{
			"package_id": map[string]interface{}{
				"type":        "string",
				"description": "Package ID",
			},
			"document_type": map[string]interface{}{
				"type":        "string",
				"description": "Document type: SHIPPING_LABEL, PACKING_SLIP (default: SHIPPING_LABEL)",
			},
			"document_size": map[string]interface{}{
				"type":        "string",
				"description": "Document size: A6, A4 (default: A6)",
			},
		}),
		[]string{"package_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			packageID := strArg(args, "package_id")
			if packageID == "" {
				return mcp.ToolResultError("package_id is required")
			}
			q := url.Values{}
			docType := strArg(args, "document_type")
			if docType == "" {
				docType = "SHIPPING_LABEL"
			}
			q.Set("document_type", docType)
			docSize := strArg(args, "document_size")
			if docSize == "" {
				docSize = "A6"
			}
			q.Set("document_size", docSize)
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/fulfillment/202309/packages/%s/shipping_documents", url.PathEscape(packageID))
			data, _, err := client.Request(c, http.MethodGet, path, q, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 4. tiktok_get_package_tracking
	add("tiktok_get_package_tracking",
		"Track package transit events, checkpoint history, and current shipping status.",
		baseProps(map[string]interface{}{
			"package_id": map[string]interface{}{
				"type":        "string",
				"description": "Package ID to track",
			},
		}),
		[]string{"package_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			packageID := strArg(args, "package_id")
			if packageID == "" {
				return mcp.ToolResultError("package_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/fulfillment/202309/packages/%s/tracking", url.PathEscape(packageID))
			data, _, err := client.Request(c, http.MethodGet, path, url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 5. tiktok_split_package
	add("tiktok_split_package",
		"Split order order line items into multiple packages (scope: seller.fulfillment.package.write).",
		baseProps(map[string]interface{}{
			"order_id": map[string]interface{}{
				"type":        "string",
				"description": "Order ID to split",
			},
			"split_packages": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "object"},
				"description": "Array of package objects containing order_line_item_ids and quantities",
			},
		}),
		[]string{"order_id", "split_packages"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			orderID := strArg(args, "order_id")
			splitPkgs := sliceMapArg(args, "split_packages")
			if orderID == "" || len(splitPkgs) == 0 {
				return mcp.ToolResultError("order_id and split_packages are required")
			}
			reqBody := map[string]interface{}{
				"order_id":       orderID,
				"split_packages": splitPkgs,
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodPost, "/fulfillment/202309/packages/split", url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 6. tiktok_combine_packages
	add("tiktok_combine_packages",
		"Combine multiple packages into a single consolidated package (scope: seller.fulfillment.package.write).",
		baseProps(map[string]interface{}{
			"package_ids": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "List of package IDs to combine",
			},
		}),
		[]string{"package_ids"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			pkgIDs := sliceStrArg(args, "package_ids")
			if len(pkgIDs) < 2 {
				return mcp.ToolResultError("package_ids must contain at least 2 package IDs")
			}
			reqBody := map[string]interface{}{
				"package_ids": pkgIDs,
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodPost, "/fulfillment/202309/packages/combine", url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 7. tiktok_update_delivery_status
	add("tiktok_update_delivery_status",
		"Update delivery status for self-shipped packages (scope: seller.delivery.status.write).",
		baseProps(map[string]interface{}{
			"package_id": map[string]interface{}{
				"type":        "string",
				"description": "Package ID",
			},
			"delivery_status": map[string]interface{}{
				"type":        "string",
				"description": "Status: DELIVERED, DELIVERY_FAILED",
			},
			"update_time": map[string]interface{}{
				"type":        "integer",
				"description": "Unix timestamp in seconds when the status occurred",
			},
		}),
		[]string{"package_id", "delivery_status"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			packageID := strArg(args, "package_id")
			status := strArg(args, "delivery_status")
			if packageID == "" || status == "" {
				return mcp.ToolResultError("package_id and delivery_status are required")
			}
			reqBody := map[string]interface{}{
				"delivery_status": status,
			}
			if ut := intArg(args, "update_time", 0); ut > 0 {
				reqBody["update_time"] = ut
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/fulfillment/202309/packages/%s/delivery_status", url.PathEscape(packageID))
			data, _, err := client.Request(c, http.MethodPost, path, url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)
}
