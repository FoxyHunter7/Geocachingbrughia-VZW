import { ref } from "vue";
import config from "../data/config.js";
import fallbackLanguages from "../data/fallbackLanguages.json";
import fallbackStaticContent from "../data/fallbackStaticContent.json";
import warnings from "../data/warnings.json";
import { showConsoleDangerWarning } from "./ConsoleService.js";
import { LanguageProvider } from "./LanguageService.js";

class StaticContentProvider {
    static DICTIONARY = {};
    static LANGUAGES;
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
    static ERRORS = "";
    static INIT_COMPLETE = false;
    static DICT_ON_FALLBACK = ref(false);

    async init() {
        await this.#fetchSupportedLanguages();
        await this.#fetchDitcionary();

        LanguageProvider.autoSetLanguage(StaticContentProvider.LANGUAGES)
        this.#setRoutes();
        
        StaticContentProvider.INIT_COMPLETE = true;
        showConsoleDangerWarning();
    }

    async #fetchSupportedLanguages() {
        try {
            const response = await fetch(`${config.apiUrl}languages`);

            if (!response.ok) {
                throw new Error("Bad fetch response", response);
            }

            const languages = await response.json();

            // If API returns empty array, use fallback
            if (!languages || languages.length === 0) {
                throw new Error("Empty languages from API");
            }

            StaticContentProvider.LANGUAGES = languages;
        } catch (err) {
            console.error("Failed to fetch (site languages)");
            this.#setErrors(err);
            StaticContentProvider.LANGUAGES = fallbackLanguages;
        }
    }

    async #fetchDitcionary() {
        try {
            const response = await fetch(`${config.apiUrl}static`);

            if (!response.ok) {
                throw new Error("Bad fetch response", response);
            }

            const unformattedDictionary = await response.json();

            // If API returns empty array, use fallback
            if (!unformattedDictionary || unformattedDictionary.length === 0) {
                throw new Error("Empty static content from API");
            }

            unformattedDictionary.forEach(item => {
                StaticContentProvider.DICTIONARY[item.property] = {}
    
                item.contents.forEach(contents => {
                    StaticContentProvider.DICTIONARY[item.property][contents.lang_code] = contents.content;
                });
            });
        } catch (err) {
            console.error("Failed to fetch (static site content)");
            this.#setErrors(err);
            StaticContentProvider.DICT_ON_FALLBACK.value = true;
            StaticContentProvider.DICTIONARY = fallbackStaticContent;
        }
    }

    #setRoutes() {
        StaticContentProvider.LANGUAGES.forEach(lang => {
            if (lang.code === config.defaultLanguage) {
                StaticContentProvider.ROUTES.navHome.path = StaticContentProvider.constructRoute(lang.code, "NavHome");
                StaticContentProvider.ROUTES.navEvents.path = StaticContentProvider.constructRoute(lang.code, "NavEvents");
                StaticContentProvider.ROUTES.navGeocaches.path = StaticContentProvider.constructRoute(lang.code, "NavGeocaches");
                StaticContentProvider.ROUTES.navShop.path = StaticContentProvider.constructRoute(lang.code, "NavShop");
            } else {
                StaticContentProvider.ROUTES.navHome.aliases.push(StaticContentProvider.constructRoute(lang.code, "NavHome"));
                StaticContentProvider.ROUTES.navEvents.aliases.push(StaticContentProvider.constructRoute(lang.code, "NavEvents"));
                StaticContentProvider.ROUTES.navGeocaches.aliases.push(StaticContentProvider.constructRoute(lang.code, "NavGeocaches"));
                StaticContentProvider.ROUTES.navShop.aliases.push(StaticContentProvider.constructRoute(lang.code, "NavShop"));
            }

            StaticContentProvider.ROUTES.navHome.aliases.push(StaticContentProvider.constructRoute(lang.code, "NavHome", lang.code));
            StaticContentProvider.ROUTES.navEvents.aliases.push(StaticContentProvider.constructRoute(lang.code, "NavEvents", lang.code));
            StaticContentProvider.ROUTES.navGeocaches.aliases.push(StaticContentProvider.constructRoute(lang.code, "NavGeocaches", lang.code));
            StaticContentProvider.ROUTES.navShop.aliases.push(StaticContentProvider.constructRoute(lang.code, "NavShop", lang.code));
        });
    }

    static constructRoute(langFromDict, routeName, langCode = "none") {
        return `/${(langCode == "none") ? "" : langCode.toLocaleLowerCase() + "/"}${StaticContentProvider.DICTIONARY[routeName][langFromDict]}`
    }

    #setErrors(err) {
        if (err.response && err.response.status) {
            switch (err.response.status) {
                case 400 || 401 || 404 || 500: StaticContentProvider.ERRORS = warnings.apiComm; break;
                case 503: StaticContentProvider.ERRORS = warnings.apiOverloaded; break;
            }
        }  else {
            StaticContentProvider.ERRORS = warnings.apiComm;
        }
    }
}

export default StaticContentProvider;
export { StaticContentProvider };