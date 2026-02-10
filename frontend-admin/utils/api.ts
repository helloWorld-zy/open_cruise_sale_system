import type { FetchOptions } from 'ofetch'

export function createApiClient() {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const api = $fetch.create({
    baseURL: config.public.apiBase,
    headers: {
      'Content-Type': 'application/json'
    },
    onRequest({ options }) {
      // Add auth token if available
      if (authStore.token) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${authStore.token}`
        }
      }
    },
    onResponseError({ response }) {
      // Handle auth errors
      if (response.status === 401) {
        authStore.logout()
      }
    }
  })

  return api
}

export type ApiClient = ReturnType<typeof createApiClient>
