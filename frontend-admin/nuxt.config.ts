// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  
  // Modules
  modules: [
    '@nuxt/ui',
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@vueuse/nuxt',
    '@vee-validate/nuxt',
  ],

  // CSS
  css: ['~/assets/css/main.css'],

  // Runtime Config
  runtimeConfig: {
    // Private keys (only available on server-side)
    apiSecret: '',
    // Public keys (exposed to client-side)
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',
      storageBase: process.env.NUXT_PUBLIC_STORAGE_BASE || 'http://localhost:9000/cruisebooking',
      appName: process.env.NUXT_PUBLIC_APP_NAME || 'CruiseBooking Admin',
      appVersion: process.env.NUXT_PUBLIC_APP_VERSION || '1.0.0',
    },
  },

  // TypeScript
  typescript: {
    strict: true,
    typeCheck: true,
  },

  // Build
  build: {
    transpile: ['vue-echarts'],
  },

  // Nitro (Server)
  nitro: {
    preset: 'node-server',
  },

  // Color Mode
  colorMode: {
    preference: 'system',
    fallback: 'light',
    hid: 'nuxt-color-mode-script',
    globalName: '__NUXT_COLOR_MODE__',
    componentName: 'ColorScheme',
    classPrefix: '',
    classSuffix: '-mode',
    storageKey: 'nuxt-color-mode',
  },

  // UI Configuration
  ui: {
    global: true,
    icons: ['heroicons', 'simple-icons'],
  },

  // Tailwind CSS
  tailwindcss: {
    cssPath: '~/assets/css/main.css',
    configPath: 'tailwind.config.ts',
    exposeConfig: false,
    injectPosition: 0,
    viewer: true,
  },

  // Pinia
  pinia: {
    storesDirs: ['./stores/**'],
  },

  // VeeValidate
  veeValidate: {
    // options
  },

  // Vite
  vite: {
    define: {
      'process.env.DEBUG': false,
    },
  },

  // Development
  devServer: {
    port: 3000,
    host: '0.0.0.0',
  },

  // Compatibility
  compatibilityDate: '2026-02-10',
})
