import config from "@/data/config.js";

// April 12, 2026 at 12:12 CEST = 10:12 UTC
const FALLBACK_ACTIVATION = new Date("2026-04-12T10:12:00Z");

async function getGoldenKeySettings() {
    try {
        const response = await fetch(`${config.apiUrl}golden-key`);
        if (!response.ok) throw new Error("Bad response");
        return await response.json();
    } catch {
        return {
            activation_time: FALLBACK_ACTIVATION.toISOString(),
            is_active: new Date() >= FALLBACK_ACTIVATION,
            banner_text: {}
        };
    }
}

async function updateGoldenKeySettings({ activation_time, banner_text = {} }) {
    const token = localStorage.getItem("admin_token");
    const response = await fetch(`${config.apiUrl}admin/golden-key`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json",
            ...(token && { "Authorization": `Bearer ${token}` })
        },
        body: JSON.stringify({ activation_time, banner_text })
    });
    if (!response.ok) throw new Error("Failed to update golden key settings");
    return await response.json();
}

export { getGoldenKeySettings, updateGoldenKeySettings };
