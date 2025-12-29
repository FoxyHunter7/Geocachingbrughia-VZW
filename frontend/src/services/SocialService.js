import { fetchFromServer } from "./fetcher";

async function getAllSocials(admin = false) {
    return fetchFromServer("socials", admin);
}

export { getAllSocials };