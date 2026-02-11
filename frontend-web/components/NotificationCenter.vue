<template>
  <div class="notification-center">
    <!-- Trigger Button -->
    <button class="notification-trigger" @click="toggleDropdown">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
      </svg>
      <span v-if="unreadCount > 0" class="notification-badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
    </button>

    <!-- Dropdown Panel -->
    <div v-if="isOpen" class="notification-dropdown" @click.stop>
      <div class="notification-header">
        <h3>通知中心</h3>
        <div class="header-actions">
          <button v-if="unreadCount > 0" class="btn-mark-all" @click="markAllRead">
            全部已读
          </button>
          <button class="btn-settings" @click="goToSettings">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Filter Tabs -->
      <div class="notification-tabs">
        <button 
          :class="['tab', { active: activeTab === 'all' }]" 
          @click="activeTab = 'all'"
        >
          全部
        </button>
        <button 
          :class="['tab', { active: activeTab === 'unread' }]" 
          @click="activeTab = 'unread'"
        >
          未读 ({{ unreadCount }})
        </button>
      </div>

      <!-- Notification List -->
      <div class="notification-list" @scroll="handleScroll">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <span>加载中...</span>
        </div>

        <template v-else-if="notifications.length > 0">
          <div 
            v-for="notification in notifications" 
            :key="notification.id"
            :class="['notification-item', { unread: !notification.is_read }]"
            @click="handleNotificationClick(notification)"
          >
            <!-- Icon based on type -->
            <div :class="['notification-icon', notification.type]">
              <svg v-if="notification.type === 'order'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/>
              </svg>
              <svg v-else-if="notification.type === 'payment'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"/>
              </svg>
              <svg v-else-if="notification.type === 'refund'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <svg v-else-if="notification.type === 'inventory'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>

            <div class="notification-content">
              <div class="notification-title">{{ notification.title }}</div>
              <div class="notification-body">{{ notification.content }}</div>
              <div class="notification-meta">
                <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
                <span v-if="notification.priority === 'urgent'" class="priority urgent">紧急</span>
                <span v-else-if="notification.priority === 'high'" class="priority high">重要</span>
              </div>
            </div>

            <div class="notification-actions">
              <button v-if="!notification.is_read" class="btn-read" @click.stop="markAsRead(notification.id)">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
              </button>
              <button class="btn-archive" @click.stop="archive(notification.id)">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Load More -->
          <div v-if="hasMore" class="load-more">
            <button @click="loadMore" :disabled="loadingMore">
              {{ loadingMore ? '加载中...' : '加载更多' }}
            </button>
          </div>
        </template>

        <!-- Empty State -->
        <div v-else class="empty-state">
          <svg class="w-16 h-16 empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
          </svg>
          <p>{{ activeTab === 'unread' ? '没有未读通知' : '暂无通知' }}</p>
        </div>
      </div>

      <!-- Footer -->
      <div class="notification-footer">
        <NuxtLink to="/profile/notifications" class="view-all">
          查看全部通知
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
          </svg>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  is_read: boolean
  is_archived: boolean
  priority: string
  action_type?: string
  action_url?: string
  source_id?: number
  source_type?: string
  created_at: string
  read_at?: string
}

const isOpen = ref(false)
const loading = ref(false)
const loadingMore = ref(false)
const notifications = ref<Notification[]>([])
const unreadCount = ref(0)
const currentPage = ref(1)
const hasMore = ref(false)
const activeTab = ref<'all' | 'unread'>('all')

// Close dropdown when clicking outside
const handleClickOutside = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.notification-center')) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  fetchUnreadCount()
  fetchNotifications()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

watch(activeTab, () => {
  currentPage.value = 1
  notifications.value = []
  fetchNotifications()
})

const toggleDropdown = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    fetchNotifications()
  }
}

const fetchUnreadCount = async () => {
  try {
    const { $api } = useNuxtApp()
    const response = await $api.get('/notifications/unread-count')
    unreadCount.value = response.data.count
  } catch (error) {
    console.error('Failed to fetch unread count:', error)
  }
}

const fetchNotifications = async () => {
  loading.value = true
  try {
    const { $api } = useNuxtApp()
    const response = await $api.get('/notifications', {
      params: {
        page: currentPage.value,
        unread_only: activeTab.value === 'unread'
      }
    })
    
    if (currentPage.value === 1) {
      notifications.value = response.data.data
    } else {
      notifications.value.push(...response.data.data)
    }
    
    hasMore.value = response.data.data.length === 20
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  currentPage.value++
  loadingMore.value = true
  await fetchNotifications()
  loadingMore.value = false
}

const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement
  if (target.scrollTop + target.clientHeight >= target.scrollHeight - 50 && hasMore.value && !loadingMore.value) {
    loadMore()
  }
}

const markAsRead = async (id: number) => {
  try {
    const { $api } = useNuxtApp()
    await $api.post(`/notifications/${id}/read`)
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

const markAllRead = async () => {
  try {
    const { $api } = useNuxtApp()
    await $api.post('/notifications/read-all')
    notifications.value.forEach(n => n.is_read = true)
    unreadCount.value = 0
  } catch (error) {
    console.error('Failed to mark all as read:', error)
  }
}

const archive = async (id: number) => {
  try {
    const { $api } = useNuxtApp()
    await $api.post(`/notifications/${id}/archive`)
    notifications.value = notifications.value.filter(n => n.id !== id)
  } catch (error) {
    console.error('Failed to archive:', error)
  }
}

const handleNotificationClick = (notification: Notification) => {
  if (!notification.is_read) {
    markAsRead(notification.id)
  }
  
  // Handle navigation based on action type
  if (notification.action_url) {
    navigateTo(notification.action_url)
  } else if (notification.source_id && notification.action_type) {
    switch (notification.action_type) {
      case 'view_order':
        navigateTo(`/orders/${notification.source_id}`)
        break
      case 'view_voyage':
        navigateTo(`/voyages/${notification.source_id}`)
        break
      case 'view_cabin':
        navigateTo(`/cabins/${notification.source_id}`)
        break
      case 'view_refund':
        navigateTo(`/orders/refunds/${notification.source_id}`)
        break
    }
  }
  
  isOpen.value = false
}

const goToSettings = () => {
  navigateTo('/profile/notifications')
  isOpen.value = false
}

const formatTime = (timestamp: string) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return date.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.notification-center {
  position: relative;
}

.notification-trigger {
  @apply relative p-2 rounded-lg hover:bg-gray-100 transition-colors;
}

.notification-badge {
  @apply absolute -top-1 -right-1 w-5 h-5 flex items-center justify-center bg-red-500 text-white text-xs rounded-full;
}

.notification-dropdown {
  @apply absolute right-0 top-full mt-2 w-96 bg-white rounded-xl shadow-2xl border border-gray-200 overflow-hidden z-50;
}

.notification-header {
  @apply flex items-center justify-between px-4 py-3 border-b border-gray-100;
}

.notification-header h3 {
  @apply font-semibold text-gray-800;
}

.header-actions {
  @apply flex items-center gap-2;
}

.btn-mark-all {
  @apply text-sm text-blue-600 hover:text-blue-700;
}

.btn-settings {
  @apply p-1.5 rounded-lg hover:bg-gray-100 text-gray-500;
}

.notification-tabs {
  @apply flex border-b border-gray-100;
}

.tab {
  @apply flex-1 py-3 text-sm font-medium text-gray-500 hover:text-gray-700 transition-colors;
}

.tab.active {
  @apply text-blue-600 border-b-2 border-blue-600;
}

.notification-list {
  @apply max-h-96 overflow-y-auto;
}

.notification-item {
  @apply flex items-start gap-3 p-4 hover:bg-gray-50 cursor-pointer transition-colors border-b border-gray-100 last:border-0;
}

.notification-item.unread {
  @apply bg-blue-50/50;
}

.notification-icon {
  @apply w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0;
}

.notification-icon.order {
  @apply bg-blue-100 text-blue-600;
}

.notification-icon.payment {
  @apply bg-green-100 text-green-600;
}

.notification-icon.refund {
  @apply bg-orange-100 text-orange-600;
}

.notification-icon.inventory {
  @apply bg-red-100 text-red-600;
}

.notification-icon.system {
  @apply bg-gray-100 text-gray-600;
}

.notification-content {
  @apply flex-1 min-w-0;
}

.notification-title {
  @apply font-medium text-gray-800 truncate;
}

.notification-body {
  @apply text-sm text-gray-500 mt-1 line-clamp-2;
}

.notification-meta {
  @apply flex items-center gap-2 mt-1;
}

.notification-time {
  @apply text-xs text-gray-400;
}

.priority {
  @apply text-xs px-1.5 py-0.5 rounded;
}

.priority.urgent {
  @apply bg-red-100 text-red-600;
}

.priority.high {
  @apply bg-orange-100 text-orange-600;
}

.notification-actions {
  @apply flex items-center gap-1 opacity-0 group-hover:opacity-100;
}

.notification-item:hover .notification-actions {
  @apply opacity-100;
}

.btn-read,
.btn-archive {
  @apply p-1.5 rounded hover:bg-gray-200 text-gray-400 hover:text-gray-600 transition-colors;
}

.loading-state {
  @apply flex flex-col items-center justify-center py-8;
}

.spinner {
  @apply w-8 h-8 border-2 border-gray-200 border-t-blue-500 rounded-full animate-spin mb-2;
}

.empty-state {
  @apply flex flex-col items-center justify-center py-12 text-gray-400;
}

.empty-icon {
  @apply mb-4 text-gray-300;
}

.load-more {
  @apply p-4 text-center;
}

.load-more button {
  @apply text-sm text-blue-600 hover:text-blue-700 disabled:text-gray-400;
}

.notification-footer {
  @apply px-4 py-3 border-t border-gray-100;
}

.view-all {
  @apply flex items-center justify-center gap-1 text-sm text-blue-600 hover:text-blue-700;
}
</style>
