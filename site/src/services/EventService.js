import { fetchFromServer } from "./fetcher"

async function getAllEvents(search, perPage, sortBy, sortDirection, page) {
    return fetchFromServer("events", search, perPage, sortBy, sortDirection, page);
}

async function getHomePageEvents() {
    return fetchFromServer("home_events");
}

export { getAllEvents, getHomePageEvents };