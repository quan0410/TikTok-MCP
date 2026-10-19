package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/datumbridge/tiktok-mcp/internal/mcp"
)

func registerProductTools(add toolAdder) {
	// 1. tiktok_search_products
	add("tiktok_search_products",
		"Search and list products on TikTok Shop with optional filters (status: ACTIVATE, SELLER_DEACTIVATED, PLATFORM_DEACTIVATED, FREEZE; page_size, page_token).",
		baseProps(map[string]interface{}{
			"status": map[string]interface{}{
				"type":        "string",
				"description": "Product status filter: ACTIVATE, SELLER_DEACTIVATED, PLATFORM_DEACTIVATED, FREEZE, DRAFT, PENDING",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of products per page (1-100, default 20)",
			},
			"page_token": map[string]interface{}{
				"type":        "string",
				"description": "Pagination token from previous search result",
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
			if status := strArg(args, "status"); status != "" {
				reqBody["status"] = status
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodPost, "/product/202309/products/search", q, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 2. tiktok_get_product
	add("tiktok_get_product",
		"Retrieve full details of a specific product by product_id (title, description, category, SKUs, attributes, images).",
		baseProps(map[string]interface{}{
			"product_id": map[string]interface{}{
				"type":        "string",
				"description": "The unique TikTok Shop product ID",
			},
		}),
		[]string{"product_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			productID := strArg(args, "product_id")
			if productID == "" {
				return mcp.ToolResultError("product_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/product/202309/products/%s", url.PathEscape(productID))
			data, _, err := client.Request(c, http.MethodGet, path, url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 3. tiktok_create_product
	add("tiktok_create_product",
		"Create a new product on TikTok Shop including title, description, category_id, brand_id, images, and SKUs.",
		baseProps(map[string]interface{}{
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Product title",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Product description in HTML or rich text",
			},
			"category_id": map[string]interface{}{
				"type":        "string",
				"description": "TikTok Shop leaf category ID",
			},
			"brand_id": map[string]interface{}{
				"type":        "string",
				"description": "Optional brand ID",
			},
			"main_images": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Array of main image URIs",
			},
			"skus": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "object"},
				"description": "List of SKU objects with price, inventory, and sales attributes",
			},
			"package_weight": map[string]interface{}{
				"type":        "object",
				"description": "Package weight object {value, unit}",
			},
			"package_dimensions": map[string]interface{}{
				"type":        "object",
				"description": "Package dimensions object {length, width, height, unit}",
			},
		}),
		[]string{"title", "category_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			reqBody := map[string]interface{}{
				"title":       strArg(args, "title"),
				"category_id": strArg(args, "category_id"),
			}
			if desc := strArg(args, "description"); desc != "" {
				reqBody["description"] = desc
			}
			if brandID := strArg(args, "brand_id"); brandID != "" {
				reqBody["brand_id"] = brandID
			}
			if images := sliceStrArg(args, "main_images"); len(images) > 0 {
				var imgList []map[string]string
				for _, img := range images {
					imgList = append(imgList, map[string]string{"uri": img})
				}
				reqBody["main_images"] = imgList
			}
			if skus := sliceMapArg(args, "skus"); len(skus) > 0 {
				reqBody["skus"] = skus
			}
			if pw := mapArg(args, "package_weight"); pw != nil {
				reqBody["package_weight"] = pw
			}
			if pd := mapArg(args, "package_dimensions"); pd != nil {
				reqBody["package_dimensions"] = pd
			}
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodPost, "/product/202309/products", url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 4. tiktok_update_product
	add("tiktok_update_product",
		"Update an existing product on TikTok Shop by product_id.",
		baseProps(map[string]interface{}{
			"product_id": map[string]interface{}{
				"type":        "string",
				"description": "The unique product ID to update",
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Updated product title",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Updated product description",
			},
			"category_id": map[string]interface{}{
				"type":        "string",
				"description": "Updated category ID",
			},
			"skus": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "object"},
				"description": "Updated SKU definitions",
			},
		}),
		[]string{"product_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			productID := strArg(args, "product_id")
			if productID == "" {
				return mcp.ToolResultError("product_id is required")
			}
			reqBody := map[string]interface{}{}
			if title := strArg(args, "title"); title != "" {
				reqBody["title"] = title
			}
			if desc := strArg(args, "description"); desc != "" {
				reqBody["description"] = desc
			}
			if catID := strArg(args, "category_id"); catID != "" {
				reqBody["category_id"] = catID
			}
			if skus := sliceMapArg(args, "skus"); len(skus) > 0 {
				reqBody["skus"] = skus
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/product/202309/products/%s", url.PathEscape(productID))
			data, _, err := client.Request(c, http.MethodPut, path, url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 5. tiktok_delete_product
	add("tiktok_delete_product",
		"Permanently delete products from TikTok Shop by product_ids.",
		baseProps(map[string]interface{}{
			"product_ids": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "List of product IDs to delete",
			},
		}),
		[]string{"product_ids"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			productIDs := sliceStrArg(args, "product_ids")
			if len(productIDs) == 0 {
				return mcp.ToolResultError("product_ids must contain at least one ID")
			}
			c, cancel := ctx()
			defer cancel()
			reqBody := map[string]interface{}{
				"product_ids": productIDs,
			}
			data, _, err := client.Request(c, http.MethodPost, "/product/202309/products/delete", url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 6. tiktok_activate_product
	add("tiktok_activate_product",
		"Activate deactivated products on TikTok Shop by product_ids.",
		baseProps(map[string]interface{}{
			"product_ids": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "List of product IDs to activate",
			},
		}),
		[]string{"product_ids"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			productIDs := sliceStrArg(args, "product_ids")
			if len(productIDs) == 0 {
				return mcp.ToolResultError("product_ids required")
			}
			c, cancel := ctx()
			defer cancel()
			reqBody := map[string]interface{}{
				"product_ids": productIDs,
			}
			data, _, err := client.Request(c, http.MethodPost, "/product/202309/products/activate", url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 7. tiktok_deactivate_product
	add("tiktok_deactivate_product",
		"Deactivate active products on TikTok Shop by product_ids.",
		baseProps(map[string]interface{}{
			"product_ids": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "List of product IDs to deactivate",
			},
		}),
		[]string{"product_ids"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			productIDs := sliceStrArg(args, "product_ids")
			if len(productIDs) == 0 {
				return mcp.ToolResultError("product_ids required")
			}
			c, cancel := ctx()
			defer cancel()
			reqBody := map[string]interface{}{
				"product_ids": productIDs,
			}
			data, _, err := client.Request(c, http.MethodPost, "/product/202309/products/deactivate", url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 8. tiktok_update_inventory
	add("tiktok_update_inventory",
		"Update stock quantity / inventory for product SKUs.",
		baseProps(map[string]interface{}{
			"product_id": map[string]interface{}{
				"type":        "string",
				"description": "Product ID",
			},
			"skus": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "object"},
				"description": "List of SKU stock updates [{id: 'sku_id', inventory: [{warehouse_id: '...', quantity: 10}]}]",
			},
		}),
		[]string{"product_id", "skus"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			productID := strArg(args, "product_id")
			skus := sliceMapArg(args, "skus")
			if productID == "" || len(skus) == 0 {
				return mcp.ToolResultError("product_id and skus are required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/product/202309/products/%s/inventory/update", url.PathEscape(productID))
			reqBody := map[string]interface{}{
				"skus": skus,
			}
			data, _, err := client.Request(c, http.MethodPost, path, url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 9. tiktok_update_price
	add("tiktok_update_price",
		"Update pricing for product SKUs.",
		baseProps(map[string]interface{}{
			"product_id": map[string]interface{}{
				"type":        "string",
				"description": "Product ID",
			},
			"skus": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "object"},
				"description": "List of SKU price updates [{id: 'sku_id', price: {amount: '19.99', currency: 'USD'}}]",
			},
		}),
		[]string{"product_id", "skus"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			productID := strArg(args, "product_id")
			skus := sliceMapArg(args, "skus")
			if productID == "" || len(skus) == 0 {
				return mcp.ToolResultError("product_id and skus are required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/product/202309/products/%s/prices/update", url.PathEscape(productID))
			reqBody := map[string]interface{}{
				"skus": skus,
			}
			data, _, err := client.Request(c, http.MethodPost, path, url.Values{}, reqBody)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 10. tiktok_get_categories
	add("tiktok_get_categories",
		"List product categories available on TikTok Shop.",
		baseProps(map[string]interface{}{
			"category_id": map[string]interface{}{
				"type":        "string",
				"description": "Optional parent category ID to query subcategories",
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
			c, cancel := ctx()
			defer cancel()
			data, _, err := client.Request(c, http.MethodGet, "/product/202309/categories", q, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 11. tiktok_get_category_rules
	add("tiktok_get_category_rules",
		"Get publication rules, required attributes, and size charts for a specific category ID.",
		baseProps(map[string]interface{}{
			"category_id": map[string]interface{}{
				"type":        "string",
				"description": "Leaf category ID",
			},
		}),
		[]string{"category_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			categoryID := strArg(args, "category_id")
			if categoryID == "" {
				return mcp.ToolResultError("category_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/product/202309/categories/%s/rules", url.PathEscape(categoryID))
			data, _, err := client.Request(c, http.MethodGet, path, url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)

	// 12. tiktok_diagnose_product
	add("tiktok_diagnose_product",
		"Retrieve optimization diagnostics and suggestions for a product (scope: seller.product.optimize).",
		baseProps(map[string]interface{}{
			"product_id": map[string]interface{}{
				"type":        "string",
				"description": "Product ID to diagnose",
			},
		}),
		[]string{"product_id"},
		func(raw json.RawMessage) map[string]interface{} {
			client, args, err := clientFrom(raw)
			if err != nil {
				return mcp.ToolResultError("credentials error: " + err.Error())
			}
			productID := strArg(args, "product_id")
			if productID == "" {
				return mcp.ToolResultError("product_id is required")
			}
			c, cancel := ctx()
			defer cancel()
			path := fmt.Sprintf("/product/202309/products/%s/diagnosis", url.PathEscape(productID))
			data, _, err := client.Request(c, http.MethodGet, path, url.Values{}, nil)
			if err != nil {
				return mcp.ToolResultError(err.Error())
			}
			return rawResult(data)
		},
	)
}
