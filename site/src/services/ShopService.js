import { fetchFromServer } from "./fetcher"

async function getAllProducts(admin = false) {
    return fetchFromServer("products", admin);
}

export { getAllProducts };