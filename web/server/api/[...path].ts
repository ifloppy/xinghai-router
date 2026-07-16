export default defineEventHandler((event) => {
  const url = getRequestURL(event)
  return proxyRequest(event, `http://127.0.0.1:8080${url.pathname.replace(/^\/api/, '')}${url.search}`)
})
