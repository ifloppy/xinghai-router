export default defineEventHandler((event) => {
  const url = getRequestURL(event)
  const upstream = process.env.API_INTERNAL_URL || 'http://127.0.0.1:8080'
  return proxyRequest(event, `${upstream}${url.pathname.replace(/^\/api/, '')}${url.search}`)
})
