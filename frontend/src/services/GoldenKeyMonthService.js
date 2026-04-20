import { fetchFromServer, fetchToServer, deleteFromServer } from "./fetcher";

// Public
export async function getGoldenKeyMonths() {
    return fetchFromServer("golden-key/months");
}

export async function getGoldenKeyMonth(id) {
    return fetchFromServer(`golden-key/months/${id}`);
}

// Admin
export async function getAdminGoldenKeyMonths() {
    return fetchFromServer("admin/golden-key/months", true);
}

export async function getAdminGoldenKeyMonth(id) {
    return fetchFromServer(`admin/golden-key/months/${id}`, true);
}

export async function updateGoldenKeyMonth(id, data) {
    return fetchToServer(`admin/golden-key/months/${id}`, "PUT", JSON.stringify(data), true);
}

// Hints
export async function addGoldenKeyHint(monthId, { content, image_url }) {
    return fetchToServer(
        `admin/golden-key/months/${monthId}/hints`,
        "POST",
        JSON.stringify({ content, image_url: image_url || "" }),
        true
    );
}

export async function updateGoldenKeyHint(hintId, { content, image_url }) {
    return fetchToServer(
        `admin/golden-key/hints/${hintId}`,
        "PUT",
        JSON.stringify({ content, image_url: image_url || "" }),
        true
    );
}

export async function deleteGoldenKeyHint(hintId) {
    return deleteFromServer(`admin/golden-key/hints/${hintId}`);
}
