export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()

  // Initialize auth state
  authStore.initAuth()

  // Public routes that don't require authentication
  const publicRoutes = ['/login', '/register', '/forgot-password']

  // If route is public, allow access
  if (publicRoutes.includes(to.path)) {
    // If already authenticated, redirect to home
    if (authStore.isAuthenticated) {
      return navigateTo('/')
    }
    return
  }

  // Check authentication for protected routes
  if (!authStore.isAuthenticated) {
    return navigateTo('/login')
  }
})
