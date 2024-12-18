import { fetchFromServer } from "./fetcher"

async function getAllEvents(admin = false, page = null, search= "", perPage = null, sortBy = "", sortDirection = "asc") {
    return fetchFromServer("events", admin, page, search, perPage, sortBy, sortDirection);
}

async function getHomePageEvents() {
    return fetchFromServer("home_events", false, null, "", null, "", "asc");
}

export { getAllEvents, getHomePageEvents };