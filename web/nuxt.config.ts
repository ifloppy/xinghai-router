export default defineNuxtConfig({
  css: ['~/src/style.css'],
  nitro: {
    prerender: {
      routes: ['/', '/auth'],
    },
  },
  routeRules: {
    '/': { prerender: true },
    '/auth': { prerender: true },
  },
  devServer: { port: 5173, host: '127.0.0.1' },
  compatibilityDate: '2026-07-16',
})
