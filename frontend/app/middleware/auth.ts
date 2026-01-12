export default defineNuxtRouteMiddleware((to, from) => {
  const token = useCookie('auth_token')

  console.log('Middleware Auth Check - Token:', token.value)

  if (!token.value && to.path !== '/') {
    return navigateTo('/')
  }
})
