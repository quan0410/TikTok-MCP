package tools_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/datumbridge/tiktok-mcp/internal/tools"
)

func clearToolFilterEnv(t *testing.T) {
	t.Helper()
	_ = os.Unsetenv("TIKTOK_TOOLS")
}

func TestRegisterFullScopeByDefault(t *testing.T) {
	clearToolFilterEnv(t)
	descs, handlers := tools.Register()
	if len(descs) != 36 {
		t.Fatalf("expected full tool scope (36 tools), got %d", len(descs))
	}

	wantedTools := []string{
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

	for _, want := range wantedTools {
		if _, ok := handlers[want]; !ok {
			t.Fatalf("missing tool handler: %s", want)
		}
	}

	// Verify schema properties
	for _, d := range descs {
		if d.Name == "" {
			t.Fatalf("found descriptor with empty name")
		}
		if d.Description == "" {
			t.Errorf("tool %s has empty description", d.Name)
		}
		props, ok := d.InputSchema["properties"].(map[string]interface{})
		if !ok {
			t.Errorf("tool %s missing inputSchema.properties", d.Name)
			continue
		}
		if _, ok := props["credentials_json"]; !ok {
			t.Errorf("tool %s missing credentials_json in properties", d.Name)
		}
		if _, ok := props["credentials_path"]; !ok {
			t.Errorf("tool %s missing credentials_path in properties", d.Name)
		}
		caps, ok := d.InputSchema["x-datumbridge-capabilities"].([]string)
		if !ok || len(caps) == 0 {
			t.Errorf("tool %s missing x-datumbridge-capabilities", d.Name)
		}
		docs, ok := d.InputSchema["x-datumbridge-docs"].(string)
		if !ok || strings.TrimSpace(docs) == "" {
			t.Errorf("tool %s missing x-datumbridge-docs", d.Name)
		}
	}
}

func TestToolAllowlist(t *testing.T) {
	os.Setenv("TIKTOK_TOOLS", "tiktok_search_products,tiktok_get_product")
	defer clearToolFilterEnv(t)

	descs, handlers := tools.Register()
	if len(descs) != 2 {
		t.Fatalf("expected 2 tools with allowlist, got %d", len(descs))
	}
	if _, ok := handlers["tiktok_search_products"]; !ok {
		t.Errorf("expected tiktok_search_products to be enabled")
	}
	if _, ok := handlers["tiktok_get_product"]; !ok {
		t.Errorf("expected tiktok_get_product to be enabled")
	}
	if _, ok := handlers["tiktok_create_product"]; ok {
		t.Errorf("expected tiktok_create_product to be excluded")
	}
}

func TestToolExecution_MockServer(t *testing.T) {
	clearToolFilterEnv(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/product/202309/products/search" {
			_, _ = w.Write([]byte(`{"code": 0, "message": "Success", "data": {"products": [{"id": "123", "title": "Test Product"}]}}`))
			return
		}
		if r.URL.Path == "/product/202309/products/123" {
			_, _ = w.Write([]byte(`{"code": 0, "message": "Success", "data": {"product": {"id": "123", "title": "Test Product"}}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code": 404, "message": "not found"}`))
	}))
	defer ts.Close()

	os.Setenv("TIKTOK_API_BASE", ts.URL)
	defer os.Unsetenv("TIKTOK_API_BASE")

	_, handlers := tools.Register()

	credsJSON, _ := json.Marshal(map[string]interface{}{
		"access_token": "valid_token",
		"app_key":      "test_key",
		"app_secret":   "test_secret",
		"shop_cipher":  "test_cipher",
	})

	// Test tiktok_search_products
	searchHandler := handlers["tiktok_search_products"]
	inputArgs, _ := json.Marshal(map[string]interface{}{
		"credentials_json": string(credsJSON),
		"status":           "ACTIVATE",
	})
	res := searchHandler(inputArgs)
	if isErr, _ := res["isError"].(bool); isErr {
		t.Fatalf("search tool returned error: %v", res)
	}

	// Test tiktok_get_product
	getHandler := handlers["tiktok_get_product"]
	inputArgsGet, _ := json.Marshal(map[string]interface{}{
		"credentials_json": string(credsJSON),
		"product_id":       "123",
	})
	resGet := getHandler(inputArgsGet)
	if isErr, _ := resGet["isError"].(bool); isErr {
		t.Fatalf("get tool returned error: %v", resGet)
	}
}
