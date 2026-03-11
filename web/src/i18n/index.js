import { createI18n } from 'vue-i18n'
import axios from 'axios'

// Detect the best language to use
function detectLanguage() {
  const saved = localStorage.getItem('locale')
  if (saved) return saved
  const browserLang = navigator.language || navigator.languages?.[0] || 'en'
  return browserLang.split('-')[0].toLowerCase() // default to browser language short code
}

const i18n = createI18n({
  legacy: false,
  locale: detectLanguage(),
  fallbackLocale: 'en',
  messages: {}
})

let loadedLocales = []

// Load locale messages from the backend API
export async function loadLocale(lang) {
  if (loadedLocales.includes(lang)) {
    i18n.global.locale.value = lang
    localStorage.setItem('locale', lang)
    document.documentElement.lang = lang === 'zh' ? 'zh-CN' : lang
    return
  }

  try {
    const { data } = await axios.get(`/api/locales/${lang}?t=${Date.now()}`)
    
    if (typeof data !== 'object' || data === null) {
      throw new Error('Received non-JSON response for locale data')
    }

    i18n.global.setLocaleMessage(lang, data)
    loadedLocales.push(lang)
    i18n.global.locale.value = lang
    localStorage.setItem('locale', lang)
    document.documentElement.lang = lang === 'zh' ? 'zh-CN' : lang
  } catch (e) {
    if (lang !== 'en') {
      console.warn(`Failed to load locale: ${lang}, falling back to en`)
      await loadLocale('en')
    } else {
      console.error(`Failed to load default locale en`, e)
    }
  }
}

// Get the current locale
export function getCurrentLocale() {
  return i18n.global.locale.value
}

export default i18n
