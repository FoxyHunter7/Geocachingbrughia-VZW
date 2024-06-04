import warnings from "../data/warnings.json";
import { LanguageProvider } from "./LanguageService";

function showConsoleDangerWarning() {
    console.log(`%c${warnings.consoleDangerTitle[LanguageProvider.CURR_LANG]}`, "font-size: 30px; font-weight: bold; color: yellow")
    console.log(`%c${warnings.consoleDangerSubtitle[LanguageProvider.CURR_LANG]}`, "font-size: 15px; font-weight: bold")
}

export { showConsoleDangerWarning }