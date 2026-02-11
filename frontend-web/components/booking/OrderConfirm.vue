<template>
  <div class="order-confirm">
    <h3 class="text-xl font-semibold mb-6">确认订单信息</h3>

    <!-- Voyage Info -->
    <div class="info-section">
      <h4 class="section-title">航次信息</h4>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">出发日期</span>
          <span class="value">{{ bookingData.voyage?.departure_date }}</span>
        </div>
        <div class="info-item">
          <span class="label">返回日期</span>
          <span class="value">{{ bookingData.voyage?.arrival_date }}</span>
        </div>
        <div class="info-item">
          <span class="label">航线</span>
          <span class="value">{{ bookingData.voyage?.route?.name }}</span>
        </div>
        <div class="info-item">
          <span class="label">航程</span>
          <span class="value">{{ bookingData.voyage?.route?.duration_days }} 天</span>
        </div>
      </div>
    </div>

    <!-- Cabin Info -->
    <div class="info-section">
      <h4 class="section-title">舱房信息</h4>
      <div
        v-for="(cabin, index) in bookingData.cabins"
        :key="index"
        class="cabin-item"
      >
        <div class="cabin-header">
          <span class="cabin-name">{{ cabin.cabin_type?.name_cn }}</span>
          <span class="cabin-number">舱号: {{ cabin.cabin_number }}</span>
        </div>
        <div class="passenger-summary">
          <span>{{ cabin.adultCount }} 成人</span>
          <span v-if="cabin.childCount">, {{ cabin.childCount }} 儿童</span>
          <span v-if="cabin.infantCount">, {{ cabin.infantCount }} 婴儿</span>
        </div>
      </div>
    </div>

    <!-- Passenger Info -->
    <div class="info-section">
      <h4 class="section-title">乘客信息</h4>
      <div class="passengers-list">
        <div
          v-for="(passenger, index) in bookingData.passengers"
          :key="index"
          class="passenger-item"
        >
          <span class="passenger-index">{{ index + 1 }}</span>
          <span class="passenger-name">{{ passenger.surname }} {{ passenger.givenName }} / {{ passenger.name }}</span>
          <span class="passenger-type">{{ getPassengerTypeLabel(passenger.type) }}</span>
        </div>
      </div>
    </div>

    <!-- Contact Info -->
    <div class="info-section">
      <h4 class="section-title">联系人信息</h4>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">姓名</span>
          <span class="value">{{ bookingData.contact?.name }}</span>
        </div>
        <div class="info-item">
          <span class="label">手机</span>
          <span class="value">{{ bookingData.contact?.phone }}</span>
        </div>
        <div class="info-item">
          <span class="label">邮箱</span>
          <span class="value">{{ bookingData.contact?.email }}</span>
        </div>
      </div>
    </div>

    <!-- Price Summary -->
    <div class="price-section">
      <h4 class="section-title">费用明细</h4>
      <div class="price-items">
        <div class="price-item">
          <span>舱房费用</span>
          <span>¥{{ formatPrice(cabinTotal) }}</span>
        </div>
        <div class="price-item">
          <span>港口费</span>
          <span>¥{{ formatPrice(portFeeTotal) }}</span>
        </div>
        <div class="price-item">
          <span>服务费</span>
          <span>¥{{ formatPrice(serviceFeeTotal) }}</span>
        </div>
      </div>
      <div class="price-total">
        <span>订单总额</span>
        <span class="total-amount">¥{{ formatPrice(orderTotal) }}</span>
      </div>
    </div>

    <!-- Terms -->
    <div class="terms-section">
      <label class="checkbox-label">
        <input type="checkbox" v-model="agreed" />
        <span>我已阅读并同意 <a href="#" class="text-blue-500">预订条款</a> 和 <a href="#" class="text-blue-500">取消政策</a></span>
      </label>
    </div>

    <div class="actions flex justify-between">
      <button class="btn-prev" @click="$emit('prev')">
        上一步
      </button>
      <button
        class="btn-submit"
        :disabled="!agreed || submitting"
        @click="onSubmit"
      >
        <span v-if="submitting" class="loading-spinner"></span>
        <span>{{ submitting ? '提交中...' : '确认并支付' }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  bookingData: {
    voyage: any
    cabins: any[]
    passengers: any[]
    contact: any
  }
}>()

const emit = defineEmits<{
  submit: []
  prev: []
}>()

const agreed = ref(false)
const submitting = ref(false)

const getPassengerTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    adult: '成人',
    child: '儿童',
    infant: '婴儿'
  }
  return labels[type] || type
}

const cabinTotal = computed(() => {
  return props.bookingData.cabins.reduce((total, cabin) => {
    const price = cabin.price
    if (!price) return total
    
    const adultTotal = price.adult_price * cabin.adultCount
    const childTotal = (price.child_price || 0) * cabin.childCount
    const infantTotal = (price.infant_price || 0) * cabin.infantCount
    
    return total + adultTotal + childTotal + infantTotal
  }, 0)
})

const portFeeTotal = computed(() => {
  return props.bookingData.cabins.reduce((total, cabin) => {
    const price = cabin.price
    if (!price) return total
    return total + price.port_fee * (cabin.adultCount + cabin.childCount)
  }, 0)
})

const serviceFeeTotal = computed(() => {
  return props.bookingData.cabins.reduce((total, cabin) => {
    const price = cabin.price
    if (!price) return total
    return total + price.service_fee * (cabin.adultCount + cabin.childCount)
  }, 0)
})

const orderTotal = computed(() => {
  return cabinTotal.value + portFeeTotal.value + serviceFeeTotal.value
})

const formatPrice = (price: number) => {
  return price?.toLocaleString('zh-CN') || '0'
}

const onSubmit = () => {
  if (!agreed.value) return
  submitting.value = true
  emit('submit')
}
</script>

<style scoped>
.order-confirm {
  @apply w-full;
}

.info-section {
  @apply mb-6 pb-6 border-b border-gray-200 last:border-0;
}

.section-title {
  @apply font-medium text-lg mb-4;
}

.info-grid {
  @apply grid grid-cols-2 md:grid-cols-4 gap-4;
}

.info-item {
  @apply flex flex-col;
}

.info-item .label {
  @apply text-sm text-gray-500;
}

.info-item .value {
  @apply font-medium text-gray-800;
}

.cabin-item {
  @apply p-4 bg-gray-50 rounded-lg mb-3;
}

.cabin-header {
  @apply flex justify-between mb-2;
}

.cabin-name {
  @apply font-medium;
}

.cabin-number {
  @apply text-sm text-gray-500;
}

.passenger-summary {
  @apply text-sm text-gray-600;
}

.passengers-list {
  @apply space-y-2;
}

.passenger-item {
  @apply flex items-center gap-3 p-3 bg-gray-50 rounded-lg;
}

.passenger-index {
  @apply w-6 h-6 flex items-center justify-center bg-blue-500 text-white text-sm rounded-full;
}

.passenger-name {
  @apply flex-1 font-medium;
}

.passenger-type {
  @apply text-sm text-gray-500;
}

.price-section {
  @apply mt-6 p-4 bg-orange-50 rounded-lg;
}

.price-items {
  @apply space-y-2 mb-4;
}

.price-item {
  @apply flex justify-between text-gray-600;
}

.price-total {
  @apply flex justify-between items-center pt-4 border-t border-orange-200;
}

.total-amount {
  @apply text-2xl font-bold text-orange-600;
}

.terms-section {
  @apply mt-6;
}

.checkbox-label {
  @apply flex items-center gap-2 text-sm text-gray-600;
}

.checkbox-label input {
  @apply w-4 h-4 text-blue-500 rounded;
}

.actions {
  @apply mt-8;
}

.btn-prev {
  @apply px-6 py-2 border border-gray-300 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors;
}

.btn-submit {
  @apply flex items-center gap-2 px-8 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors;
}

.loading-spinner {
  @apply w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin;
}
</style>
