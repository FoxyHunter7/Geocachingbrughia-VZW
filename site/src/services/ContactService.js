import { fetchToServer } from "./fetcher";

async function postContact(email, subject, message) {
    const json = JSON.stringify({
        email: email,
        subject: subject,
        message: message
    });

    const result = await fetchToServer("contact/form/responses", "POST", json);
}

export { postContact };