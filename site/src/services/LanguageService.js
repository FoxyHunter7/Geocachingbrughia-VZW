function getLanguageFromPath(path) {
    return path = path.split('/')[1]
}

export { getLanguageFromPath }