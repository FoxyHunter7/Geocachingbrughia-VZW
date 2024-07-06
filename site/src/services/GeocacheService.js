import { fetchFromServer } from "./fetcher";

async function getAllGeocaches(search, perPage, sortBy, sortDirection, page) {
    return fetchFromServer("geocaches", search, perPage, sortBy, sortDirection, page);
}

export { getAllGeocaches };