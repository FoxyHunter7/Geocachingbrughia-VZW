import config from "../data/config.js"
import LanguageProvider from "./LanguageService.js";

// Helper to get JWT token from localStorage
function getAuthToken() {
    return localStorage.getItem('admin_token');
}

async function fetchFromServer(endpoint, includeCreds = false, page = null, search = "", perPage = null, sortBy = "", sortDirection = "") {
    const queryParams = {
        lang: (includeCreds) ? "" : `?lang=${LanguageProvider.CURR_LANG.value}`,
        page: (page) ? (includeCreds) ? `?page=${page}` : `&page=${page}` : "",
        search: (search) ? `&search=${search}` : "",
        perPage: (perPage) ? `&per_page=${perPage}` : "",
        sortBy: (sortBy) ? `&sort_by=${sortBy}` : "",
        sortDirection: (sortDirection) ? `&sort_direction=${sortDirection}` : ""
    };
    queryParams.all = `${queryParams.lang}${queryParams.search}${queryParams.perPage}${queryParams.sortBy}${queryParams.sortDirection}${queryParams.page}`;

    const headers = {
        "Content-Type": "application/json",
        "Accept": "application/json"
    };

    // Add JWT token for authenticated requests
    if (includeCreds) {
        const token = getAuthToken();
        if (token) {
            headers["Authorization"] = `Bearer ${token}`;
        }
    }

    try {
        const response = await fetch(`${config.apiUrl}${endpoint}${queryParams.all}`, {
            method: "GET",
            headers: headers
        });

        // Handle auth failures specifically
        if (response.status === 401 || response.status === 403) {
            console.error(`Auth failed for endpoint: ${endpoint}`);
            return { access_denied: true };
        }

        if (!response.ok) {
            throw new Error("Bad fetch response", response);
        }

        return await response.json();
    } catch (err) {
        console.error(`Failed to fetch (endpoint: ${endpoint})`);
        return [];
    }
}

async function fetchToServer(endpoint, method = "POST", body = "", includeCreds = false, headerContentType = "application/json") {
    try {
        const headers = {
            "Accept": "application/json"
        };

        // Add JWT token for authenticated requests
        if (includeCreds) {
            const token = getAuthToken();
            if (token) {
                headers["Authorization"] = `Bearer ${token}`;
            }
        }

        if (headerContentType) {
            headers["Content-Type"] = headerContentType;
        }

        const response = await fetch(`${config.apiUrl}${endpoint}`, {
            method: method,
            headers: headers,
            body: body
        });

        // Handle auth failures specifically
        if (response.status === 401 || response.status === 403) {
            console.error(`Auth failed for ${method} ${endpoint}`);
            return {
                success: false,
                error: "Unauthorized",
                access_denied: true
            };
        }

        if (!response.ok) {
            // Try to get error message from response body
            let errorData;
            try {
                errorData = await response.json();
            } catch {
                errorData = null;
            }
            return {
                success: false,
                error: `Bad ${method} response: ${response.statusText}`,
                data: errorData
            };
        }

        return {
            success: true,
            data: await response.json()
        };
    } catch (err) {
        console.error(`Failed to ${method} (endpoint: ${endpoint})`, err);
        return {
            success: false,
            error: err.message
        };
    }
}

async function deleteFromServer(endpoint) {
    try {
        const headers = {
            "Accept": "application/json"
        };

        // Add JWT token for authenticated requests
        const token = getAuthToken();
        if (token) {
            headers["Authorization"] = `Bearer ${token}`;
        }

        const response = await fetch(`${config.apiUrl}${endpoint}`, {
            method: "DELETE",
            headers: headers
        });

        // Handle auth failures specifically
        if (response.status === 401 || response.status === 403) {
            console.error(`Auth failed for DELETE ${endpoint}`);
            return {
                success: false,
                error: "Unauthorized",
                access_denied: true
            };
        }

        if (!response.ok) {
            // Try to get error message from response body
            let errorData;
            try {
                errorData = await response.json();
            } catch {
                errorData = null;
            }
            return {
                success: false,
                error: `Bad response: ${response.statusText}`,
                data: errorData
            };
        }

        return {
            success: true,
            data: await response.json()
        };
    } catch (err) {
        console.error(`Failed to delete (endpoint: ${endpoint})`, err);
        return {
            success: false,
            error: err.message
        };
    }
}

export { fetchFromServer, fetchToServer, deleteFromServer };