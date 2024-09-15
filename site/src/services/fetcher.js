import config from "../data/config.js"
import LanguageProvider from "./LanguageService";

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

    if (includeCreds) {
        headers.credentials = "include";
    }

    try {
        const response = await fetch(`${config.apiUrl}${endpoint}${queryParams.all}`, {
            method: "GET",
            headers: headers,
            credentials: (includeCreds) ? "include" : "omit"
        });

        if (!response.ok) {
            throw new error("Bad fetch response", response);
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

        if (includeCreds) {
            headers.credentials = "include";
        }

        if (headerContentType) {
            headers["Content-Type"] = headerContentType;
        }

        const response = await fetch(`${config.apiUrl}${endpoint}`, {
            method: method,
            headers: headers,
            credentials: (includeCreds) ? "include" : "omit",
            body: body
        });

        if (!response.ok) {
            throw new Error(`Bad ${method} response: ${response.statusText}`);
        }

        return {
            success: true,
            data: await response.json()
        };
    } catch (err) {
        console.error(`Failed to ${method} (endpoint: ${endpoint})`);
        return {
            success: false,
            error: err.message
        };
    }
}

async function deleteFromServer(endpoint) {
    try {
        const response = await fetch(`${config.apiUrl}${endpoint}`, {
            method: "DELETE",
            credentials: "include"
        });

        if (!response.ok) {
            throw new Error(`Bad response: ${response.statusText}`);
        }

        return {
            success: true,
            data: await response.json()
        };
    } catch (err) {
        console.error(`Failed to delete (endpoint: ${endpoint})`);
        return {
            success: false,
            error: err.message
        };
    }
}

export { fetchFromServer, fetchToServer, deleteFromServer };