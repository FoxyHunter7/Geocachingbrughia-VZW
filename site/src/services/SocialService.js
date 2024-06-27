import { fetchFromServer } from "./fetcher";

async function getAllSocials() {
    return fetchFromServer("socials");
}

export { getAllSocials };