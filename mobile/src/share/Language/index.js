import i18next from "i18next";
import { initReactI18next } from "react-i18next";

import Vietnamese from "../../assets/languages/vi.json";
import English from "../../assets/languages/en.json";

i18next.use(initReactI18next).init({
    lng: "vi",
    resources: {
        vi: { translation: Vietnamese },
        en: { translation: English }
    },
    react: {
        useSuspense: false
    },
    compatibilityJSON: "v3"
});

export default i18next