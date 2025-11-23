import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import HttpBackend from 'i18next-http-backend';

i18n
  // Loads translations from /public/locales/{{lng}}.json
  .use(HttpBackend)

  // Detects language (from browser, localStorage, etc.)
  .use(LanguageDetector)

  // Connects i18next to React
  .use(initReactI18next)

  .init({
    fallbackLng: 'en',
    debug: false,

    interpolation: {
      escapeValue: false, // React already escapes
    },

    backend: {
      loadPath: '/locales/{{lng}}.json',
    },
  });

export default i18n;
