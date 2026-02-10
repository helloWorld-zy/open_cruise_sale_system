import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '~/types'

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref<string | null>(null)
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'super_admin')
  const userRole = computed(() => user.value?.role || '')

  // Actions
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const clearToken = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  const login = async (username: string, password: string) => {
    loading.value = true
    error.value = null
    
    try {
      const { $api } = useNuxtApp()
      const response = await $api('/auth/login', {
        method: 'POST',
        body: { username, password }
      })
      
      setToken(response.token)
      user.value = response.user
      
      return response
    } catch (err: any) {
      error.value = err.message || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const logout = () => {
    clearToken()
    navigateTo('/login')
  }

  const fetchUser = async () => {
    if (!token.value) return
    
    try {
      const { $api } = useNuxtApp()
      const response = await $api('/auth/me')
      user.value = response
    } catch (err) {
      clearToken()
    }
  }

  const initAuth = () => {
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      token.value = savedToken
      fetchUser()
    }
  }

  return {
    // State
    token,
    user,
    loading,
    error,
    // Getters
    isAuthenticated,
    isAdmin,
    userRole,
    // Actions
    login,
    logout,
    setToken,
    clearToken,
    fetchUser,
    initAuth
  }
})
