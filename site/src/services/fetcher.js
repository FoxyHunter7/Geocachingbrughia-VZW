import config from "../data/config.json";
import LanguageProvider from "./LanguageService";

async function fetchFromServer(endpoint, search = "", perPage = null, sortBy = "", sortDirection = "", page = "",) {
    const queryParams = {
        lang: `?lang=${LanguageProvider.CURR_LANG.value}`,
        page: (page) ? `&page=${page}` : "",
        search: (search) ? `&search=${search}` : "",
        perPage: (perPage) ? `&per_page=${perPage}` : "",
        sortBy: (sortBy) ? `&sort_by=${sortBy}` : "",
        sortDirection: (sortDirection) ? `&sort_direction=${sortDirection}` : ""
    };
    queryParams.all = `${queryParams.lang}${queryParams.search}${queryParams.perPage}${queryParams.sortBy}${queryParams.sortDirection}${queryParams.page}`;

    try {
        const response = await fetch(`${config.apiUrl}${endpoint}${queryParams.all}`);

        if (!response.ok) {
            throw new error("Bad fetch response", response);
        }

        return await response.json();
    } catch (err) {
        console.error(`Failed to fetch (endpoint: ${endpoint})`);
        return [];
    }
}

async function postToServer(endpoint, json) {
    try {
        const response = await fetch(`${config.apiUrl}${endpoint}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: json
        });

        if (!response.ok) {
            throw new Error(`Bad post response: ${response.statusText}`);
        }

        return {
            success: true,
            data: await response.json()
        };
    } catch (err) {
        console.error(`Failed to post (endpoint: ${endpoint})`);
        return {
            success: false,
            error: err.message
        };
    }
}

export { fetchFromServer, postToServer };