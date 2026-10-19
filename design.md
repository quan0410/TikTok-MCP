# Design — tiktok-mcp

## Class

**Tool-server** (Streamable HTTP MCP), matching `Bright-Data-MCP`, `shopify-mcp`, and `magento-mcp`'s shape.

## Language

Go 1.23, matching the platform standard established in `doc/MCP_TOOL_REFERENCE_ARCHITECTURE.md` (BR1).

## Upstream

- TikTok Shop Open API (`https://open.tiktokshop.com`, version `2023-09` / `2024-05`)
- `shop_cipher` and `shop_id` come directly from the injected credential bundle (`Credentials.ShopCipher`), never from tool arguments (preventing multi-tenant confusion and IDOR risks).

## Auth & Signing

- Production: Vault `credentials_json` injected per tool call via `POST /internal/v1/inject` on `datumbridge-integrations`.
- Fallback: `TIKTOK_SHOP_ACCESS_TOKEN`, `TIKTOK_SHOP_APP_KEY`, `TIKTOK_SHOP_APP_SECRET`, `TIKTOK_SHOP_CIPHER` env vars (local dev only).
- `credentials_path` under `TIKTOK_CREDENTIALS_DIR` supported as a path-jailed alternative.
- Request Signing: Every request to TikTok Shop Open API includes timestamp, app_key, shop_cipher (when shop-specific), and an HMAC-SHA256 signature (`sign`) calculated across sorted query parameters and raw request body using `app_secret`.
- Refresh token handling: Managed upstream by `datumbridge-integrations`. The tool-server uses the token in `x-tts-access-token` directly.

## Scope Coverage (24 Scopes)

Mapped across 8 tool categories:
1. **Products & Global Products** (`internal/tools/product_tools.go`):
   - `tiktok_search_products`: Search products (`seller.product.basic`, `seller.global_product.info`)
   - `tiktok_get_product`: Retrieve product detail (`seller.product.basic`, `seller.global_product.info`)
   - `tiktok_create_product`: Create product (`seller.product.write`, `seller.global_product.write`)
   - `tiktok_update_product`: Update product info (`seller.product.write`, `seller.global_product.write`)
   - `tiktok_delete_product`: Delete product (`seller.product.delete`, `seller.global_product.delete`)
   - `tiktok_activate_product`: Activate product (`seller.product.delete`)
   - `tiktok_deactivate_product`: Deactivate product (`seller.product.delete`)
   - `tiktok_update_inventory`: Update inventory (`seller.product.write`)
   - `tiktok_update_price`: Update SKU price (`seller.product.write`)
   - `tiktok_get_categories`: Category tree (`seller.global_product.category.info`)
   - `tiktok_get_category_rules`: Category rules & attributes (`seller.global_product.category.info`)
   - `tiktok_diagnose_product`: Product diagnosis (`seller.product.optimize`)
2. **Orders** (`internal/tools/order_tools.go`):
   - `tiktok_search_orders`: Search orders (`seller.order.info`)
   - `tiktok_get_order`: Retrieve order details (`seller.order.info`)
   - `tiktok_cancel_order`: Cancel order (`seller.order.info`)
3. **Fulfillment, Logistics & Delivery** (`internal/tools/fulfillment_tools.go`):
   - `tiktok_get_package_shipping_services`: Look up carriers/services (`seller.logistics`, `seller.fulfillment.basic`)
   - `tiktok_ship_package`: Ship package (`seller.fulfillment.basic`)
   - `tiktok_get_shipping_document`: Shipping label & packing slip URL (`seller.fulfillment.basic`)
   - `tiktok_get_package_tracking`: Tracking trajectory (`seller.logistics`)
   - `tiktok_split_package`: Split order into packages (`seller.fulfillment.package.write`)
   - `tiktok_combine_packages`: Combine packages (`seller.fulfillment.package.write`)
   - `tiktok_update_delivery_status`: Update delivery status (`seller.delivery.status.write`)
4. **Return & Refund** (`internal/tools/return_tools.go`):
   - `tiktok_search_returns`: Search returns/refunds (`seller.return_refund.basic`)
   - `tiktok_get_return`: View return details (`seller.return_refund.basic`)
   - `tiktok_approve_return`: Approve return (`seller.return_refund.basic`)
   - `tiktok_reject_return`: Reject return (`seller.return_refund.basic`)
5. **Promotions** (`internal/tools/promotion_tools.go`):
   - `tiktok_search_promotions`: Search activities (`seller.promotion.info`)
   - `tiktok_get_promotion`: View promotion activity (`seller.promotion.info`)
   - `tiktok_update_promotion_status`: Deactivate promotion (`seller.promotion.write`)
6. **Finance** (`internal/tools/finance_tools.go`):
   - `tiktok_search_settlements`: Search settlements (`seller.finance.info`)
   - `tiktok_search_statements`: Search statements (`seller.finance.info`)
   - `tiktok_get_order_statement_transactions`: Order financial details (`seller.finance.info`)
7. **Shop & Authorized Accounts** (`internal/tools/shop_tools.go`):
   - `tiktok_get_authorized_shops`: Authorized shops (`seller.authorization.info`)
   - `tiktok_get_shop_profile`: Shop profile & warehouses (`seller.shop.info`)
8. **Analytics & Bestsellers** (`internal/tools/analytics_tools.go`):
   - `tiktok_get_shop_performance`: Analytics summary (`seller.data.shop_analytics.public.read`, `creator.data.live.read.public`)
   - `tiktok_get_bestsellers`: Bestselling rankings (`seller.data.bestselling.public.read`)

## Security Baseline

- Never log secrets, app secret, or access tokens.
- Mask tokens in logs and status responses.
- Enforce shop cipher pinning from credential bundle.
- 10 MiB response limit on upstream reads.
