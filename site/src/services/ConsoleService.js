import warnings from "../data/warnings.json";
import { StaticContentProvider } from "./StaticContentService";

function showConsoleDangerWarning() {
    console.log(`%c${warnings.consoleDangerTitle[StaticContentProvider.CURR_LANG]}`, "font-size: 30px; font-weight: bold; color: yellow")
    console.log(`%c${warnings.consoleDangerSubtitle[StaticContentProvider.CURR_LANG]}`, "font-size: 15px; font-weight: bold")
}

export { showConsoleDangerWarning }