async function fetchAllEvents(lang) {
    const response = await fetch(`http://127.0.0.1:8000/api/events?lang=${lang}`)
    return await response.json()
}

export { fetchAllEvents }