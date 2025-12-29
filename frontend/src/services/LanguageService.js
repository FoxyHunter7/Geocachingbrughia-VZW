import { StaticContentProvider } from "./StaticContentService";
import { ref } from 'vue';
import config from '../data/config';

class LanguageProvider {
    static #CURR_LANG = ref('');

    static get CURR_LANG() {
        return LanguageProvider.#CURR_LANG;
    }

    static set CURR_LANG(newLang) {
        LanguageProvider.#CURR_LANG.value = newLang;
        localStorage.setItem("language", newLang);
    }

    static autoSetLanguage(supportLanguages = StaticContentProvider.LANGUAGES) {
        const storedLang = localStorage.getItem("language");

        if (supportLanguages.some(lang => lang.code === storedLang)) {
            LanguageProvider.#CURR_LANG.value = localStorage.getItem("language");
            return;
        }

        const browserPrefferedLang = navigator.language.split("-")[0].toLocaleUpperCase() ?? "";

        if (supportLanguages.some(lang => lang.code === browserPrefferedLang)) {
            LanguageProvider.#CURR_LANG.value = browserPrefferedLang;
            localStorage.setItem("language", browserPrefferedLang);
            return
        } else {
            LanguageProvider.#CURR_LANG.value = config.defaultLanguage;
        }
    }
}

export default LanguageProvider;
export { LanguageProvider };