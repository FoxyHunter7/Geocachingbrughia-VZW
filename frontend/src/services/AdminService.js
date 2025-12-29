import { deleteFromServer, fetchFromServer, fetchToServer } from "./fetcher"

async function login(email, password) {
    const json = JSON.stringify({
        email: email,
        password: password
    });
    
    return fetchToServer("login", "POST", json, false);
}

async function changePassword(currentPassword, newPassword) {
    const json = JSON.stringify({
        current_password: currentPassword,
        new_password: newPassword
    });
    
    return fetchToServer("admin/change-password", "POST", json, true);
}

async function getProfileData() {
    return fetchFromServer("admin/profile", true);
}

async function logout() {
    return fetchToServer("admin/logout", "POST", "", true);
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
    return fetchToServer(`admin/events/${id}`, "PUT", formData, true, "");
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
    return fetchToServer(`admin/geocaches/${id}`, "PUT", formData, true, "");
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
    return fetchToServer(`admin/socials/${id}`, "PUT", formData, true, "");
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
    return fetchToServer(`admin/messages/${id}`, "PUT", formData, true, "");
}

async function deleteMessage(id) {
    return deleteFromServer(`admin/messages/${id}`);
}

/* -----
CONTACTS (formerly FORMRESPONSES)
----- */

async function fetchContactFormResponses(curr_page) {
    return fetchFromServer("admin/contacts", true, curr_page);
}

async function fetchContact(id) {
    return fetchFromServer(`admin/contacts/${id}`, true);
}

async function updateContact(id, data) {
    const json = JSON.stringify(data);
    return fetchToServer(`admin/contacts/${id}`, "PUT", json, true);
}

async function addContactNote(contactId, note) {
    const json = JSON.stringify({ content: note });
    return fetchToServer(`admin/contacts/${contactId}/notes`, "POST", json, true);
}

/* -----
USERS
----- */

async function fetchUsers() {
    return fetchFromServer("admin/users", true);
}

async function createUser(name, email) {
    const json = JSON.stringify({ name, email });
    return fetchToServer("admin/users", "POST", json, true);
}

async function deleteUser(id) {
    return deleteFromServer(`admin/users/${id}`);
}

async function resendInvitation(id) {
    return fetchToServer(`admin/users/${id}/resend-invitation`, "POST", "", true);
}

export {
    login,
    changePassword,
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
    fetchContactFormResponses,
    fetchContact,
    updateContact,
    addContactNote,
    fetchUsers,
    createUser,
    deleteUser,
    resendInvitation
};