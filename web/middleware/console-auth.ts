import { getToken } from '~/src/api'

export default defineNuxtRouteMiddleware(() => {
  if (import.meta.server) {
    const token = useCookie('xinghai.admin-token').value
    if (!token) return navigateTo('/auth', { redirectCode: 302 })
    return
  }
  if (!getToken()) return navigateTo('/auth', { redirectCode: 302 })
})
