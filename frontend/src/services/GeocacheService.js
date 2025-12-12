import { fetchFromServer } from "./fetcher";

async function getAllGeocaches(admin = false, page = null, search = "", perPage = null, sortBy = "", sortDirection = "") {
    return fetchFromServer("geocaches", admin, page, search, perPage, sortBy, sortDirection);
}

export { getAllGeocaches };