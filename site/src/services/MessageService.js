import fetchFromServer from "./fetcher";

async function getAllMessages() {
    return fetchFromServer("messages");
}

export { getAllMessages };