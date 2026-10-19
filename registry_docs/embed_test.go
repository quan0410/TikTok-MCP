package registrydocs

import (
	"strings"
	"testing"
)

var allTikTokTools = []string{
	// Products (12)
	"tiktok_search_products", "tiktok_get_product", "tiktok_create_product",
	"tiktok_update_product", "tiktok_delete_product", "tiktok_activate_product",
	"tiktok_deactivate_product", "tiktok_update_inventory", "tiktok_update_price",
	"tiktok_get_categories", "tiktok_get_category_rules", "tiktok_diagnose_product",
	// Orders (3)
	"tiktok_search_orders", "tiktok_get_order", "tiktok_cancel_order",
	// Fulfillment & Logistics (7)
	"tiktok_get_package_shipping_services", "tiktok_ship_package", "tiktok_get_shipping_document",
	"tiktok_get_package_tracking", "tiktok_split_package", "tiktok_combine_packages",
	"tiktok_update_delivery_status",
	// Returns & Refunds (4)
	"tiktok_search_returns", "tiktok_get_return", "tiktok_approve_return",
	"tiktok_reject_return",
	// Promotions (3)
	"tiktok_search_promotions", "tiktok_get_promotion", "tiktok_update_promotion_status",
	// Finance (3)
	"tiktok_search_settlements", "tiktok_search_statements", "tiktok_get_order_statement_transactions",
	// Shop (2)
	"tiktok_get_authorized_shops", "tiktok_get_shop_profile",
	// Analytics (2)
	"tiktok_get_shop_performance", "tiktok_get_bestsellers",
}

func TestMarkdown_AllToolsExist(t *testing.T) {
	for _, name := range allTikTokTools {
		doc := Markdown(name)
		if strings.TrimSpace(doc) == "" {
			t.Errorf("missing or empty documentation file for tool: %s", name)
		}
		if !strings.HasPrefix(doc, "# "+name) {
			t.Errorf("doc for %s should start with header '# %s'", name, name)
		}
		for _, section := range []string{"## Agent rules", "## Parameters", "## Success fields", "## Error codes", "## Cases"} {
			if !strings.Contains(doc, section) {
				t.Errorf("doc for %q missing section %q", name, section)
			}
		}
		if !strings.Contains(doc, "Related tool") {
			t.Errorf("doc for %q missing 'Related tool' in Parameters table", name)
		}
	}
}

func TestMarkdown_Sanitization(t *testing.T) {
	if got := Markdown("../etc/passwd"); got != "" {
		t.Errorf("expected empty string for path traversal, got %q", got)
	}
	if got := Markdown("non_existent_tool_xyz"); got != "" {
		t.Errorf("expected empty string for non-existent tool, got %q", got)
	}
}
