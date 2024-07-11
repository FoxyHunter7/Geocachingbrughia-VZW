import { deleteFromServer, fetchFromServer, fetchToServer } from "./fetcher"

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

/* -----
EVENTS
----- */

async function fetchEvents() {
    return fetchFromServer("admin/events", true);
}

async function postEvent(formData) {
    return fetchToServer("admin/events", "POST", formData, true, "");
}

async function updateEvent(id, formData) {
    return fetchToServer(`admin/events/${id}?_method=PUT`, "POST", formData, true, "");
}

async function deleteEvent(id) {
    return deleteFromServer(`admin/events/${id}`);
}

/* -----
GEOCACHES
----- */

async function fetchGeocaches() {
    return fetchFromServer("admin/geocaches", true);
}

async function postGeocache(formData) {
    return fetchToServer("admin/geocaches", "POST", formData, true, "");
}

async function updateGeocache(id, formData) {
    return fetchToServer(`admin/geocaches/${id}?_method=PUT`, "POST", formData, true, "");
}

async function deleteGeocache(id) {
    return deleteFromServer(`admin/geocaches/${id}`);
}

export {
    login, 
    getProfileData, 
    logout, 
    fetchEvents, 
    postEvent, 
    updateEvent, 
    deleteEvent, 
    fetchGeocaches,
    postGeocache,
    updateGeocache,
    deleteGeocache
};