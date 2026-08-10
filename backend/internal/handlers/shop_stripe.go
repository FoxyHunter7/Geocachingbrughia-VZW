package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type createCheckoutRequest struct {
	ItemID          int64  `json:"item_id"`
	Quantity        int    `json:"quantity"`
	FulfillmentType string `json:"fulfillment_type"`
	BuyerEmail      string `json:"buyer_email"`
	ShippingName    string `json:"shipping_name"`
	ShippingAddress string `json:"shipping_address"`
	ShippingCity    string `json:"shipping_city"`
	ShippingPostal  string `json:"shipping_postal_code"`
	ShippingCountry string `json:"shipping_country"`
}

func (h *Handler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var req createCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.ItemID == 0 || !validateQuantity(req.Quantity) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid item or quantity (max 999)"})
		return
	}

	if req.FulfillmentType != "pickup" && req.FulfillmentType != "shipping" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid fulfillment type"})
		return
	}

	if !validateEmail(req.BuyerEmail) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "A valid email is required"})
		return
	}

	if req.FulfillmentType == "shipping" {
		if !validateCountryCode(req.ShippingCountry) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid or unsupported shipping country"})
			return
		}
		req.ShippingName = truncateString(strings.TrimSpace(req.ShippingName), 200)
		req.ShippingAddress = truncateString(strings.TrimSpace(req.ShippingAddress), 500)
		req.ShippingCity = truncateString(strings.TrimSpace(req.ShippingCity), 100)
		req.ShippingPostal = truncateString(strings.TrimSpace(req.ShippingPostal), 20)
		if req.ShippingName == "" || req.ShippingAddress == "" || req.ShippingCity == "" || req.ShippingPostal == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "All shipping fields are required"})
			return
		}
	}

	settings, err := h.getShopSettings()
	if err != nil || settings.StripeSecretKey == "" {
		respondJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Shop is not configured"})
		return
	}

	var item ShopItem
	var imageURL sql.NullString
	var stockQty sql.NullInt64
	var allowPickup, allowShipping, autoConfirm, active int

	err = h.db.QueryRow(`
		SELECT id, title, description, price_cents, image_url, stock_quantity,
		       allow_pickup, pickup_label, allow_shipping, shipping_regions,
		       auto_confirm, active, sort_order
		FROM shop_items WHERE id = ? AND active = 1
	`, req.ItemID).Scan(
		&item.ID, &item.Title, &item.Description, &item.PriceCents,
		&imageURL, &stockQty, &allowPickup, &item.PickupLabel,
		&allowShipping, &item.ShippingRegions, &autoConfirm, &active, &item.SortOrder,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Item not found or inactive"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if req.FulfillmentType == "pickup" && allowPickup != 1 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Pickup not available for this item"})
		return
	}
	if req.FulfillmentType == "shipping" && allowShipping != 1 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Shipping not available for this item"})
		return
	}

	if stockQty.Valid {
		remaining := int(stockQty.Int64)
		if remaining < req.Quantity {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Insufficient stock"})
			return
		}
	}

	totalCents := item.PriceCents * req.Quantity
	if totalCents > maxPriceCents || totalCents < item.PriceCents {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Total amount exceeds maximum allowed"})
		return
	}

	result, err := h.db.Exec(`
		INSERT INTO shop_orders (item_id, buyer_email, quantity, amount_cents,
		                         fulfillment_type, shipping_name, shipping_address,
		                         shipping_city, shipping_postal_code, shipping_country,
		                         status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')
	`, req.ItemID, req.BuyerEmail, req.Quantity, totalCents,
		req.FulfillmentType, req.ShippingName, req.ShippingAddress,
		req.ShippingCity, req.ShippingPostal, req.ShippingCountry)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create order"})
		return
	}

	orderID, _ := result.LastInsertId()

	frontendURL := h.cfg.FrontendURL
	successURL := fmt.Sprintf("%s/shop/success?order=%d", frontendURL, orderID)
	cancelURL := fmt.Sprintf("%s/shop/cancel?order=%d", frontendURL, orderID)

	itemName := item.Title
	if req.Quantity > 1 {
		itemName = fmt.Sprintf("%s (x%d)", item.Title, req.Quantity)
	}

	formData := url.Values{}
	formData.Set("mode", "payment")
	formData.Set("success_url", successURL)
	formData.Set("cancel_url", cancelURL)
	formData.Set("customer_email", req.BuyerEmail)
	formData.Set("line_items[0][quantity]", strconv.Itoa(req.Quantity))
	formData.Set("line_items[0][price_data][currency]", strings.ToLower(settings.Currency))
	formData.Set("line_items[0][price_data][unit_amount]", strconv.Itoa(item.PriceCents))
	formData.Set("line_items[0][price_data][product_data][name]", itemName)
	if item.Description != "" {
		formData.Set("line_items[0][price_data][product_data][description]", item.Description)
	}
	if imageURL.Valid && imageURL.String != "" {
		formData.Set("line_items[0][price_data][product_data][images][0]", imageURL.String)
	}

	formData.Set("metadata[order_id]", strconv.FormatInt(orderID, 10))
	formData.Set("metadata[item_id]", strconv.FormatInt(req.ItemID, 10))
	formData.Set("metadata[fulfillment_type]", req.FulfillmentType)

	stripeReq, err := http.NewRequest("POST", "https://api.stripe.com/v1/checkout/sessions", strings.NewReader(formData.Encode()))
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create checkout request"})
		return
	}
	stripeReq.SetBasicAuth(settings.StripeSecretKey, "")
	stripeReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	stripeResp, err := client.Do(stripeReq)
	if err != nil {
		h.db.Exec("UPDATE shop_orders SET status = 'cancelled' WHERE id = ?", orderID)
		respondJSON(w, http.StatusBadGateway, map[string]string{"error": "Failed to connect to Stripe"})
		return
	}
	defer stripeResp.Body.Close()

	body, _ := io.ReadAll(stripeResp.Body)

	if stripeResp.StatusCode != http.StatusOK {
		var stripeErr struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		json.Unmarshal(body, &stripeErr)
		msg := stripeErr.Error.Message
		if msg == "" {
			msg = "Stripe checkout creation failed"
		}
		h.db.Exec("UPDATE shop_orders SET status = 'cancelled' WHERE id = ?", orderID)
		respondJSON(w, http.StatusBadGateway, map[string]string{"error": msg})
		return
	}

	var session struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	json.Unmarshal(body, &session)

	h.db.Exec("UPDATE shop_orders SET stripe_session_id = ? WHERE id = ?", session.ID, orderID)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"checkout_url": session.URL,
		"session_id":   session.ID,
		"order_id":     orderID,
	})
}

func (h *Handler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	settings, err := h.getShopSettings()
	if err != nil || settings.StripeWebhookSecret == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sig := r.Header.Get("Stripe-Signature")
	if sig == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !verifyStripeSignature(rawBody, sig, settings.StripeWebhookSecret) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var event struct {
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID            string `json:"id"`
				PaymentIntent string `json:"payment_intent"`
				Metadata      struct {
					OrderID string `json:"order_id"`
				} `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}

	if err := json.Unmarshal(rawBody, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		orderIDStr := event.Data.Object.Metadata.OrderID
		if orderIDStr == "" {
			w.WriteHeader(http.StatusOK)
			return
		}

		orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		var itemID int64
		var autoConfirm int
		var qty int
		err = h.db.QueryRow(`
			SELECT o.item_id, COALESCE(i.auto_confirm, 0), o.quantity
			FROM shop_orders o
			LEFT JOIN shop_items i ON o.item_id = i.id
			WHERE o.id = ?
		`, orderID).Scan(&itemID, &autoConfirm, &qty)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		newStatus := "paid"
		if autoConfirm == 1 {
			newStatus = "confirmed"
		}

		_, err = h.db.Exec(`
			UPDATE shop_orders SET
				status = ?, stripe_payment_intent_id = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ? AND status = 'pending'
		`, newStatus, event.Data.Object.PaymentIntent, orderID)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		if autoConfirm == 1 {
			h.db.Exec(`
				UPDATE shop_items SET stock_quantity = stock_quantity - ?
				WHERE id = ? AND stock_quantity IS NOT NULL
			`, qty, itemID)
		}

	case "payment_intent.payment_failed":
		orderIDStr := event.Data.Object.Metadata.OrderID
		if orderIDStr != "" {
			if orderID, err := strconv.ParseInt(orderIDStr, 10, 64); err == nil {
				h.db.Exec("UPDATE shop_orders SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP WHERE id = ? AND status = 'pending'", orderID)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func verifyStripeSignature(payload []byte, header, secret string) bool {
	var timestamp string
	var signatures []string

	parts := strings.Split(header, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			timestamp = kv[1]
		case "v1":
			signatures = append(signatures, kv[1])
		}
	}

	if timestamp == "" || len(signatures) == 0 {
		return false
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	age := time.Since(time.Unix(ts, 0))
	if age > 5*time.Minute || age < -5*time.Minute {
		return false
	}

	signedPayload := fmt.Sprintf("%d.%s", ts, string(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	for _, sig := range signatures {
		if hmac.Equal([]byte(expectedSig), []byte(sig)) {
			return true
		}
	}

	return false
}
