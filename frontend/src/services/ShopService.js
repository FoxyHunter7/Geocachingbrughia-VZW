import { fetchFromServer, fetchToServer, deleteFromServer } from "./fetcher";

// Public shop
export async function getShopSettings() {
    return fetchFromServer("shop/settings");
}

export async function getShopItems() {
    return fetchFromServer("shop/items");
}

export async function getShopItem(id) {
    return fetchFromServer(`shop/items/${id}`);
}

export async function createCheckoutSession(data) {
    return fetchToServer("shop/checkout", "POST", JSON.stringify(data), false);
}

// Admin shop settings
export async function getAdminShopSettings() {
    return fetchFromServer("admin/shop/settings", true);
}

export async function updateShopSettings(data) {
    return fetchToServer("admin/shop/settings", "PUT", JSON.stringify(data), true);
}

// Admin shop items
export async function getAdminShopItems() {
    return fetchFromServer("admin/shop/items", true);
}

export async function getAdminShopItem(id) {
    return fetchFromServer(`admin/shop/items/${id}`, true);
}

export async function createShopItem(data) {
    return fetchToServer("admin/shop/items", "POST", JSON.stringify(data), true);
}

export async function updateShopItem(id, data) {
    return fetchToServer(`admin/shop/items/${id}`, "PUT", JSON.stringify(data), true);
}

export async function deleteShopItem(id) {
    return deleteFromServer(`admin/shop/items/${id}`);
}

// Admin shop orders
export async function getAdminShopOrders(status) {
    const endpoint = status ? `admin/shop/orders?status=${status}` : "admin/shop/orders";
    return fetchFromServer(endpoint, true);
}

export async function getAdminShopOrder(id) {
    return fetchFromServer(`admin/shop/orders/${id}`, true);
}

export async function updateShopOrderStatus(id, status, notes) {
    return fetchToServer(`admin/shop/orders/${id}/status`, "PUT", JSON.stringify({ status, notes }), true);
}
