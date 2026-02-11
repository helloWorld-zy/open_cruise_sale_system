<template>
  <div class="order-detail-page">
    <div class="container" v-if="order">
      <!-- Header -->
      <div class="page-header">
        <button class="btn-back" @click="goBack">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
          </svg>
          返回
        </button>
        <h1 class="text-2xl font-bold">订单详情</h1>
        <span class="order-status" :class="getStatusClass(order.status)">
          {{ getStatusLabel(order.status) }}
        </span>
      </div>

      <!-- Order Info Card -->
      <div class="info-card">
        <h2 class="card-title">订单信息</h2>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">订单号</span>
            <span class="value">{{ order.order_number }}</span>
          </div>
          <div class="info-item">
            <span class="label">下单时间</span>
            <span class="value">{{ formatDateTime(order.created_at) }}</span>
          </div>
          <div class="info-item">
            <span class="label">订单状态</span>
            <span class="value">{{ getStatusLabel(order.status) }}</span>
          </div>
          <div class="info-item">
            <span class="label">支付状态</span>
            <span class="value">{{ order.payment_status === 'paid' ? '已支付' : '未支付' }}</span>
          </div>
        </div>
      </div>

      <!-- Voyage Info -->
      <div class="info-card">
        <h2 class="card-title">航次信息</h2>
        <div class="voyage-detail">
          <h3 class="route-name">{{ order.voyage?.route?.name }}</h3>
          <div class="voyage-meta">
            <p><span>出发日期:</span> {{ order.voyage?.departure_date }} {{ order.voyage?.departure_time }}</p>
            <p><span>返回日期:</span> {{ order.voyage?.arrival_date }} {{ order.voyage?.arrival_time }}</p>
            <p><span>航程:</span> {{ order.voyage?.route?.duration_days }} 天</p>
            <p><span>邮轮:</span> {{ order.cruise?.name_cn }}</p>
          </div>
        </div>
      </div>

      <!-- Cabin Info -->
      <div class="info-card">
        <h2 class="card-title">舱房信息</h2>
        <div class="cabin-list">
          <div
            v-for="(item, index) in order.items"
            :key="index"
            class="cabin-item"
          >
            <div class="cabin-header">
              <span class="cabin-name">{{ item.cabin_type?.name_cn }}</span>
              <span class="cabin-status" :class="item.status">{{ item.status }}</span>
            </div>
            <div class="cabin-meta">
              <span>舱号: {{ item.cabin_number || '待分配' }}</span>
              <span>{{ item.adult_count }} 成人</span>
              <span v-if="item.child_count">{{ item.child_count }} 儿童</span>
            </div>
            <div class="cabin-price">¥{{ formatPrice(item.subtotal) }}</div>
          </div>
        </div>
      </div>

      <!-- Passenger Info -->
      <div class="info-card">
        <h2 class="card-title">乘客信息</h2>
        <div class="passenger-list">
          <div
            v-for="(passenger, index) in order.passengers"
            :key="index"
            class="passenger-item"
          >
            <div class="passenger-index">{{ index + 1 }}</div>
            <div class="passenger-info">
              <div class="passenger-name">
                {{ passenger.surname }} {{ passenger.given_name }} / {{ passenger.name }}
              </div>
              <div class="passenger-meta">
                <span>{{ getPassengerTypeLabel(passenger.passenger_type) }}</span>
                <span>{{ passenger.gender === 'male' ? '男' : '女' }}</span>
                <span>{{ formatDate(passenger.birth_date) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Contact Info -->
      <div class="info-card">
        <h2 class="card-title">联系人信息</h2>
        <div class="contact-info">
          <p><span>姓名:</span> {{ order.contact_name }}</p>
          <p><span>手机:</span> {{ order.contact_phone }}</p>
          <p v-if="order.contact_email"><span>邮箱:</span> {{ order.contact_email }}</p>
        </div>
      </div>

      <!-- Price Summary -->
      <div class="info-card price-card">
        <h2 class="card-title">费用明细</h2>
        <div class="price-breakdown">
          <div class="price-row">
            <span>舱房费用</span>
            <span>¥{{ formatPrice(calculateCabinTotal()) }}</span>
          </div>
          <div class="price-row" v-if="order.discount_amount > 0">
            <span>优惠金额</span>
            <span class="discount">-¥{{ formatPrice(order.discount_amount) }}</span>
          </div>
          <div class="price-row total">
            <span>订单总额</span>
            <span class="total-amount">¥{{ formatPrice(order.total_amount) }}</span>
          </div>
        </div>
      </div>

      <!-- Payment Info -->
      <div class="info-card" v-if="order.payments && order.payments.length > 0">
        <h2 class="card-title">支付记录</h2>
        <div class="payment-list">
          <div
            v-for="(payment, index) in order.payments"
            :key="index"
            class="payment-item"
          >
            <div class="payment-info">
              <span class="payment-no">支付单号: {{ payment.payment_no }}</span>
              <span class="payment-method">{{ getPaymentMethodLabel(payment.payment_method) }}</span>
            </div>
            <div class="payment-amount">¥{{ formatPrice(payment.amount) }}</div>
            <div class="payment-status" :class="payment.status">{{ payment.status }}</div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="actions">
        <button
          v-if="order.status === 'pending'"
          class="btn-primary"
          @click="goToPayment"
        >
          立即支付
        </button>
        <button
          v-if="canRefund(order)"
          class="btn-secondary"
          @click="showRefundModal = true"
        >
          申请退款
        </button>
        <button
          v-if="order.status === 'pending'"
          class="btn-danger"
          @click="cancelOrder"
        >
          取消订单
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-else-if="loading" class="loading-state">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
      <p class="text-gray-500 mt-4 text-center">加载中...</p>
    </div>

    <!-- Refund Modal -->
    <RefundRequestModal
      v-if="showRefundModal"
      :order="order"
      @close="showRefundModal = false"
      @submit="handleRefund"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const orderId = route.params.id as string

const order = ref<any>(null)
const loading = ref(false)
const showRefundModal = ref(false)

onMounted(() => {
  loadOrder()
})

const loadOrder = async () => {
  loading.value = true
  try {
    const response = await $fetch(`/api/orders/${orderId}/detail`)
    order.value = response.data
  } catch (err) {
    console.error('Failed to load order:', err)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/orders')
}

const goToPayment = () => {
  router.push(`/payment/${orderId}`)
}

const cancelOrder = async () => {
  if (!confirm('确定要取消此订单吗？')) return
  
  try {
    await $fetch(`/api/orders/${orderId}/cancel`, { method: 'POST' })
    loadOrder()
  } catch (err: any) {
    alert(err.message || '取消订单失败')
  }
}

const handleRefund = async (refundData: any) => {
  try {
    await $fetch('/api/refunds', {
      method: 'POST',
      body: refundData
    })
    showRefundModal.value = false
    loadOrder()
    alert('退款申请已提交，等待审核')
  } catch (err: any) {
    alert(err.message || '申请退款失败')
  }
}

const calculateCabinTotal = () => {
  if (!order.value?.items) return 0
  return order.value.items.reduce((sum: number, item: any) => sum + item.subtotal, 0)
}

const canRefund = (order: any) => {
  return order?.status === 'paid' || order?.status === 'confirmed'
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    pending: 'status-pending',
    paid: 'status-paid',
    confirmed: 'status-confirmed',
    completed: 'status-completed',
    cancelled: 'status-cancelled',
    refunded: 'status-refunded'
  }
  return classes[status] || ''
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    confirmed: '已确认',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return labels[status] || status
}

const getPassengerTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    adult: '成人',
    child: '儿童',
    infant: '婴儿'
  }
  return labels[type] || type
}

const getPaymentMethodLabel = (method: string) => {
  const labels: Record<string, string> = {
    wechat: '微信支付',
    alipay: '支付宝',
    card: '银行卡'
  }
  return labels[method] || method
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const formatPrice = (price: number) => {
  return price?.toLocaleString('zh-CN') || '0'
}
</script>

<style scoped>
.order-detail-page {
  @apply min-h-screen bg-gray-50 py-8;
}

.container {
  @apply max-w-3xl mx-auto px-4 space-y-4;
}

.page-header {
  @apply flex items-center gap-4 mb-6;
}

.btn-back {
  @apply flex items-center gap-2 text-gray-600 hover:text-gray-800;
}

.order-status {
  @apply px-3 py-1 rounded-full text-sm font-medium ml-auto;
}

.info-card {
  @apply bg-white rounded-xl p-6 shadow-sm;
}

.card-title {
  @apply text-lg font-bold text-gray-800 mb-4 pb-2 border-b;
}

.info-grid {
  @apply grid grid-cols-2 gap-4;
}

.info-item {
  @apply flex flex-col;
}

.info-item .label {
  @apply text-sm text-gray-500 mb-1;
}

.info-item .value {
  @apply font-medium text-gray-800;
}

.voyage-detail .route-name {
  @apply text-xl font-bold text-gray-800 mb-3;
}

.voyage-meta p {
  @apply text-gray-600 mb-2;
}

.voyage-meta span {
  @apply font-medium text-gray-700;
}

.cabin-list {
  @apply space-y-4;
}

.cabin-item {
  @apply p-4 bg-gray-50 rounded-lg;
}

.cabin-header {
  @apply flex justify-between items-center mb-2;
}

.cabin-name {
  @apply font-medium text-gray-800;
}

.cabin-status {
  @apply text-sm;
}

.cabin-status.confirmed {
  @apply text-green-600;
}

.cabin-meta {
  @apply flex gap-4 text-sm text-gray-500 mb-2;
}

.cabin-price {
  @apply font-bold text-orange-500;
}

.passenger-list {
  @apply space-y-3;
}

.passenger-item {
  @apply flex items-center gap-3 p-3 bg-gray-50 rounded-lg;
}

.passenger-index {
  @apply w-8 h-8 bg-blue-500 text-white rounded-full flex items-center justify-center text-sm font-medium;
}

.passenger-name {
  @apply font-medium text-gray-800;
}

.passenger-meta {
  @apply text-sm text-gray-500 mt-1;
}

.passenger-meta span {
  @apply mr-3;
}

.contact-info p {
  @apply mb-2 text-gray-600;
}

.contact-info span {
  @apply font-medium text-gray-700 mr-2;
}

.price-card {
  @apply bg-orange-50;
}

.price-breakdown {
  @apply space-y-3;
}

.price-row {
  @apply flex justify-between text-gray-600;
}

.price-row.discount span:last-child {
  @apply text-green-600;
}

.price-row.total {
  @apply pt-3 border-t border-orange-200 text-lg font-bold;
}

.total-amount {
  @apply text-orange-600;
}

.payment-list {
  @apply space-y-3;
}

.payment-item {
  @apply flex justify-between items-center p-3 bg-gray-50 rounded-lg;
}

.payment-info {
  @apply flex flex-col;
}

.payment-no {
  @apply text-sm font-medium text-gray-700;
}

.payment-method {
  @apply text-xs text-gray-500 mt-1;
}

.payment-amount {
  @apply font-bold text-gray-800;
}

.payment-status {
  @apply text-xs px-2 py-1 rounded-full ml-2;
}

.payment-status.success {
  @apply bg-green-100 text-green-600;
}

.actions {
  @apply flex gap-3 pt-4;
}

.btn-primary {
  @apply flex-1 py-3 bg-blue-500 text-white font-medium rounded-lg hover:bg-blue-600 transition-colors;
}

.btn-secondary {
  @apply flex-1 py-3 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 transition-colors;
}

.btn-danger {
  @apply flex-1 py-3 bg-red-500 text-white font-medium rounded-lg hover:bg-red-600 transition-colors;
}

.loading-state {
  @apply py-12;
}
</style>
