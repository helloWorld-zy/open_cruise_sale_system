<template>
  <div class="select-voyage">
    <h3 class="text-xl font-semibold mb-6">选择出发日期</h3>
    
    <div v-if="loading" class="loading-state">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
      <p class="text-gray-500 mt-4 text-center">加载航次中...</p>
    </div>

    <div v-else-if="voyages.length === 0" class="empty-state text-center py-12">
      <p class="text-gray-500">暂无可用航次</p>
    </div>

    <div v-else class="voyages-grid">
      <div
        v-for="voyage in voyages"
        :key="voyage.id"
        class="voyage-card"
        :class="{ 'selected': selected?.id === voyage.id }"
        @click="selectVoyage(voyage)"
      >
        <div class="voyage-date">
          <div class="date-day">{{ formatDay(voyage.departure_date) }}</div>
          <div class="date-month">{{ formatMonth(voyage.departure_date) }}</div>
        </div>
        
        <div class="voyage-info">
          <div class="info-row">
            <span class="label">出发:</span>
            <span class="value">{{ voyage.departure_date }} {{ voyage.departure_time }}</span>
          </div>
          <div class="info-row">
            <span class="label">返回:</span>
            <span class="value">{{ voyage.arrival_date }} {{ voyage.arrival_time }}</span>
          </div>
          <div class="info-row">
            <span class="label">航线:</span>
            <span class="value">{{ voyage.route?.name }}</span>
          </div>
          <div class="info-row">
            <span class="label">航程:</span>
            <span class="value">{{ voyage.route?.duration_days }}天</span>
          </div>
        </div>

        <div class="voyage-price" v-if="voyage.min_price">
          <span class="price-label">起</span>
          <span class="price-value">¥{{ formatPrice(voyage.min_price) }}</span>
        </div>

        <div class="selection-indicator" v-if="selected?.id === voyage.id">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
      </div>
    </div>

    <div class="actions mt-8 flex justify-between">
      <button
        class="btn-prev"
        @click="$emit('cancel')"
      >
        取消
      </button>
      <button
        class="btn-next"
        :disabled="!selected"
        @click="onNext"
      >
        下一步
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const props = defineProps<{
  cruiseId: string
  initialSelection?: any
}>()

const emit = defineEmits<{
  select: [voyage: any]
  next: []
  cancel: []
}>()

const voyages = ref<any[]>([])
const selected = ref<any>(props.initialSelection)
const loading = ref(false)

onMounted(async () => {
  await loadVoyages()
})

const loadVoyages = async () => {
  loading.value = true
  try {
    const response = await $fetch(`/api/voyages`, {
      params: {
        cruise_id: props.cruiseId,
        booking_status: 'open'
      }
    })
    voyages.value = response.data || []
  } catch (err) {
    console.error('Failed to load voyages:', err)
  } finally {
    loading.value = false
  }
}

const selectVoyage = (voyage: any) => {
  selected.value = voyage
  emit('select', voyage)
}

const onNext = () => {
  if (selected.value) {
    emit('next')
  }
}

const formatDay = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.getDate()
}

const formatMonth = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月`
}

const formatPrice = (price: number) => {
  return price.toLocaleString('zh-CN')
}
</script>

<style scoped>
.select-voyage {
  @apply w-full;
}

.voyages-grid {
  @apply grid gap-4;
}

.voyage-card {
  @apply flex items-center p-4 border-2 border-gray-200 rounded-lg cursor-pointer transition-all duration-200 hover:border-blue-300;
}

.voyage-card.selected {
  @apply border-blue-500 bg-blue-50;
}

.voyage-date {
  @apply flex flex-col items-center justify-center w-16 h-16 bg-blue-100 rounded-lg mr-4;
}

.date-day {
  @apply text-2xl font-bold text-blue-600;
}

.date-month {
  @apply text-sm text-blue-500;
}

.voyage-info {
  @apply flex-1;
}

.info-row {
  @apply flex text-sm mb-1;
}

.info-row .label {
  @apply text-gray-500 w-12;
}

.info-row .value {
  @apply text-gray-800;
}

.voyage-price {
  @apply flex flex-col items-end mr-4;
}

.price-label {
  @apply text-xs text-gray-500;
}

.price-value {
  @apply text-xl font-bold text-orange-500;
}

.selection-indicator {
  @apply w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center;
}

.btn-prev {
  @apply px-6 py-2 border border-gray-300 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors;
}

.btn-next {
  @apply px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors;
}

.loading-state {
  @apply py-12;
}
</style>
