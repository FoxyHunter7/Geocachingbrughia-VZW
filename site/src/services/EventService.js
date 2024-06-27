import { fetchFromServer } from "./fetcher"

async function getAllEvents(search, perPage, sortBy, sortDirection) {
    return fetchFromServer("events", search, perPage, sortBy, sortDirection);
}

async function getHomePageEvents() {
    return fetchFromServer("home_events");
}

export { getAllEvents, getHomePageEvents };