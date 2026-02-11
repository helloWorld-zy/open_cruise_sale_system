<template>
  <div class="admin-refunds">
    <div class="page-header">
      <h1 class="text-2xl font-bold">退款管理</h1>
    </div>

    <div class="filter-section">
      <select v-model="filters.status" class="filter-select" @change="loadRefunds">
        <option value="">全部状态</option>
        <option value="pending">待审核</option>
        <option value="approved">已批准</option>
        <option value="rejected">已拒绝</option>
        <option value="processing">处理中</option>
        <option value="completed">已完成</option>
      </select>
    </div>

    <table class="data-table">
      <thead>
        <tr>
          <th>退款单号</th>
          <th>订单号</th>
          <th>退款金额</th>
          <th>退款原因</th>
          <th>状态</th>
          <th>申请时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="refund in refunds" :key="refund.id">
          <td>{{ refund.id.slice(0, 8) }}</td>
          <td>{{ refund.order?.order_number }}</td>
          <td class="amount">¥{{ formatPrice(refund.refund_amount) }}</td>
          <td>{{ refund.refund_reason }}</td>
          <td>
            <span class="badge" :class="getStatusClass(refund.status)">
              {{ getStatusLabel(refund.status) }}
            </span>
          </td>
          <td>{{ formatDate(refund.requested_at) }}</td>
          <td class="actions">
            <button v-if="refund.status === 'pending'" class="btn-approve" @click="approve(refund.id)">
              批准
            </button>
            <button v-if="refund.status === 'pending'" class="btn-reject" @click="reject(refund.id)">
              拒绝
            </button>
            <button v-if="refund.status === 'approved'" class="btn-process" @click="process(refund.id)">
              处理退款
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
const refunds = ref([])
const filters = ref({ status: '' })

const loadRefunds = async () => {
  const res = await $fetch('/api/admin/refunds', { params: filters.value })
  refunds.value = res.data
}

const approve = async (id) => {
  const note = prompt('请输入审核备注:')
  await $fetch(`/api/admin/refunds/${id}/approve`, {
    method: 'POST',
    body: { note }
  })
  loadRefunds()
}

const reject = async (id) => {
  const note = prompt('请输入拒绝原因:')
  if (!note) return
  await $fetch(`/api/admin/refunds/${id}/reject`, {
    method: 'POST',
    body: { note }
  })
  loadRefunds()
}

const process = async (id) => {
  if (!confirm('确定要处理此退款吗？')) return
  await $fetch(`/api/admin/refunds/${id}/process`, { method: 'POST' })
  loadRefunds()
}

const getStatusClass = (s) => ({
  pending: 'bg-yellow-100 text-yellow-600',
  approved: 'bg-blue-100 text-blue-600',
  rejected: 'bg-red-100 text-red-600',
  processing: 'bg-orange-100 text-orange-600',
  completed: 'bg-green-100 text-green-600'
})[s] || 'bg-gray-100'

const getStatusLabel = (s) => ({
  pending: '待审核', approved: '已批准', rejected: '已拒绝',
  processing: '处理中', completed: '已完成'
})[s] || s

const formatPrice = (p) => p?.toLocaleString('zh-CN') || '0'
const formatDate = (d) => d ? new Date(d).toLocaleDateString('zh-CN') : '-'

onMounted(loadRefunds)
</script>

<style scoped>
.admin-refunds { @apply p-6; }
.page-header { @apply mb-6; }
.filter-section { @apply mb-4; }
.filter-select { @apply px-3 py-2 border border-gray-300 rounded-lg; }
.data-table { @apply w-full bg-white rounded-xl shadow-sm; }
.data-table th { @apply px-4 py-3 bg-gray-50 text-left text-sm font-medium border-b; }
.data-table td { @apply px-4 py-3 border-b; }
.data-table .amount { @apply font-bold text-orange-500; }
.badge { @apply px-2 py-1 rounded-full text-xs font-medium; }
.actions { @apply flex gap-2; }
.btn-approve { @apply px-3 py-1 bg-green-500 text-white text-sm rounded-lg; }
.btn-reject { @apply px-3 py-1 bg-red-500 text-white text-sm rounded-lg; }
.btn-process { @apply px-3 py-1 bg-blue-500 text-white text-sm rounded-lg; }
</style>
