<template>
  <div class="payment-result">
    <div class="result-container">
      <!-- Success State -->
      <div v-if="status === 'success'" class="result-content success">
        <div class="icon-wrapper">
          <svg class="success-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <h1 class="title">支付成功！</h1>
        <p class="subtitle">您的订单已确认，我们将尽快为您安排</p>
        
        <div class="order-info" v-if="order">
          <div class="info-item">
            <span class="label">订单号</span>
            <span class="value">{{ order.order_number }}</span>
          </div>
          <div class="info-item">
            <span class="label">支付金额</span>
            <span class="value amount">¥{{ formatPrice(order.total_amount) }}</span>
          </div>
          <div class="info-item">
            <span class="label">支付时间</span>
            <span class="value">{{ formatDate(order.paid_at) }}</span>
          </div>
        </div>

        <div class="actions">
          <NuxtLink :to="`/orders/${orderId}`" class="btn-primary">
            查看订单详情
          </NuxtLink>
          <NuxtLink to="/" class="btn-secondary">
            返回首页
          </NuxtLink>
        </div>
      </div>

      <!-- Failed State -->
      <div v-else-if="status === 'failed'" class="result-content failed">
        <div class="icon-wrapper">
          <svg class="failed-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <h1 class="title">支付失败</h1>
        <p class="subtitle">抱歉，支付过程中出现问题</p>
        <p class="error-message" v-if="errorMessage">{{ errorMessage }}</p>
        
        <div class="actions">
          <NuxtLink :to="`/payment/${orderId}`" class="btn-primary">
            重新支付
          </NuxtLink>
          <NuxtLink to="/orders" class="btn-secondary">
            查看订单
          </NuxtLink>
        </div>
      </div>

      <!-- Cancelled State -->
      <div v-else-if="status === 'cancelled'" class="result-content cancelled">
        <div class="icon-wrapper">
          <svg class="cancelled-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </div>
        <h1 class="title">支付已取消</h1>
        <p class="subtitle">您已取消本次支付</p>
        
        <div class="actions">
          <NuxtLink :to="`/payment/${orderId}`" class="btn-primary">
            重新支付
          </NuxtLink>
          <NuxtLink to="/cruises" class="btn-secondary">
            浏览其他邮轮
          </NuxtLink>
        </div>
      </div>

      <!-- Loading State -->
      <div v-else class="result-content loading">
        <div class="loading-spinner"></div>
        <h1 class="title">正在查询支付结果...</h1>
        <p class="subtitle">请稍候，正在确认您的支付状态</p>
      </div>
    </div>

    <!-- Help Section -->
    <div class="help-section">
      <h3>遇到问题？</h3>
      <div class="help-links">
        <NuxtLink to="/help/payment">支付常见问题</NuxtLink>
        <NuxtLink to="/contact">联系客服</NuxtLink>
        <NuxtLink to="/orders">查看我的订单</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const orderId = route.query.order_id as string
const status = ref<'loading' | 'success' | 'failed' | 'cancelled'>('loading')
const order = ref<any>(null)
const errorMessage = ref('')

onMounted(async () => {
  await checkPaymentStatus()
})

const checkPaymentStatus = async () => {
  if (!orderId) {
    status.value = 'failed'
    errorMessage.value = '订单信息缺失'
    return
  }

  try {
    // Query order status
    const response = await $fetch(`/api/orders/${orderId}`)
    order.value = response.data

    // Determine status based on order state
    if (order.value.status === 'paid' || order.value.status === 'confirmed') {
      status.value = 'success'
    } else if (order.value.status === 'cancelled') {
      status.value = 'cancelled'
    } else if (order.value.payment_status === 'failed') {
      status.value = 'failed'
      errorMessage.value = '支付失败，请重试'
    } else {
      // Still pending, wait and retry
      setTimeout(checkPaymentStatus, 2000)
    }
  } catch (err: any) {
    status.value = 'failed'
    errorMessage.value = err.message || '查询订单状态失败'
  }
}

const formatPrice = (price: number) => {
  return price?.toLocaleString('zh-CN') || '0'
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.payment-result {
  @apply min-h-screen bg-gray-100 py-12 px-4;
}

.result-container {
  @apply max-w-lg mx-auto bg-white rounded-xl shadow-lg p-8;
}

.result-content {
  @apply text-center;
}

.icon-wrapper {
  @apply w-20 h-20 mx-auto mb-6 flex items-center justify-center rounded-full;
}

.success .icon-wrapper {
  @apply bg-green-100;
}

.failed .icon-wrapper {
  @apply bg-red-100;
}

.cancelled .icon-wrapper {
  @apply bg-gray-100;
}

.success-icon {
  @apply w-12 h-12 text-green-500;
}

.failed-icon {
  @apply w-12 h-12 text-red-500;
}

.cancelled-icon {
  @apply w-12 h-12 text-gray-500;
}

.title {
  @apply text-2xl font-bold text-gray-800 mb-2;
}

.subtitle {
  @apply text-gray-500 mb-6;
}

.error-message {
  @apply text-red-500 mb-6;
}

.order-info {
  @apply bg-gray-50 rounded-lg p-4 mb-6 text-left;
}

.info-item {
  @apply flex justify-between py-2;
}

.info-item:not(:last-child) {
  @apply border-b border-gray-200;
}

.label {
  @apply text-gray-500;
}

.value {
  @apply font-medium text-gray-800;
}

.value.amount {
  @apply text-orange-500 font-bold;
}

.actions {
  @apply space-y-3;
}

.btn-primary {
  @apply block w-full py-3 bg-blue-500 text-white rounded-lg font-medium hover:bg-blue-600 transition-colors;
}

.btn-secondary {
  @apply block w-full py-3 border border-gray-300 text-gray-600 rounded-lg font-medium hover:bg-gray-50 transition-colors;
}

.loading-spinner {
  @apply w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4;
}

.help-section {
  @apply max-w-lg mx-auto mt-8 text-center;
}

.help-section h3 {
  @apply text-gray-600 mb-4;
}

.help-links {
  @apply flex justify-center gap-6;
}

.help-links a {
  @apply text-blue-500 hover:text-blue-600;
}
</style>
