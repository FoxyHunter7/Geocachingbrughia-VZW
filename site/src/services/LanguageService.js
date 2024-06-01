import config from "../data/config.json";

function getLanguageFromPath(path) {
    return path = path.split('/')[1]
}

class LanguageProvider {
    static DICTIONARY = {};
    static LANGUAGES;
    static CURR_LANG;

    async init() {
        await this.#fetchSupportedLanguages();
        await this.#fetchDitcionary();
        this.#setCurrentLanguage();
    }

    async #fetchSupportedLanguages() {
        const response = await fetch(`${config.apiUrl}languages`);
        LanguageProvider.LANGUAGES = await response.json();
    }

    async #fetchDitcionary() {
        const response = await fetch(`${config.apiUrl}static`);
        const unformattedDictionary = await response.json();

        unformattedDictionary.forEach(item => {
            LanguageProvider.DICTIONARY[item.property] = {}

            item.contents.forEach(contents => {
                LanguageProvider.DICTIONARY[item.property][contents.lang_code] = contents.content;
            });
        });
    }

    #setCurrentLanguage() {
        const storedLang = localStorage.getItem("language")

        if (LanguageProvider.LANGUAGES.some(lang => lang.code === storedLang)) {
            LanguageProvider.CURR_LANG = localStorage.getItem("language");
            return;
        }

        const browserPrefferedLang = navigator.language.split("-")[0].toLocaleUpperCase() ?? "";

        if (LanguageProvider.LANGUAGES.some(lang => lang.code === browserPrefferedLang)) {
            LanguageProvider.CURR_LANG = browserPrefferedLang;
            localStorage.setItem("language", browserPrefferedLang);
            return
        } else {
            LanguageProvider.CURR_LANG = config.defaultLang;
        }
    }
}

export { getLanguageFromPath, LanguageProvider }