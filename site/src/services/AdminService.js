import { fetchFromServer, fetchToServer } from "./fetcher"

async function login(email, password) {
    const json = JSON.stringify({
        email: email,
        password: password
    });
    
    return fetchToServer("login", "POST", json, true);
}

async function getProfileData() {
    return fetchFromServer("profile", true);
}

async function logout() {
    return fetchFromServer("logout", true);
}

async function fetchEvents() {
    return fetchFromServer("admin/events", true);
}

async function postEvent(formData) {
    return fetchToServer("admin/events", "POST", formData, true, "");
}

async function updateEvent(id, formData) {
    return fetchToServer(`admin/events/${id}?_method=PUT`, "POST", formData, true, "");
}

export { login, getProfileData, logout, fetchEvents, postEvent, updateEvent };