import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  modules: ['@nuxt/eslint'],
  css: ['~/src/style.css'],
  components: [
    { path: '~/components', pattern: '**/*.vue' },
  ],
  app: {
    head: {
      link: [
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=DM+Mono:wght@400;500&family=Manrope:wght@400;500;600;700;800&display=swap',
        },
      ],
    },
  },
  vite: {
    plugins: [tailwindcss()],
  },
  nitro: {
    prerender: {
      routes: ['/', '/auth', '/rankings', '/terms', '/privacy'],
    },
  },
  routeRules: {
    '/': { prerender: true },
    '/auth': { prerender: true },
    '/rankings': { prerender: true },
    '/terms': { prerender: true },
    '/privacy': { prerender: true },
  },
  devServer: { port: 5173, host: '127.0.0.1' },
  compatibilityDate: '2026-07-16',
})
