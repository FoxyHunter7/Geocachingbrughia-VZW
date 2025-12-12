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

async function fetchEvents(curr_page) {
    return fetchFromServer("admin/events", true, curr_page);
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

async function fetchGeocaches(curr_page) {
    return fetchFromServer("admin/geocaches", true, curr_page);
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

/* -----
SOCIALS
----- */

async function fetchSocials(curr_page) {
    return fetchFromServer("admin/socials", true, curr_page);
}

async function postSocial(formData) {
    return fetchToServer("admin/socials", "POST", formData, true, "");
}

async function updateSocial(id, formData) {
    return fetchToServer(`admin/socials/${id}?_method=PUT`, "POST", formData, true, "");
}

async function deleteSocial(id) {
    return deleteFromServer(`admin/socials/${id}`);
}

/* -----
MESSAGES
----- */

async function fetchMessages(curr_page) {
    return fetchFromServer("admin/messages", true, curr_page);
}

async function postMessage(formData) {
    return fetchToServer("admin/messages", "POST", formData, true, "");
}

async function updateMessage(id, formData) {
    return fetchToServer(`admin/messages/${id}?_method=PUT`, "POST", formData, true, "");
}

async function deleteMessage(id) {
    return deleteFromServer(`admin/messages/${id}`);
}

/* -----
FORMRESPONSES
----- */

async function fetchContactFormResponses(curr_page) {
    return fetchFromServer("admin/formresponses", true, curr_page);
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
    deleteGeocache,
    fetchSocials,
    postSocial,
    updateSocial,
    deleteSocial,
    fetchMessages,
    postMessage,
    updateMessage,
    deleteMessage,
    fetchContactFormResponses
};