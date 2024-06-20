import fetchFromServer from "./fetcher";

async function getAllMessages() {
    const messages = await fetchFromServer("messages");
    return messages
        .sort((a, b) => new Date(a.updated_at) - new Date(b.updated_at))
        .flatMap(message => message.translations.map(translation => ({
            title: translation.title,
            body: translation.body
        })));
}

export { getAllMessages };