import { StaticContentProvider } from "./StaticContentService";

class LanguageProvider {
    static CURR_LANG;

    static autoSetLanguage(supportLanguages = StaticContentProvider.LANGUAGES) {
        const storedLang = localStorage.getItem("language");

        if (supportLanguages.some(lang => lang.code === storedLang)) {
            LanguageProvider.CURR_LANG = localStorage.getItem("language");
            return;
        }

        const browserPrefferedLang = navigator.language.split("-")[0].toLocaleUpperCase() ?? "";

        if (supportLanguages.some(lang => lang.code === browserPrefferedLang)) {
            LanguageProvider.CURR_LANG = browserPrefferedLang;
            localStorage.setItem("language", browserPrefferedLang);
            return
        } else {
            LanguageProvider.CURR_LANG = config.defaultLanguage;
        }
    }
}

export default LanguageProvider;
export { LanguageProvider };