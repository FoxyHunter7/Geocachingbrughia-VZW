import config from "../data/config.json";

function getLanguageFromPath(path) {
    return path = path.split('/')[1]
}

class LanguageProvider {
    static DICTIONARY = {};
    static LANGUAGES;
    static CURR_LANG;
    static ROUTES = {
        navHome: {
            path: "",
            aliases: []
        },
        navEvents: {
            path: "",
            aliases: []
        },
        navGeocaches: {
            path: "",
            aliases: []
        },
        navShop: {
            path: "",
            aliases: []
        }
    };

    async init() {
        await this.#fetchSupportedLanguages();
        await this.#fetchDitcionary();
        this.#setCurrentLanguage();
        this.#setRoutes();
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
            LanguageProvider.CURR_LANG = config.defaultLanguage;
        }
    }

    #setRoutes() {
        LanguageProvider.LANGUAGES.forEach(lang => {
            if (lang.code === config.defaultLanguage) {
                LanguageProvider.ROUTES.navHome.path = `/${lang.code.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY.navHome[lang.code]}`;
                LanguageProvider.ROUTES.navEvents.path = `/${lang.code.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY.navEvents[lang.code]}`;
                LanguageProvider.ROUTES.navGeocaches.path = `/${lang.code.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY.navGeocaches[lang.code]}`;
                LanguageProvider.ROUTES.navShop.path = `/${lang.code.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY.navShop[lang.code]}`;
            } else {
                LanguageProvider.ROUTES.navHome.aliases.push(`/${lang.code.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY.navHome[lang.code]}`);
                LanguageProvider.ROUTES.navEvents.aliases.push(`/${lang.code.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY.navEvents[lang.code]}`);
                LanguageProvider.ROUTES.navGeocaches.aliases.push(`/${lang.code.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY.navGeocaches[lang.code]}`);
                LanguageProvider.ROUTES.navShop.aliases.push(`/${lang.code.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY.navShop[lang.code]}`);
            }
        });
    }
}

export { getLanguageFromPath, LanguageProvider }