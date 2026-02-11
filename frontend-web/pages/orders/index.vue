<template>
  <div class="my-orders-page">
    <div class="container">
      <!-- Header -->
      <div class="page-header">
        <h1 class="text-2xl font-bold">我的订单</h1>
        <p class="text-gray-500 mt-2">查看和管理您的邮轮预订</p>
      </div>

      <!-- Statistics -->
      <div class="statistics-grid" v-if="statistics">
        <div class="stat-card">
          <div class="stat-value">{{ statistics.total_orders || 0 }}</div>
          <div class="stat-label">全部订单</div>
        </div>
        <div class="stat-card pending">
          <div class="stat-value">{{ statistics.pending_orders || 0 }}</div>
          <div class="stat-label">待支付</div>
        </div>
        <div class="stat-card paid">
          <div class="stat-value">{{ statistics.paid_orders || 0 }}</div>
          <div class="stat-label">待出行</div>
        </div>
        <div class="stat-card completed">
          <div class="stat-value">{{ statistics.completed_orders || 0 }}</div>
          <div class="stat-label">已完成</div>
        </div>
      </div>

      <!-- Filter Tabs -->
      <div class="filter-tabs">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          class="tab-button"
          :class="{ active: activeTab === tab.value }"
          @click="activeTab = tab.value"
        >
          {{ tab.label }}
          <span v-if="tab.count" class="tab-count">({{ tab.count }})</span>
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
        <p class="text-gray-500 mt-4 text-center">加载订单中...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredOrders.length === 0" class="empty-state">
        <svg class="w-24 h-24 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
        </svg>
        <p class="text-gray-500 text-center">暂无{{ activeTabLabel }}订单</p>
        <NuxtLink to="/cruises" class="btn-primary mt-4 inline-block">
          去预订邮轮
        </NuxtLink>
      </div>

      <!-- Order List -->
      <div v-else class="order-list">
        <div
          v-for="order in filteredOrders"
          :key="order.id"
          class="order-card"
          @click="goToDetail(order.id)"
        >
          <div class="order-header">
            <div class="order-info">
              <span class="order-number">订单号: {{ order.order_number }}</span>
              <span class="order-date">{{ formatDate(order.created_at) }}</span>
            </div>
            <div class="order-status" :class="getStatusClass(order.status)">
              {{ getStatusLabel(order.status) }}
            </div>
          </div>

          <div class="order-body">
            <div class="cruise-info">
              <h3 class="cruise-name">{{ order.voyage?.route?.name || '未知航线' }}</h3>
              <p class="cruise-detail">
                {{ order.voyage?.departure_date }} 出发 | {{ order.voyage?.route?.duration_days }}天航程
              </p>
              <p class="cruise-ship">{{ order.cruise?.name_cn }}</p>
            </div>
            <div class="cabin-info">
              <span class="cabin-count">{{ order.cabin_count }} 间舱房</span>
              <span class="passenger-count">{{ order.passenger_count }} 位乘客</span>
            </div>
          </div>

          <div class="order-footer">
            <div class="order-amount">
              <span class="amount-label">订单金额</span>
              <span class="amount-value">¥{{ formatPrice(order.total_amount) }}</span>
            </div>
            <div class="order-actions">
              <button
                v-if="order.status === 'pending'"
                class="btn-pay"
                @click.stop="goToPayment(order.id)"
              >
                立即支付
              </button>
              <button
                v-if="canRefund(order)"
                class="btn-refund"
                @click.stop="showRefundModal(order)"
              >
                申请退款
              </button>
              <button class="btn-detail">
                查看详情
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button
          class="page-btn"
          :disabled="currentPage === 1"
          @click="currentPage--"
        >
          上一页
        </button>
        <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
        <button
          class="page-btn"
          :disabled="currentPage === totalPages"
          @click="currentPage++"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- Refund Modal -->
    <RefundRequestModal
      v-if="showRefundModal"
      :order="selectedOrder"
      @close="showRefundModal = false"
      @submit="handleRefundSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const orders = ref<any[]>([])
const statistics = ref<any>(null)
const activeTab = ref('all')
const currentPage = ref(1)
const totalPages = ref(1)

const showRefundModal = ref(false)
const selectedOrder = ref<any>(null)

const tabs = [
  { label: '全部', value: 'all', count: 0 },
  { label: '待支付', value: 'pending', count: 0 },
  { label: '待出行', value: 'confirmed', count: 0 },
  { label: '已完成', value: 'completed', count: 0 },
  { label: '已取消', value: 'cancelled', count: 0 }
]

const activeTabLabel = computed(() => {
  const tab = tabs.find(t => t.value === activeTab.value)
  return tab?.label || ''
})

const filteredOrders = computed(() => {
  if (activeTab.value === 'all') {
    return orders.value
  }
  return orders.value.filter(order => order.status === activeTab.value)
})

onMounted(() => {
  loadOrders()
  loadStatistics()
})

const loadOrders = async () => {
  loading.value = true
  try {
    const response = await $fetch('/api/orders/my', {
      params: {
        page: currentPage.value,
        page_size: 10
      }
    })
    orders.value = response.data || []
    totalPages.value = Math.ceil((response.pagination?.total || 0) / 10)
  } catch (err) {
    console.error('Failed to load orders:', err)
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
  try {
    const response = await $fetch('/api/orders/statistics')
    statistics.value = response.data
  } catch (err) {
    console.error('Failed to load statistics:', err)
  }
}

const goToDetail = (orderId: string) => {
  router.push(`/orders/${orderId}`)
}

const goToPayment = (orderId: string) => {
  router.push(`/payment/${orderId}`)
}

const showRefundModal = (order: any) => {
  selectedOrder.value = order
  showRefundModal.value = true
}

const handleRefundSubmit = async (refundData: any) => {
  try {
    await $fetch('/api/refunds', {
      method: 'POST',
      body: refundData
    })
    showRefundModal.value = false
    // Refresh orders
    loadOrders()
  } catch (err: any) {
    alert(err.message || '申请退款失败')
  }
}

const canRefund = (order: any) => {
  return order.status === 'paid' || order.status === 'confirmed'
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

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

const formatPrice = (price: number) => {
  return price?.toLocaleString('zh-CN') || '0'
}
</script>

<style scoped>
.my-orders-page {
  @apply min-h-screen bg-gray-50 py-8;
}

.container {
  @apply max-w-4xl mx-auto px-4;
}

.page-header {
  @apply mb-8;
}

.statistics-grid {
  @apply grid grid-cols-4 gap-4 mb-8;
}

.stat-card {
  @apply bg-white rounded-xl p-4 text-center shadow-sm;
}

.stat-card.pending {
  @apply bg-orange-50;
}

.stat-card.paid {
  @apply bg-blue-50;
}

.stat-card.completed {
  @apply bg-green-50;
}

.stat-value {
  @apply text-2xl font-bold text-gray-800;
}

.stat-card.pending .stat-value {
  @apply text-orange-500;
}

.stat-card.paid .stat-value {
  @apply text-blue-500;
}

.stat-card.completed .stat-value {
  @apply text-green-500;
}

.stat-label {
  @apply text-sm text-gray-500 mt-1;
}

.filter-tabs {
  @apply flex gap-2 mb-6 overflow-x-auto pb-2;
}

.tab-button {
  @apply px-4 py-2 rounded-lg text-sm font-medium text-gray-600 bg-white whitespace-nowrap transition-colors;
}

.tab-button.active {
  @apply bg-blue-500 text-white;
}

.tab-count {
  @apply ml-1 text-xs;
}

.order-list {
  @apply space-y-4;
}

.order-card {
  @apply bg-white rounded-xl shadow-sm overflow-hidden cursor-pointer transition-shadow hover:shadow-md;
}

.order-header {
  @apply flex justify-between items-center p-4 bg-gray-50 border-b;
}

.order-info {
  @apply flex flex-col gap-1;
}

.order-number {
  @apply text-sm font-medium text-gray-800;
}

.order-date {
  @apply text-xs text-gray-500;
}

.order-status {
  @apply px-3 py-1 rounded-full text-xs font-medium;
}

.status-pending {
  @apply bg-orange-100 text-orange-600;
}

.status-paid {
  @apply bg-blue-100 text-blue-600;
}

.status-confirmed {
  @apply bg-indigo-100 text-indigo-600;
}

.status-completed {
  @apply bg-green-100 text-green-600;
}

.status-cancelled,
.status-refunded {
  @apply bg-gray-100 text-gray-600;
}

.order-body {
  @apply p-4;
}

.cruise-name {
  @apply font-bold text-lg text-gray-800 mb-2;
}

.cruise-detail,
.cruise-ship {
  @apply text-sm text-gray-500 mb-1;
}

.cabin-info {
  @apply flex gap-4 mt-3 text-sm text-gray-600;
}

.order-footer {
  @apply flex justify-between items-center p-4 border-t bg-gray-50;
}

.order-amount {
  @apply flex flex-col;
}

.amount-label {
  @apply text-xs text-gray-500;
}

.amount-value {
  @apply text-xl font-bold text-orange-500;
}

.order-actions {
  @apply flex gap-2;
}

.btn-pay {
  @apply px-4 py-2 bg-orange-500 text-white text-sm font-medium rounded-lg hover:bg-orange-600 transition-colors;
}

.btn-refund {
  @apply px-4 py-2 border border-gray-300 text-gray-600 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors;
}

.btn-detail {
  @apply px-4 py-2 text-blue-500 text-sm font-medium hover:bg-blue-50 rounded-lg transition-colors;
}

.btn-primary {
  @apply px-6 py-3 bg-blue-500 text-white font-medium rounded-lg hover:bg-blue-600 transition-colors;
}

.pagination {
  @apply flex justify-center items-center gap-4 mt-8;
}

.page-btn {
  @apply px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed;
}

.page-info {
  @apply text-sm text-gray-600;
}

.loading-state,
.empty-state {
  @apply py-12;
}
</style>
