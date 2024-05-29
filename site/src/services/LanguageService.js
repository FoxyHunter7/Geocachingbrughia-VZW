function defaultLanguage() {
    // TODO: actually some logic to detirmine the best language for the user
    return 'en'
}

function getLanguageFromPath(path) {
    return path = path.split('/')[1]
}

export { defaultLanguage, getLanguageFromPath }