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
                LanguageProvider.ROUTES.navHome.path = this.#constructRoute(lang.code, "navHome");
                LanguageProvider.ROUTES.navEvents.path = this.#constructRoute(lang.code, "navEvents");
                LanguageProvider.ROUTES.navGeocaches.path = this.#constructRoute(lang.code, "navGeocaches");
                LanguageProvider.ROUTES.navShop.path = this.#constructRoute(lang.code, "navShop");
            } else {
                LanguageProvider.ROUTES.navHome.aliases.push(this.#constructRoute(lang.code, "navHome"));
                LanguageProvider.ROUTES.navEvents.aliases.push(this.#constructRoute(lang.code, "navEvents"));
                LanguageProvider.ROUTES.navGeocaches.aliases.push(this.#constructRoute(lang.code, "navGeocaches"));
                LanguageProvider.ROUTES.navShop.aliases.push(this.#constructRoute(lang.code, "navShop"));
            }
        });
    }

    #constructRoute(langCode, routeName) {
        return `/${langCode.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY[routeName][langCode]}`
    }
}

export { getLanguageFromPath, LanguageProvider }