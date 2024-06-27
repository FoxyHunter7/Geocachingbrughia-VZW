import { postToServer } from "./fetcher";

async function postContact(email, subject, message) {
    const json = JSON.stringify({
        email: email,
        subject: subject,
        message: message
    });

    const result = await postToServer("contact/form/responses", json);
    console.log(result);
}

export { postContact };