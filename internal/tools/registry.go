package tools

import (
	"strings"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
	registrydocs "github.com/datumbridge/tiktok-mcp/registry_docs"
)

var toolOutputs = map[string]string{
	// Products (12)
	"tiktok_search_products":    "products, total_count, next_page_token",
	"tiktok_get_product":        "product, id, title, status, skus",
	"tiktok_create_product":     "product_id, skus",
	"tiktok_update_product":     "product_id",
	"tiktok_delete_product":     "product_ids",
	"tiktok_activate_product":   "product_ids",
	"tiktok_deactivate_product": "product_ids",
	"tiktok_update_inventory":   "failed_skus",
	"tiktok_update_price":       "failed_skus",
	"tiktok_get_categories":     "categories",
	"tiktok_get_category_rules": "rules, attributes",
	"tiktok_diagnose_product":   "diagnoses, suggestions",

	// Orders (3)
	"tiktok_search_orders": "orders, total_count, next_page_token",
	"tiktok_get_order":     "order, id, status, recipient_address, line_items, packages",
	"tiktok_cancel_order":  "order_id",

	// Fulfillment & Logistics (7)
	"tiktok_get_package_shipping_services": "shipping_services, providers",
	"tiktok_ship_package":                  "package_id, tracking_number",
	"tiktok_get_shipping_document":         "doc_url",
	"tiktok_get_package_tracking":          "tracking_info, checkpoints",
	"tiktok_split_package":                 "packages",
	"tiktok_combine_packages":              "package_id",
	"tiktok_update_delivery_status":        "package_id",

	// Returns & Refunds (4)
	"tiktok_search_returns": "returns, total_count, next_page_token",
	"tiktok_get_return":     "return_order, id, reason, return_amount, proof_images",
	"tiktok_approve_return": "return_id",
	"tiktok_reject_return":  "return_id",

	// Promotions (3)
	"tiktok_search_promotions":       "activities, total_count, next_page_token",
	"tiktok_get_promotion":           "activity, id, title, status, discounts",
	"tiktok_update_promotion_status": "activity_id",

	// Finance (3)
	"tiktok_search_settlements":             "settlements, next_page_token",
	"tiktok_search_statements":              "statements, next_page_token",
	"tiktok_get_order_statement_transactions": "statement_transactions, order_id",

	// Shop & Authorized (2)
	"tiktok_get_authorized_shops": "shops",
	"tiktok_get_shop_profile":     "shops, warehouses",

	// Analytics (2)
	"tiktok_get_shop_performance": "metrics, gmv, orders, visitors",
	"tiktok_get_bestsellers":      "bestsellers, categories, products",
}

type relatedMeta struct {
	tool  string
	field string
}

var paramRelated = map[string]relatedMeta{
	"product_id":  {tool: "tiktok_search_products", field: "products[0].id"},
	"product_ids": {tool: "tiktok_search_products", field: "products[0].id"},
	"order_id":    {tool: "tiktok_search_orders", field: "orders[0].id"},
	"package_id":  {tool: "tiktok_search_orders", field: "orders[0].packages[0].id"},
	"package_ids": {tool: "tiktok_search_orders", field: "orders[0].packages[0].id"},
	"return_id":   {tool: "tiktok_search_returns", field: "returns[0].id"},
	"activity_id": {tool: "tiktok_search_promotions", field: "activities[0].id"},
	"category_id": {tool: "tiktok_get_categories", field: "categories[0].id"},
}

func annotateRelated(toolName string, props map[string]interface{}) {
	for k, v := range props {
		propMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if rel, found := paramRelated[k]; found {
			propMap["x-datumbridge-related-tool"] = rel.tool
			propMap["x-datumbridge-output-field"] = rel.field
		}
	}
}

func formatToolDescription(name, desc string) string {
	res := strings.TrimSpace(desc)
	if !strings.Contains(strings.ToLower(res), "capabilities:") {
		res += "\n\nCapabilities: tiktok." + name
	}
	if !strings.Contains(strings.ToLower(res), "outputs:") {
		outs := toolOutputs[name]
		if outs != "" {
			res += "\nOutputs: " + outs + ", success"
		} else {
			res += "\nOutputs: success"
		}
	}
	return res
}

// Register returns the full tiktok-mcp tool catalogue across all 24 scopes.
// Full scope is the default; TIKTOK_TOOLS narrows it to an explicit allowlist if set.
func Register() ([]mcp.ToolDesc, map[string]mcp.ToolHandler) {
	cfg := LoadServerConfig()
	handlers := map[string]mcp.ToolHandler{}
	var descs []mcp.ToolDesc

	add := func(name, desc string, props map[string]interface{}, required []string, h mcp.ToolHandler) {
		if !cfg.toolEnabled(name) {
			return
		}
		annotateRelated(name, props)
		input := schema(props, required)
		input["x-datumbridge-capabilities"] = []string{"tiktok." + name}
		if doc := registrydocs.Markdown(name); doc != "" {
			input["x-datumbridge-docs"] = doc
		}
		descs = append(descs, mcp.ToolDesc{
			Name:        name,
			Description: formatToolDescription(name, desc),
			InputSchema: input,
		})
		handlers[name] = h
	}

	registerProductTools(add)
	registerOrderTools(add)
	registerFulfillmentTools(add)
	registerReturnTools(add)
	registerPromotionTools(add)
	registerFinanceTools(add)
	registerShopTools(add)
	registerAnalyticsTools(add)

	return descs, handlers
}
