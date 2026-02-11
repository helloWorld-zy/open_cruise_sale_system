<template>
  <div class="payment-page">
    <div class="payment-container">
      <!-- Header -->
      <div class="payment-header">
        <h1 class="text-2xl font-bold">订单支付</h1>
        <p class="text-gray-500 mt-2">订单号: {{ order?.order_number }}</p>
      </div>

      <!-- Order Summary -->
      <div class="order-summary" v-if="order">
        <h3 class="text-lg font-semibold mb-4">订单详情</h3>
        <div class="summary-items">
          <div class="summary-item">
            <span>邮轮</span>
            <span>{{ order.cruise?.name_cn }}</span>
          </div>
          <div class="summary-item">
            <span>航次</span>
            <span>{{ order.voyage?.voyage_number }}</span>
          </div>
          <div class="summary-item">
            <span>出发日期</span>
            <span>{{ order.voyage?.departure_date }}</span>
          </div>
          <div class="summary-item">
            <span>舱房数量</span>
            <span>{{ order.cabin_count }} 间</span>
          </div>
          <div class="summary-item">
            <span>乘客人数</span>
            <span>{{ order.passenger_count }} 人</span>
          </div>
        </div>
        <div class="summary-total">
          <span>应付金额</span>
          <span class="amount">¥{{ formatPrice(order.total_amount) }}</span>
        </div>
      </div>

      <!-- Payment Methods -->
      <div class="payment-methods">
        <h3 class="text-lg font-semibold mb-4">选择支付方式</h3>
        
        <div class="methods-list">
          <div
            v-for="method in paymentMethods"
            :key="method.id"
            class="method-item"
            :class="{ 'selected': selectedMethod === method.id }"
            @click="selectedMethod = method.id"
          >
            <div class="method-icon">
              <img :src="method.icon" :alt="method.name" />
            </div>
            <div class="method-info">
              <div class="method-name">{{ method.name }}</div>
              <div class="method-desc">{{ method.description }}</div>
            </div>
            <div class="check-icon" v-if="selectedMethod === method.id">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Payment QR Code (for WeChat/Alipay) -->
      <div v-if="qrCodeUrl" class="qr-section">
        <div class="qr-container">
          <img :src="qrCodeUrl" alt="支付二维码" class="qr-code" />
          <p class="qr-hint">请使用微信扫描二维码完成支付</p>
        </div>
      </div>

      <!-- Countdown -->
      <div class="countdown-section" v-if="timeLeft > 0">
        <p class="text-gray-600">
          支付剩余时间: 
          <span class="countdown">{{ formatTime(timeLeft) }}</span>
        </p>
      </div>

      <!-- Actions -->
      <div class="actions">
        <button
          class="btn-pay"
          :disabled="!selectedMethod || paying"
          @click="processPayment"
        >
          <span v-if="paying" class="loading-spinner"></span>
          <span>{{ paying ? '处理中...' : '立即支付' }}</span>
        </button>
        
        <button class="btn-cancel" @click="cancelOrder">
          取消订单
        </button>
      </div>

      <!-- Status Messages -->
      <div v-if="paymentStatus" class="status-message" :class="paymentStatus">
        {{ statusMessage }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const orderId = route.params.orderId as string

const order = ref<any>(null)
const selectedMethod = ref('wechat')
const qrCodeUrl = ref('')
const paying = ref(false)
const paymentStatus = ref('')
const statusMessage = ref('')
const timeLeft = ref(1800) // 30 minutes in seconds
let countdownInterval: any = null

const paymentMethods = [
  {
    id: 'wechat',
    name: '微信支付',
    description: '使用微信扫码支付',
    icon: '/icons/wechat-pay.svg'
  },
  {
    id: 'alipay',
    name: '支付宝',
    description: '使用支付宝支付',
    icon: '/icons/alipay.svg'
  }
]

onMounted(async () => {
  await loadOrder()
  startCountdown()
})

onUnmounted(() => {
  if (countdownInterval) {
    clearInterval(countdownInterval)
  }
})

const loadOrder = async () => {
  try {
    const response = await $fetch(`/api/orders/${orderId}`)
    order.value = response.data
    
    // Check if order is expired
    if (order.value.status !== 'pending') {
      router.push(`/orders/${orderId}`)
    }
  } catch (err) {
    console.error('Failed to load order:', err)
  }
}

const processPayment = async () => {
  if (!selectedMethod.value) return
  
  paying.value = true
  try {
    const response = await $fetch('/api/payments', {
      method: 'POST',
      body: {
        order_id: orderId,
        method: selectedMethod.value,
        description: `订单支付 - ${order.value?.order_number}`
      }
    })

    const payment = response.data
    
    if (selectedMethod.value === 'wechat') {
      // For WeChat, show QR code
      if (payment.pay_url) {
        qrCodeUrl.value = payment.pay_url
        // Start polling for payment status
        startPaymentPolling(payment.id)
      }
    } else if (selectedMethod.value === 'alipay') {
      // For Alipay, redirect or open in new window
      if (payment.pay_url) {
        window.open(payment.pay_url, '_blank')
        startPaymentPolling(payment.id)
      }
    }
  } catch (err: any) {
    paymentStatus.value = 'error'
    statusMessage.value = err.message || '支付失败，请重试'
  } finally {
    paying.value = false
  }
}

const startPaymentPolling = (paymentId: string) => {
  const pollInterval = setInterval(async () => {
    try {
      const response = await $fetch(`/api/payments/${paymentId}`)
      const payment = response.data
      
      if (payment.status === 'success') {
        clearInterval(pollInterval)
        paymentStatus.value = 'success'
        statusMessage.value = '支付成功！'
        setTimeout(() => {
          router.push(`/orders/${orderId}/success`)
        }, 1500)
      } else if (payment.status === 'failed') {
        clearInterval(pollInterval)
        paymentStatus.value = 'error'
        statusMessage.value = '支付失败，请重试'
      }
    } catch (err) {
      console.error('Polling error:', err)
    }
  }, 3000) // Poll every 3 seconds

  // Stop polling after 5 minutes
  setTimeout(() => {
    clearInterval(pollInterval)
  }, 300000)
}

const cancelOrder = async () => {
  if (!confirm('确定要取消此订单吗？')) return
  
  try {
    await $fetch(`/api/orders/${orderId}/cancel`, {
      method: 'POST'
    })
    router.push('/orders')
  } catch (err) {
    console.error('Failed to cancel order:', err)
  }
}

const startCountdown = () => {
  countdownInterval = setInterval(() => {
    if (timeLeft.value > 0) {
      timeLeft.value--
    } else {
      clearInterval(countdownInterval)
      // Order expired
      paymentStatus.value = 'error'
      statusMessage.value = '订单已过期'
    }
  }, 1000)
}

const formatPrice = (price: number) => {
  return price?.toLocaleString('zh-CN') || '0'
}

const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}
</script>

<style scoped>
.payment-page {
  @apply min-h-screen bg-gray-100 py-8 px-4;
}

.payment-container {
  @apply max-w-2xl mx-auto bg-white rounded-xl shadow-lg p-8;
}

.payment-header {
  @apply text-center mb-8;
}

.order-summary {
  @apply bg-gray-50 rounded-lg p-6 mb-6;
}

.summary-items {
  @apply space-y-2;
}

.summary-item {
  @apply flex justify-between text-gray-600;
}

.summary-total {
  @apply flex justify-between items-center mt-4 pt-4 border-t border-gray-300;
}

.summary-total .amount {
  @apply text-2xl font-bold text-orange-600;
}

.payment-methods {
  @apply mb-6;
}

.methods-list {
  @apply space-y-3;
}

.method-item {
  @apply flex items-center p-4 border-2 border-gray-200 rounded-lg cursor-pointer transition-all duration-200 hover:border-blue-300;
}

.method-item.selected {
  @apply border-blue-500 bg-blue-50;
}

.method-icon {
  @apply w-12 h-12 mr-4;
}

.method-icon img {
  @apply w-full h-full object-contain;
}

.method-info {
  @apply flex-1;
}

.method-name {
  @apply font-medium text-gray-800;
}

.method-desc {
  @apply text-sm text-gray-500;
}

.check-icon {
  @apply text-blue-500;
}

.qr-section {
  @apply mb-6;
}

.qr-container {
  @apply flex flex-col items-center p-6 bg-gray-50 rounded-lg;
}

.qr-code {
  @apply w-48 h-48 mb-4;
}

.qr-hint {
  @apply text-sm text-gray-600;
}

.countdown-section {
  @apply text-center mb-6;
}

.countdown {
  @apply text-orange-500 font-mono font-bold;
}

.actions {
  @apply space-y-3;
}

.btn-pay {
  @apply w-full flex items-center justify-center gap-2 py-3 bg-orange-500 text-white rounded-lg font-medium hover:bg-orange-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors;
}

.btn-cancel {
  @apply w-full py-3 text-gray-500 hover:text-gray-700 transition-colors;
}

.loading-spinner {
  @apply w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin;
}

.status-message {
  @apply mt-4 p-4 rounded-lg text-center;
}

.status-message.success {
  @apply bg-green-50 text-green-600;
}

.status-message.error {
  @apply bg-red-50 text-red-600;
}
</style>
