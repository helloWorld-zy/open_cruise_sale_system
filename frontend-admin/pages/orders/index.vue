<template>
  <div class="admin-orders-page">
    <div class="page-header">
      <h1 class="text-2xl font-bold">订单管理</h1>
      <p class="text-gray-500 mt-2">管理所有用户订单和退款申请</p>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid" v-if="statistics">
      <div class="stat-card">
        <div class="stat-value">{{ statistics.total_orders || 0 }}</div>
        <div class="stat-label">全部订单</div>
      </div>
      <div class="stat-card warning">
        <div class="stat-value">{{ statistics.pending_orders || 0 }}</div>
        <div class="stat-label">待支付</div>
      </div>
      <div class="stat-card info">
        <div class="stat-value">{{ statistics.paid_orders || 0 }}</div>
        <div class="stat-label">已支付</div>
      </div>
      <div class="stat-card success">
        <div class="stat-value">{{ statistics.completed_orders || 0 }}</div>
        <div class="stat-label">已完成</div>
      </div>
      <div class="stat-card danger">
        <div class="stat-value">{{ statistics.pending_refunds || 0 }}</div>
        <div class="stat-label">待处理退款</div>
      </div>
    </div>

    <!-- Filter Section -->
    <div class="filter-section">
      <div class="filter-row">
        <input
          v-model="filters.orderNumber"
          type="text"
          placeholder="搜索订单号"
          class="filter-input"
        />
        <select v-model="filters.status" class="filter-select">
          <option value="">全部状态</option>
          <option value="pending">待支付</option>
          <option value="paid">已支付</option>
          <option value="confirmed">已确认</option>
          <option value="completed">已完成</option>
          <option value="cancelled">已取消</option>
          <option value="refunded">已退款</option>
        </select>
        <input
          v-model="filters.dateFrom"
          type="date"
          class="filter-input"
          placeholder="开始日期"
        />
        <input
          v-model="filters.dateTo"
          type="date"
          class="filter-input"
          placeholder="结束日期"
        />
        <button class="btn-search" @click="loadOrders">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
          </svg>
          搜索
        </button>
        <button class="btn-reset" @click="resetFilters">
          重置
        </button>
      </div>
    </div>

    <!-- Orders Table -->
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>订单号</th>
            <th>用户信息</th>
            <th>航线/航次</th>
            <th>订单金额</th>
            <th>支付状态</th>
            <th>订单状态</th>
            <th>下单时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in orders" :key="order.id">
            <td class="order-number">{{ order.order_number }}</td>
            <td>
              <div class="user-info">
                <div class="user-name">{{ order.contact_name }}</div>
                <div class="user-phone">{{ order.contact_phone }}</div>
              </div>
            </td>
            <td>
              <div class="voyage-info">
                <div class="route-name">{{ order.voyage?.route?.name }}</div>
                <div class="voyage-date">{{ order.voyage?.departure_date }}</div>
              </div>
            </td>
            <td class="order-amount">¥{{ formatPrice(order.total_amount) }}</td>
            <td>
              <span class="badge" :class="getPaymentStatusClass(order.payment_status)">
                {{ getPaymentStatusLabel(order.payment_status) }}
              </span>
            </td>
            <td>
              <span class="badge" :class="getOrderStatusClass(order.status)">
                {{ getOrderStatusLabel(order.status) }}
              </span>
            </td>
            <td class="text-sm text-gray-500">{{ formatDate(order.created_at) }}</td>
            <td class="actions">
              <button class="btn-action" @click="viewDetail(order.id)">
                查看
              </button>
              <button 
                v-if="order.status === 'pending'"
                class="btn-action warning"
                @click="cancelOrder(order.id)"
              >
                取消
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="pagination" v-if="totalPages > 1">
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
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const orders = ref<any[]>([])
const statistics = ref<any>({})
const loading = ref(false)
const currentPage = ref(1)
const totalPages = ref(1)

const filters = ref({
  orderNumber: '',
  status: '',
  dateFrom: '',
  dateTo: ''
})

onMounted(() => {
  loadOrders()
  loadStatistics()
})

const loadOrders = async () => {
  loading.value = true
  try {
    const response = await $fetch('/api/admin/orders', {
      params: {
        page: currentPage.value,
        page_size: 20,
        order_number: filters.value.orderNumber,
        status: filters.value.status,
        date_from: filters.value.dateFrom,
        date_to: filters.value.dateTo
      }
    })
    orders.value = response.data || []
    totalPages.value = Math.ceil((response.pagination?.total || 0) / 20)
  } catch (err) {
    console.error('Failed to load orders:', err)
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
  try {
    const response = await $fetch('/api/admin/orders/statistics')
    statistics.value = response.data
  } catch (err) {
    console.error('Failed to load statistics:', err)
  }
}

const resetFilters = () => {
  filters.value = {
    orderNumber: '',
    status: '',
    dateFrom: '',
    dateTo: ''
  }
  loadOrders()
}

const viewDetail = (orderId: string) => {
  router.push(`/admin/orders/${orderId}`)
}

const cancelOrder = async (orderId: string) => {
  if (!confirm('确定要取消此订单吗？')) return
  
  try {
    await $fetch(`/api/orders/${orderId}/cancel`, { method: 'POST' })
    loadOrders()
  } catch (err: any) {
    alert(err.message || '取消订单失败')
  }
}

const getOrderStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    pending: 'bg-orange-100 text-orange-600',
    paid: 'bg-blue-100 text-blue-600',
    confirmed: 'bg-indigo-100 text-indigo-600',
    completed: 'bg-green-100 text-green-600',
    cancelled: 'bg-gray-100 text-gray-600',
    refunded: 'bg-red-100 text-red-600'
  }
  return classes[status] || 'bg-gray-100 text-gray-600'
}

const getOrderStatusLabel = (status: string) => {
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

const getPaymentStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    unpaid: 'bg-orange-100 text-orange-600',
    partial: 'bg-yellow-100 text-yellow-600',
    paid: 'bg-green-100 text-green-600',
    refunded: 'bg-red-100 text-red-600'
  }
  return classes[status] || 'bg-gray-100 text-gray-600'
}

const getPaymentStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    unpaid: '未支付',
    partial: '部分支付',
    paid: '已支付',
    refunded: '已退款'
  }
  return labels[status] || status
}

const formatPrice = (price: number) => {
  return price?.toLocaleString('zh-CN') || '0'
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.admin-orders-page {
  @apply p-6;
}

.page-header {
  @apply mb-6;
}

.stats-grid {
  @apply grid grid-cols-5 gap-4 mb-6;
}

.stat-card {
  @apply bg-white rounded-xl p-4 shadow-sm;
}

.stat-card.warning {
  @apply bg-orange-50;
}

.stat-card.info {
  @apply bg-blue-50;
}

.stat-card.success {
  @apply bg-green-50;
}

.stat-card.danger {
  @apply bg-red-50;
}

.stat-value {
  @apply text-2xl font-bold text-gray-800;
}

.stat-card.warning .stat-value {
  @apply text-orange-600;
}

.stat-card.info .stat-value {
  @apply text-blue-600;
}

.stat-card.success .stat-value {
  @apply text-green-600;
}

.stat-card.danger .stat-value {
  @apply text-red-600;
}

.stat-label {
  @apply text-sm text-gray-500 mt-1;
}

.filter-section {
  @apply bg-white rounded-xl p-4 shadow-sm mb-6;
}

.filter-row {
  @apply flex gap-3 flex-wrap;
}

.filter-input,
.filter-select {
  @apply px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.filter-input {
  @apply w-48;
}

.filter-select {
  @apply w-40;
}

.btn-search,
.btn-reset {
  @apply px-4 py-2 rounded-lg font-medium flex items-center gap-2;
}

.btn-search {
  @apply bg-blue-500 text-white hover:bg-blue-600;
}

.btn-reset {
  @apply border border-gray-300 text-gray-700 hover:bg-gray-50;
}

.table-container {
  @apply bg-white rounded-xl shadow-sm overflow-hidden;
}

.data-table {
  @apply w-full;
}

.data-table th {
  @apply px-4 py-3 bg-gray-50 text-left text-sm font-medium text-gray-700 border-b;
}

.data-table td {
  @apply px-4 py-3 border-b border-gray-100;
}

.order-number {
  @apply font-mono text-sm text-gray-600;
}

.user-info .user-name {
  @apply font-medium text-gray-800;
}

.user-info .user-phone {
  @apply text-sm text-gray-500;
}

.voyage-info .route-name {
  @apply font-medium text-gray-800;
}

.voyage-info .voyage-date {
  @apply text-sm text-gray-500;
}

.order-amount {
  @apply font-bold text-gray-800;
}

.badge {
  @apply px-2 py-1 rounded-full text-xs font-medium;
}

.actions {
  @apply flex gap-2;
}

.btn-action {
  @apply px-3 py-1 text-sm text-blue-500 hover:bg-blue-50 rounded-lg transition-colors;
}

.btn-action.warning {
  @apply text-orange-500 hover:bg-orange-50;
}

.pagination {
  @apply flex justify-center items-center gap-4 mt-6;
}

.page-btn {
  @apply px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed;
}

.page-info {
  @apply text-sm text-gray-600;
}
</style>
