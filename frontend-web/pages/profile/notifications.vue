<template>
  <div class="notification-settings-page">
    <div class="max-w-4xl mx-auto py-8 px-4">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900">通知设置</h1>
        <p class="text-gray-500 mt-1">管理您希望接收的通知类型和渠道</p>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="w-8 h-8 border-2 border-gray-200 border-t-blue-500 rounded-full animate-spin mr-3"></div>
        <span class="text-gray-500">加载中...</span>
      </div>

      <!-- Settings Content -->
      <div v-else-if="settings" class="space-y-6">
        <!-- Notification Types -->
        <section class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div class="px-6 py-4 border-b border-gray-100">
            <h2 class="text-lg font-semibold text-gray-900">通知类型</h2>
            <p class="text-sm text-gray-500 mt-1">选择您希望接收的通知类别</p>
          </div>

          <div class="divide-y divide-gray-100">
            <div v-for="item in notificationTypes" :key="item.key" class="px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors">
              <div class="flex items-center gap-3">
                <div :class="['w-10 h-10 rounded-full flex items-center justify-center', item.iconBg]">
                  <svg class="w-5 h-5" :class="item.iconColor" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="item.icon"/>
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-gray-900">{{ item.label }}</div>
                  <div class="text-sm text-gray-500">{{ item.description }}</div>
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" :checked="settings[item.key]" class="sr-only peer" @change="toggleSetting(item.key)">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
          </div>
        </section>

        <!-- Notification Channels -->
        <section class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div class="px-6 py-4 border-b border-gray-100">
            <h2 class="text-lg font-semibold text-gray-900">通知渠道</h2>
            <p class="text-sm text-gray-500 mt-1">选择您偏好的通知接收方式</p>
          </div>

          <div class="divide-y divide-gray-100">
            <!-- WeChat -->
            <div class="px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
                  <svg class="w-5 h-5 text-green-600" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18z"/>
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-gray-900">微信消息</div>
                  <div class="text-sm text-gray-500">通过微信小程序或公众号接收通知</div>
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="settings.wechat_enabled" class="sr-only peer" @change="saveSettings">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>

            <!-- SMS -->
            <div class="px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
                  <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-gray-900">短信通知</div>
                  <div class="text-sm text-gray-500">重要通知将通过短信发送到您的手机</div>
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="settings.sms_enabled" class="sr-only peer" @change="saveSettings">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>

            <!-- Email -->
            <div class="px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-purple-100 flex items-center justify-center">
                  <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-gray-900">邮件通知</div>
                  <div class="text-sm text-gray-500">通知将发送到您的注册邮箱</div>
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="settings.email_enabled" class="sr-only peer" @change="saveSettings">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
          </div>
        </section>

        <!-- Quiet Hours -->
        <section class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div class="px-6 py-4 border-b border-gray-100">
            <h2 class="text-lg font-semibold text-gray-900">免打扰设置</h2>
            <p class="text-sm text-gray-500 mt-1">设置在特定时间段内不接收非紧急通知</p>
          </div>

          <div class="px-6 py-6 space-y-4">
            <label class="flex items-center gap-3 cursor-pointer">
              <input type="checkbox" v-model="settings.quiet_hours_enabled" class="w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500" @change="saveSettings">
              <span class="text-gray-700">启用免打扰模式</span>
            </label>

            <div v-if="settings.quiet_hours_enabled" class="pl-7 space-y-4">
              <div class="flex items-center gap-4">
                <div class="flex-1">
                  <label class="block text-sm font-medium text-gray-700 mb-1">开始时间</label>
                  <input 
                    type="time" 
                    v-model="settings.quiet_hours_start" 
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    @change="saveSettings"
                  >
                </div>
                <div class="text-gray-400 pt-6">至</div>
                <div class="flex-1">
                  <label class="block text-sm font-medium text-gray-700 mb-1">结束时间</label>
                  <input 
                    type="time" 
                    v-model="settings.quiet_hours_end" 
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    @change="saveSettings"
                  >
                </div>
              </div>
              <p class="text-sm text-gray-500">免打扰期间，仅紧急通知（如订单状态变更）会发送</p>
            </div>
          </div>
        </section>

        <!-- Save Button (Mobile) -->
        <div class="md:hidden">
          <button 
            @click="saveSettings" 
            :disabled="saving"
            class="w-full py-3 px-4 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ saving ? '保存中...' : '保存设置' }}
          </button>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex flex-col items-center justify-center py-12">
        <svg class="w-16 h-16 text-red-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
        <p class="text-gray-500 mb-4">加载设置失败</p>
        <button @click="fetchSettings" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          重试
        </button>
      </div>
    </div>

    <!-- Toast Notification -->
    <div v-if="showToast" class="fixed bottom-4 right-4 bg-gray-800 text-white px-4 py-3 rounded-lg shadow-lg flex items-center gap-2 z-50">
      <svg class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
      </svg>
      <span>设置已保存</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface NotificationSettings {
  order_enabled: boolean
  payment_enabled: boolean
  inventory_enabled: boolean
  system_enabled: boolean
  refund_enabled: boolean
  voyage_enabled: boolean
  promotion_enabled: boolean
  wechat_enabled: boolean
  sms_enabled: boolean
  email_enabled: boolean
  quiet_hours_start: string
  quiet_hours_end: string
  quiet_hours_enabled: boolean
}

const settings = ref<NotificationSettings | null>(null)
const loading = ref(true)
const saving = ref(false)
const error = ref(false)
const showToast = ref(false)

const notificationTypes = [
  {
    key: 'order_enabled',
    label: '订单通知',
    description: '订单创建、支付、确认等状态变更通知',
    icon: 'M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z',
    iconBg: 'bg-blue-100',
    iconColor: 'text-blue-600'
  },
  {
    key: 'payment_enabled',
    label: '支付通知',
    description: '支付成功、失败、退款等通知',
    icon: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z',
    iconBg: 'bg-green-100',
    iconColor: 'text-green-600'
  },
  {
    key: 'refund_enabled',
    label: '退款通知',
    description: '退款申请进度、审批结果通知',
    icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    iconBg: 'bg-orange-100',
    iconColor: 'text-orange-600'
  },
  {
    key: 'voyage_enabled',
    label: '航次提醒',
    description: '出发提醒、行程变更等通知',
    icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    iconBg: 'bg-indigo-100',
    iconColor: 'text-indigo-600'
  },
  {
    key: 'promotion_enabled',
    label: '促销优惠',
    description: '优惠活动、限时折扣通知',
    icon: 'M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z',
    iconBg: 'bg-pink-100',
    iconColor: 'text-pink-600'
  },
  {
    key: 'system_enabled',
    label: '系统通知',
    description: '账户安全、系统维护等重要通知',
    icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z',
    iconBg: 'bg-gray-100',
    iconColor: 'text-gray-600'
  }
]

const fetchSettings = async () => {
  loading.value = true
  error.value = false

  try {
    const { $api } = useNuxtApp()
    const response = await $api.get('/notifications/settings')
    settings.value = response.data
  } catch (err) {
    console.error('Failed to fetch notification settings:', err)
    error.value = true
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  if (!settings.value) return

  saving.value = true
  try {
    const { $api } = useNuxtApp()
    await $api.put('/notifications/settings', settings.value)
    
    showToast.value = true
    setTimeout(() => {
      showToast.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to save notification settings:', err)
    alert('保存设置失败，请重试')
  } finally {
    saving.value = false
  }
}

const toggleSetting = (key: keyof NotificationSettings) => {
  if (!settings.value) return
  settings.value[key] = !settings.value[key]
  saveSettings()
}

onMounted(() => {
  fetchSettings()
})
</script>

<style scoped>
.notification-settings-page {
  min-height: calc(100vh - 200px);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
