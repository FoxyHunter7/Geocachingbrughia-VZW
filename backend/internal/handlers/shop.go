package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ShopSettings struct {
	StripeSecretKey      string `json:"stripe_secret_key"`
	StripePublishableKey string `json:"stripe_publishable_key"`
	StripeWebhookSecret  string `json:"stripe_webhook_secret"`
	PretixWidgetURL      string `json:"pretix_widget_url"`
	Currency             string `json:"currency"`
}

type ShopItemTranslation struct {
	LangCode    string `json:"lang_code"`
	Description string `json:"description"`
}

type ShopItem struct {
	ID                int64                 `json:"id"`
	Title             string                `json:"title"`
	Description       string                `json:"description"`
	PriceCents        int                   `json:"price_cents"`
	PriceDisplay      string                `json:"price_display"`
	ImageURL          string                `json:"image_url,omitempty"`
	StockQuantity     *int                  `json:"stock_quantity,omitempty"`
	AllowPickup       bool                  `json:"allow_pickup"`
	PickupLabel       string                `json:"pickup_label,omitempty"`
	AllowShipping     bool                  `json:"allow_shipping"`
	ShippingCountries []string              `json:"shipping_countries,omitempty"`
	AutoConfirm       bool                  `json:"auto_confirm"`
	Active            bool                  `json:"active"`
	SortOrder         int                   `json:"sort_order"`
	Translations      []ShopItemTranslation `json:"translations,omitempty"`
}

type ShopOrder struct {
	ID                  int64  `json:"id"`
	ItemID              int64  `json:"item_id"`
	ItemTitle           string `json:"item_title"`
	StripeSessionID     string `json:"stripe_session_id,omitempty"`
	StripePaymentIntent string `json:"stripe_payment_intent_id,omitempty"`
	BuyerEmail          string `json:"buyer_email"`
	Quantity            int    `json:"quantity"`
	AmountCents         int    `json:"amount_cents"`
	AmountDisplay       string `json:"amount_display"`
	FulfillmentType     string `json:"fulfillment_type"`
	ShippingName        string `json:"shipping_name,omitempty"`
	ShippingAddress     string `json:"shipping_address,omitempty"`
	ShippingCity        string `json:"shipping_city,omitempty"`
	ShippingPostalCode  string `json:"shipping_postal_code,omitempty"`
	ShippingCountry     string `json:"shipping_country,omitempty"`
	Status              string `json:"status"`
	Notes               string `json:"notes,omitempty"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

func formatPrice(cents int, currency string) string {
	symbol := "\u20ac"
	if currency == "USD" {
		symbol = "$"
	}
	return fmt.Sprintf("%s %d.%02d", symbol, cents/100, cents%100)
}

func (h *Handler) getShopSettings() (ShopSettings, error) {
	var s ShopSettings
	err := h.db.QueryRow(`SELECT stripe_secret_key, stripe_publishable_key, stripe_webhook_secret, pretix_widget_url, currency FROM shop_settings WHERE id = 1`).
		Scan(&s.StripeSecretKey, &s.StripePublishableKey, &s.StripeWebhookSecret, &s.PretixWidgetURL, &s.Currency)
	if err != nil {
		return s, err
	}
	return s, nil
}

func scanShopItem(rows interface {
	Scan(dest ...any) error
}) (ShopItem, error) {
	var item ShopItem
	var imageURL sql.NullString
	var stockQty sql.NullInt64
	var allowPickup, allowShipping, autoConfirm, active int
	var shippingRegionsJSON string

	err := rows.Scan(
		&item.ID, &item.Title, &item.Description, &item.PriceCents,
		&imageURL, &stockQty, &allowPickup, &item.PickupLabel,
		&allowShipping, &shippingRegionsJSON, &autoConfirm, &active, &item.SortOrder,
	)
	if err != nil {
		return item, err
	}

	if imageURL.Valid {
		item.ImageURL = imageURL.String
	}
	if stockQty.Valid {
		q := int(stockQty.Int64)
		item.StockQuantity = &q
	}
	if shippingRegionsJSON != "" && shippingRegionsJSON != "[]" {
		json.Unmarshal([]byte(shippingRegionsJSON), &item.ShippingCountries)
	}
	item.AllowPickup = allowPickup == 1
	item.AllowShipping = allowShipping == 1
	item.AutoConfirm = autoConfirm == 1
	item.Active = active == 1
	return item, nil
}

func (h *Handler) GetPublicShopSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.getShopSettings()
	if err != nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"stripe_publishable_key": "",
			"pretix_widget_url":      "",
			"currency":               "EUR",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"stripe_publishable_key": s.StripePublishableKey,
		"pretix_widget_url":      s.PretixWidgetURL,
		"currency":               s.Currency,
	})
}

func (h *Handler) GetAdminShopSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.getShopSettings()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch shop settings"})
		return
	}
	respondJSON(w, http.StatusOK, s)
}

func (h *Handler) UpdateShopSettings(w http.ResponseWriter, r *http.Request) {
	var req ShopSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	req.StripeSecretKey = strings.TrimSpace(req.StripeSecretKey)
	req.StripePublishableKey = strings.TrimSpace(req.StripePublishableKey)
	req.StripeWebhookSecret = strings.TrimSpace(req.StripeWebhookSecret)
	req.PretixWidgetURL = truncateString(strings.TrimSpace(req.PretixWidgetURL), 500)

	if len(req.StripeSecretKey) > 200 || len(req.StripePublishableKey) > 200 || len(req.StripeWebhookSecret) > 200 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Stripe key values are too long"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE shop_settings SET
			stripe_secret_key = ?, stripe_publishable_key = ?, stripe_webhook_secret = ?,
			pretix_widget_url = ?, currency = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`, req.StripeSecretKey, req.StripePublishableKey, req.StripeWebhookSecret, req.PretixWidgetURL, "EUR")

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update shop settings"})
		return
	}

	h.GetAdminShopSettings(w, r)
}

func (h *Handler) GetPublicShopItems(w http.ResponseWriter, r *http.Request) {
	settings, _ := h.getShopSettings()
	lang := r.URL.Query().Get("lang")

	rows, err := h.db.Query(`
		SELECT id, title, description, price_cents, image_url, stock_quantity,
		       allow_pickup, pickup_label, allow_shipping, shipping_regions,
		       auto_confirm, active, sort_order
		FROM shop_items WHERE active = 1 ORDER BY sort_order, id
	`)
	if err != nil {
		respondJSON(w, http.StatusOK, []ShopItem{})
		return
	}
	defer rows.Close()

	items := []ShopItem{}
	for rows.Next() {
		item, err := scanShopItem(rows)
		if err != nil {
			continue
		}
		item.PriceDisplay = formatPrice(item.PriceCents, settings.Currency)
		if lang != "" {
			item.Translations = h.getShopItemTranslations(item.ID, lang)
			for _, t := range item.Translations {
				if t.Description != "" {
					item.Description = t.Description
					break
				}
			}
		}
		items = append(items, item)
	}

	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetAdminShopItems(w http.ResponseWriter, r *http.Request) {
	settings, _ := h.getShopSettings()

	rows, err := h.db.Query(`
		SELECT id, title, description, price_cents, image_url, stock_quantity,
		       allow_pickup, pickup_label, allow_shipping, shipping_regions,
		       auto_confirm, active, sort_order
		FROM shop_items ORDER BY sort_order, id
	`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []ShopItem{})
		return
	}
	defer rows.Close()

	items := []ShopItem{}
	for rows.Next() {
		item, err := scanShopItem(rows)
		if err != nil {
			continue
		}
		item.PriceDisplay = formatPrice(item.PriceCents, settings.Currency)
		item.Translations = h.getShopItemTranslations(item.ID, "")
		items = append(items, item)
	}

	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetShopItemByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	settings, _ := h.getShopSettings()

	item, err := scanShopItem(h.db.QueryRow(`
		SELECT id, title, description, price_cents, image_url, stock_quantity,
		       allow_pickup, pickup_label, allow_shipping, shipping_regions,
		       auto_confirm, active, sort_order
		FROM shop_items WHERE id = ?
	`, id))

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Item not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	item.PriceDisplay = formatPrice(item.PriceCents, settings.Currency)
	item.Translations = h.getShopItemTranslations(item.ID, "")
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateShopItem(w http.ResponseWriter, r *http.Request) {
	var item ShopItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	item.Title = strings.TrimSpace(item.Title)
	item.Description = truncateString(item.Description, maxStringLength)
	item.PickupLabel = truncateString(strings.TrimSpace(item.PickupLabel), maxTitleLength)
	item.ImageURL = truncateString(strings.TrimSpace(item.ImageURL), 500)

	// Validate and serialize shipping countries
	var shippingCountriesJSON string
	if item.AllowShipping {
		if len(item.ShippingCountries) == 0 {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Shipping countries are required when shipping is enabled"})
			return
		}
		validCountries := []string{}
		for _, code := range item.ShippingCountries {
			code = strings.ToUpper(strings.TrimSpace(code))
			if validateCountryCode(code) {
				validCountries = append(validCountries, code)
			}
		}
		if len(validCountries) == 0 {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "No valid shipping countries selected"})
			return
		}
		item.ShippingCountries = validCountries
		bytes, _ := json.Marshal(validCountries)
		shippingCountriesJSON = string(bytes)
	}

	if item.Title == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
		return
	}
	if len(item.Title) > maxTitleLength {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is too long (max 200 characters)"})
		return
	}
	if !validatePriceCents(item.PriceCents) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Price must be between 0.01 and 1000000.00"})
		return
	}
	if !validateStock(item.StockQuantity) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Stock must be between 0 and 999999"})
		return
	}
	if !item.AllowPickup && !item.AllowShipping {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "At least one fulfillment option (pickup or shipping) is required"})
		return
	}
	if item.AllowPickup && item.PickupLabel == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Pickup label is required when pickup is enabled"})
		return
	}

	result, err := h.db.Exec(`
		INSERT INTO shop_items (title, description, price_cents, image_url, stock_quantity,
		                        allow_pickup, pickup_label, allow_shipping, shipping_regions,
		                        auto_confirm, active, sort_order)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, item.Title, item.Description, item.PriceCents, item.ImageURL, nullableInt(item.StockQuantity),
		boolToInt(item.AllowPickup), item.PickupLabel, boolToInt(item.AllowShipping), shippingCountriesJSON,
		boolToInt(item.AutoConfirm), boolToInt(item.Active), item.SortOrder)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create item"})
		return
	}

	item.ID, _ = result.LastInsertId()
	h.saveShopItemTranslations(item.ID, item.Translations)
	settings, _ := h.getShopSettings()
	item.PriceDisplay = formatPrice(item.PriceCents, settings.Currency)
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateShopItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var item ShopItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	item.Title = strings.TrimSpace(item.Title)
	item.Description = truncateString(item.Description, maxStringLength)
	item.PickupLabel = truncateString(strings.TrimSpace(item.PickupLabel), maxTitleLength)
	item.ImageURL = truncateString(strings.TrimSpace(item.ImageURL), 500)

	// Validate and serialize shipping countries
	var shippingCountriesJSON string
	if item.AllowShipping {
		if len(item.ShippingCountries) == 0 {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Shipping countries are required when shipping is enabled"})
			return
		}
		validCountries := []string{}
		for _, code := range item.ShippingCountries {
			code = strings.ToUpper(strings.TrimSpace(code))
			if validateCountryCode(code) {
				validCountries = append(validCountries, code)
			}
		}
		if len(validCountries) == 0 {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "No valid shipping countries selected"})
			return
		}
		item.ShippingCountries = validCountries
		bytes, _ := json.Marshal(validCountries)
		shippingCountriesJSON = string(bytes)
	}

	if item.Title == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
		return
	}
	if len(item.Title) > maxTitleLength {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is too long (max 200 characters)"})
		return
	}
	if !validatePriceCents(item.PriceCents) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Price must be between 0.01 and 1000000.00"})
		return
	}
	if !validateStock(item.StockQuantity) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Stock must be between 0 and 999999"})
		return
	}
	if !item.AllowPickup && !item.AllowShipping {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "At least one fulfillment option (pickup or shipping) is required"})
		return
	}
	if item.AllowPickup && item.PickupLabel == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Pickup label is required when pickup is enabled"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE shop_items SET
			title = ?, description = ?, price_cents = ?, image_url = ?, stock_quantity = ?,
			allow_pickup = ?, pickup_label = ?, allow_shipping = ?, shipping_regions = ?,
			auto_confirm = ?, active = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, item.Title, item.Description, item.PriceCents, item.ImageURL, nullableInt(item.StockQuantity),
		boolToInt(item.AllowPickup), item.PickupLabel, boolToInt(item.AllowShipping), shippingCountriesJSON,
		boolToInt(item.AutoConfirm), boolToInt(item.Active), item.SortOrder, id)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update item"})
		return
	}

	idInt, _ := strconv.ParseInt(id, 10, 64)
	item.ID = idInt
	h.saveShopItemTranslations(item.ID, item.Translations)
	settings, _ := h.getShopSettings()
	item.PriceDisplay = formatPrice(item.PriceCents, settings.Currency)
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteShopItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.db.Exec("DELETE FROM shop_items WHERE id = ?", id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete item"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Item deleted"})
}

func (h *Handler) GetAdminShopOrders(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	settings, _ := h.getShopSettings()

	var rows *sql.Rows
	var err error

	if status != "" {
		rows, err = h.db.Query(`
			SELECT o.id, o.item_id, COALESCE(i.title, '(deleted)'), o.stripe_session_id,
			       o.stripe_payment_intent_id, o.buyer_email, o.quantity, o.amount_cents,
			       o.fulfillment_type, o.shipping_name, o.shipping_address, o.shipping_city,
			       o.shipping_postal_code, o.shipping_country, o.status, o.notes,
			       o.created_at, o.updated_at
			FROM shop_orders o
			LEFT JOIN shop_items i ON o.item_id = i.id
			WHERE o.status = ?
			ORDER BY o.created_at DESC
		`, status)
	} else {
		rows, err = h.db.Query(`
			SELECT o.id, o.item_id, COALESCE(i.title, '(deleted)'), o.stripe_session_id,
			       o.stripe_payment_intent_id, o.buyer_email, o.quantity, o.amount_cents,
			       o.fulfillment_type, o.shipping_name, o.shipping_address, o.shipping_city,
			       o.shipping_postal_code, o.shipping_country, o.status, o.notes,
			       o.created_at, o.updated_at
			FROM shop_orders o
			LEFT JOIN shop_items i ON o.item_id = i.id
			ORDER BY
				CASE o.status
					WHEN 'pending' THEN 1
					WHEN 'paid' THEN 2
					WHEN 'confirmed' THEN 3
					WHEN 'shipped' THEN 4
					WHEN 'fulfilled' THEN 5
					ELSE 6
				END,
				o.created_at DESC
		`)
	}

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []ShopOrder{})
		return
	}
	defer rows.Close()

	orders := []ShopOrder{}
	for rows.Next() {
		var o ShopOrder
		if err := rows.Scan(
			&o.ID, &o.ItemID, &o.ItemTitle, &o.StripeSessionID, &o.StripePaymentIntent,
			&o.BuyerEmail, &o.Quantity, &o.AmountCents, &o.FulfillmentType,
			&o.ShippingName, &o.ShippingAddress, &o.ShippingCity, &o.ShippingPostalCode,
			&o.ShippingCountry, &o.Status, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			continue
		}
		o.AmountDisplay = formatPrice(o.AmountCents, settings.Currency)
		orders = append(orders, o)
	}

	respondJSON(w, http.StatusOK, orders)
}

func (h *Handler) GetShopOrderByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	settings, _ := h.getShopSettings()

	var o ShopOrder
	err := h.db.QueryRow(`
		SELECT o.id, o.item_id, COALESCE(i.title, '(deleted)'), o.stripe_session_id,
		       o.stripe_payment_intent_id, o.buyer_email, o.quantity, o.amount_cents,
		       o.fulfillment_type, o.shipping_name, o.shipping_address, o.shipping_city,
		       o.shipping_postal_code, o.shipping_country, o.status, o.notes,
		       o.created_at, o.updated_at
		FROM shop_orders o
		LEFT JOIN shop_items i ON o.item_id = i.id
		WHERE o.id = ?
	`, id).Scan(
		&o.ID, &o.ItemID, &o.ItemTitle, &o.StripeSessionID, &o.StripePaymentIntent,
		&o.BuyerEmail, &o.Quantity, &o.AmountCents, &o.FulfillmentType,
		&o.ShippingName, &o.ShippingAddress, &o.ShippingCity, &o.ShippingPostalCode,
		&o.ShippingCountry, &o.Status, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Order not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	o.AmountDisplay = formatPrice(o.AmountCents, settings.Currency)
	respondJSON(w, http.StatusOK, o)
}

func (h *Handler) UpdateShopOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var update struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	update.Status = strings.TrimSpace(update.Status)
	update.Notes = truncateString(strings.TrimSpace(update.Notes), maxStringLength)

	validStatuses := map[string]bool{
		"pending": true, "paid": true, "confirmed": true,
		"shipped": true, "fulfilled": true, "cancelled": true,
	}
	if !validStatuses[update.Status] {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid status"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE shop_orders SET status = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, update.Status, update.Notes, id)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update order"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Order updated"})
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullableInt(v *int) interface{} {
	if v == nil {
		return nil
	}
	return *v
}

func (h *Handler) getShopItemTranslations(itemID int64, langFilter string) []ShopItemTranslation {
	var rows *sql.Rows
	var err error

	if langFilter != "" {
		rows, err = h.db.Query("SELECT lang_code, description FROM shop_item_translations WHERE item_id = ? AND lang_code = ?", itemID, langFilter)
	} else {
		rows, err = h.db.Query("SELECT lang_code, description FROM shop_item_translations WHERE item_id = ?", itemID)
	}

	if err != nil {
		return []ShopItemTranslation{}
	}
	defer rows.Close()

	translations := []ShopItemTranslation{}
	for rows.Next() {
		var t ShopItemTranslation
		if err := rows.Scan(&t.LangCode, &t.Description); err != nil {
			continue
		}
		translations = append(translations, t)
	}
	return translations
}

func (h *Handler) saveShopItemTranslations(itemID int64, translations []ShopItemTranslation) {
	for _, t := range translations {
		desc := truncateString(strings.TrimSpace(t.Description), maxStringLength)
		h.db.Exec(`INSERT OR REPLACE INTO shop_item_translations (item_id, lang_code, description) VALUES (?, ?, ?)`,
			itemID, t.LangCode, desc)
	}
}
