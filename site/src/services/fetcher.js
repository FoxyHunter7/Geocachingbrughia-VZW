import config from "../data/config.json";
import LanguageProvider from "./LanguageService";

async function fetchFromServer(endpoint, search = "", perPage = null, sortBy = "", sortDirection = "") {
    const queryParams = {
        lang: `?lang=${LanguageProvider.CURR_LANG}`,
        search: (search) ? `&search=${search}` : "",
        perPage: (perPage) ? `&per_page=${perPage}` : "",
        sortBy: (sortBy) ? `&sort_by=${sortBy}` : "",
        sortDirection: (sortDirection) ? `&sort_direction=${sortDirection}` : ""
    };
    queryParams.all = `${queryParams.lang}${queryParams.search}${queryParams.perPage}${queryParams.sortBy}${queryParams.sortDirection}`;

    try {
        const response = await fetch(`${config.apiUrl}${endpoint}${queryParams.all}`);

        if (!response.ok) {
            throw new error("Bad fetch response", response);
        }

        return await response.json();
    } catch (err) {
        console.error(`Failed to fetch (endpoint: ${endpoint})`, err);
        return [];
    }
}

export default fetchFromServer;