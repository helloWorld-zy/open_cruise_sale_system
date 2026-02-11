<template>
  <div class="admin-order-detail">
    <div class="page-header">
      <button class="btn-back" @click="goBack">返回列表</button>
      <h1 class="text-2xl font-bold">订单详情</h1>
      <div class="header-actions">
        <select v-model="order.status" class="status-select" @change="updateStatus">
          <option value="pending">待支付</option>
          <option value="paid">已支付</option>
          <option value="confirmed">已确认</option>
          <option value="completed">已完成</option>
          <option value="cancelled">已取消</option>
          <option value="refunded">已退款</option>
        </select>
      </div>
    </div>

    <div v-if="order" class="detail-grid">
      <div class="info-card">
        <h3>订单信息</h3>
        <p><strong>订单号:</strong> {{ order.order_number }}</p>
        <p><strong>下单时间:</strong> {{ formatDateTime(order.created_at) }}</p>
        <p><strong>订单金额:</strong> ¥{{ formatPrice(order.total_amount) }}</p>
        <p><strong>支付状态:</strong> {{ order.payment_status }}</p>
      </div>

      <div class="info-card">
        <h3>用户信息</h3>
        <p><strong>联系人:</strong> {{ order.contact_name }}</p>
        <p><strong>手机:</strong> {{ order.contact_phone }}</p>
        <p><strong>邮箱:</strong> {{ order.contact_email || '-' }}</p>
      </div>

      <div class="info-card wide">
        <h3>航次信息</h3>
        <p><strong>航线:</strong> {{ order.voyage?.route?.name }}</p>
        <p><strong>出发日期:</strong> {{ order.voyage?.departure_date }}</p>
        <p><strong>邮轮:</strong> {{ order.cruise?.name_cn }}</p>
      </div>

      <div class="info-card wide">
        <h3>舱房信息</h3>
        <table class="data-table">
          <tr v-for="item in order.items" :key="item.id">
            <td>{{ item.cabin_type?.name_cn }}</td>
            <td>{{ item.adult_count }}成人 {{ item.child_count }}儿童</td>
            <td>¥{{ formatPrice(item.subtotal) }}</td>
          </tr>
        </table>
      </div>

      <div class="info-card wide">
        <h3>乘客信息</h3>
        <table class="data-table">
          <tr v-for="p in order.passengers" :key="p.id">
            <td>{{ p.name }} ({{ p.surname }} {{ p.given_name }})</td>
            <td>{{ p.passenger_type }}</td>
            <td>{{ p.gender === 'male' ? '男' : '女' }}</td>
            <td>{{ formatDate(p.birth_date) }}</td>
          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const orderId = route.params.id
const order = ref({})

const loadOrder = async () => {
  const res = await $fetch(`/api/admin/orders/${orderId}`)
  order.value = res.data
}

const updateStatus = async () => {
  await $fetch(`/api/admin/orders/${orderId}/status`, {
    method: 'PUT',
    body: { status: order.value.status }
  })
}

const goBack = () => router.push('/admin/orders')
const formatPrice = (p) => p?.toLocaleString('zh-CN') || '0'
const formatDate = (d) => d ? new Date(d).toLocaleDateString('zh-CN') : '-'
const formatDateTime = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(loadOrder)
</script>

<style scoped>
.admin-order-detail { @apply p-6; }
.page-header { @apply flex items-center gap-4 mb-6; }
.btn-back { @apply px-4 py-2 border border-gray-300 rounded-lg; }
.status-select { @apply px-3 py-2 border border-gray-300 rounded-lg; }
.detail-grid { @apply grid grid-cols-2 gap-4; }
.info-card { @apply bg-white rounded-xl p-4 shadow-sm; }
.info-card.wide { @apply col-span-2; }
.info-card h3 { @apply font-bold mb-3 pb-2 border-b; }
.info-card p { @apply mb-2; }
.data-table { @apply w-full; }
.data-table td { @apply py-2 border-b; }
</style>
