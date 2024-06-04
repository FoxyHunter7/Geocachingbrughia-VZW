import config from "../data/config.json";
import fallbackLanguages from "../data/fallbackLanguages.json";
import fallbackStaticContent from "../data/fallbackStaticContent.json";
import warnings from "../data/warnings.json";

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
    static ERRORS = {
        "langFetch": "",
        "staticFetch": ""
    };

    async init() {
        this.#setCurrentLanguage(fallbackLanguages);
        await this.#fetchSupportedLanguages();
        await this.#fetchDitcionary();
        this.#setCurrentLanguage(LanguageProvider.LANGUAGES);
        this.#setRoutes();

        console.log(`%c${warnings.consoleDangerTitle[LanguageProvider.CURR_LANG]}`, "font-size: 30px; font-weight: bold; color: yellow")
        console.log(`%c${warnings.consoleDangerSubtitle[LanguageProvider.CURR_LANG]}`, "font-size: 15px; font-weight: bold")
    }

    async #fetchSupportedLanguages() {
        try {
            const response = await fetch(`${config.apiUrl}languages`);

            if (!response.ok) {
                throw new ResponseError("Bad fetch response", response);
            }

            LanguageProvider.LANGUAGES = await response.json();
        } catch (err) {
            console.log("ERROR While fetching languages: ", err);

            if (err.response && err.response.status) {
                switch (err.response.status) {
                    case 400 || 401 || 404 || 500: LanguageProvider.ERRORS.langFetch = warnings.apiLangComm; break;
                    case 503: LanguageProvider.ERRORS.langFetch = warnings.apiLangOverloaded; break;
                }
            }

            LanguageProvider.LANGUAGES = fallbackLanguages;
        }
    }

    async #fetchDitcionary() {
        try {
            const response = await fetch(`${config.apiUrl}static`);

            if (!response.ok) {
                throw new ResponseError("Bad fetch response", response);
            }

            const unformattedDictionary = await response.json();

            unformattedDictionary.forEach(item => {
                LanguageProvider.DICTIONARY[item.property] = {}
    
                item.contents.forEach(contents => {
                    LanguageProvider.DICTIONARY[item.property][contents.lang_code] = contents.content;
                });
            });
        } catch (err) {
            console.log("ERROR While fetching static site content: ", err);

            if (err.response && err.response.status) {
                switch (err.response.status) {
                    case 400 || 401 || 404 || 500: LanguageProvider.ERRORS.staticFetch = warnings.apiLangComm; break;
                    case 503: LanguageProvider.ERRORS.staticFetch = warnings.apiLangOverloaded; break;
                }
            }

            LanguageProvider.DICTIONARY = fallbackStaticContent;
        }
    }

    #setCurrentLanguage(languages) {
        const storedLang = localStorage.getItem("language")

        if (languages.some(lang => lang.code === storedLang)) {
            LanguageProvider.CURR_LANG = localStorage.getItem("language");
            return;
        }

        const browserPrefferedLang = navigator.language.split("-")[0].toLocaleUpperCase() ?? "";

        if (languages.some(lang => lang.code === browserPrefferedLang)) {
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
                LanguageProvider.ROUTES.navHome.path = this.#constructRoute(lang.code, "NavHome");
                LanguageProvider.ROUTES.navEvents.path = this.#constructRoute(lang.code, "NavEvents");
                LanguageProvider.ROUTES.navGeocaches.path = this.#constructRoute(lang.code, "NavGeocaches");
                LanguageProvider.ROUTES.navShop.path = this.#constructRoute(lang.code, "NavShop");
            } else {
                LanguageProvider.ROUTES.navHome.aliases.push(this.#constructRoute(lang.code, "NavHome"));
                LanguageProvider.ROUTES.navEvents.aliases.push(this.#constructRoute(lang.code, "NavEvents"));
                LanguageProvider.ROUTES.navGeocaches.aliases.push(this.#constructRoute(lang.code, "NavGeocaches"));
                LanguageProvider.ROUTES.navShop.aliases.push(this.#constructRoute(lang.code, "NavShop"));
            }
        });
    }

    #constructRoute(langCode, routeName) {
        return `/${langCode.toLocaleLowerCase()}/${LanguageProvider.DICTIONARY[routeName][langCode]}`
    }
}

export { getLanguageFromPath, LanguageProvider }