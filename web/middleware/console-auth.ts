import { getToken } from '~/src/api'

export default defineNuxtRouteMiddleware(() => {
  if (import.meta.server) return
  if (!getToken()) return navigateTo('/auth')
})
