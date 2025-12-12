import { fetchToServer } from "./fetcher";

async function postContact(email, subject, message) {
    const json = JSON.stringify({
        email: email,
        subject: subject,
        message: message
    });

    const response = await fetchToServer("contact/form/responses", "POST", json);
    return response;
}

export { postContact };