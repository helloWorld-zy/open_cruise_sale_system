<template>
  <div class="select-cabin">
    <h3 class="text-xl font-semibold mb-2">选择舱房</h3>
    <p class="text-gray-500 mb-6">{{ voyage.voyage_number }} | {{ voyage.departure_date }} 出发</p>
    
    <div v-if="loading" class="loading-state">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
      <p class="text-gray-500 mt-4 text-center">加载舱房信息...</p>
    </div>

    <div v-else-if="cabins.length === 0" class="empty-state text-center py-12">
      <p class="text-gray-500">暂无可用舱房</p>
    </div>

    <div v-else class="cabins-container">
      <div
        v-for="cabin in cabins"
        :key="cabin.id"
        class="cabin-section"
      >
        <div class="cabin-header">
          <div class="cabin-type-info">
            <h4 class="font-semibold text-lg">{{ cabin.cabin_type?.name_cn }}</h4>
            <p class="text-sm text-gray-500">{{ cabin.cabin_type?.name_en }}</p>
            <div class="cabin-features mt-2">
              <span
                v-for="feature in cabin.cabin_type?.features"
                :key="feature"
                class="feature-tag"
              >
                {{ feature }}
              </span>
            </div>
          </div>
          
          <div class="cabin-price" v-if="cabin.price">
            <span class="price-value">¥{{ formatPrice(cabin.price.adult_price) }}</span>
            <span class="price-unit">/人起</span>
          </div>
        </div>

        <div class="cabin-selection">
          <div class="passenger-count">
            <label>成人 (12岁以上)</label>
            <div class="counter">
              <button @click="decrement(cabin.id, 'adult')" :disabled="getCount(cabin.id, 'adult') <= 1">-</button>
              <span>{{ getCount(cabin.id, 'adult') }}</span>
              <button @click="increment(cabin.id, 'adult')">+</button>
            </div>
          </div>
          
          <div class="passenger-count">
            <label>儿童 (2-11岁)</label>
            <div class="counter">
              <button @click="decrement(cabin.id, 'child')" :disabled="getCount(cabin.id, 'child') <= 0">-</button>
              <span>{{ getCount(cabin.id, 'child') }}</span>
              <button @click="increment(cabin.id, 'child')">+</button>
            </div>
          </div>
          
          <div class="passenger-count">
            <label>婴儿 (2岁以下)</label>
            <div class="counter">
              <button @click="decrement(cabin.id, 'infant')" :disabled="getCount(cabin.id, 'infant') <= 0">-</button>
              <span>{{ getCount(cabin.id, 'infant') }}</span>
              <button @click="increment(cabin.id, 'infant')">+</button>
            </div>
          </div>

          <div class="available-count">
            剩余: {{ cabin.inventory?.available_cabins || 0 }} 间
          </div>

          <button
            class="btn-add"
            :class="{ 'added': isSelected(cabin.id) }"
            :disabled="cabin.inventory?.available_cabins === 0"
            @click="toggleCabin(cabin)"
          >
            {{ isSelected(cabin.id) ? '已选择' : '选择此舱房' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Selected Summary -->
    <div v-if="selectedCabins.length > 0" class="selected-summary">
      <h4 class="font-semibold mb-3">已选择舱房</h4>
      <div
        v-for="(item, index) in selectedCabins"
        :key="index"
        class="selected-item"
      >
        <span>{{ item.cabin_type?.name_cn }} - {{ item.adultCount }}成人</span>
        <span v-if="item.childCount">, {{ item.childCount }}儿童</span>
        <span v-if="item.infantCount">, {{ item.infantCount }}婴儿</span>
        <button class="btn-remove" @click="removeCabin(index)">×</button>
      </div>
      <div class="total-price">
        <span>预估总价:</span>
        <span class="price">¥{{ formatPrice(estimatedTotal) }}</span>
      </div>
    </div>

    <div class="actions mt-8 flex justify-between">
      <button class="btn-prev" @click="$emit('prev')">
        上一步
      </button>
      <button
        class="btn-next"
        :disabled="selectedCabins.length === 0"
        @click="onNext"
      >
        下一步
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const props = defineProps<{
  voyage: any
  initialSelection?: any[]
}>()

const emit = defineEmits<{
  select: [cabins: any[]]
  prev: []
  next: []
}>()

const cabins = ref<any[]>([])
const selectedCabins = ref<any[]>(props.initialSelection || [])
const passengerCounts = ref<Record<string, { adult: number; child: number; infant: number }>>({})
const loading = ref(false)

onMounted(async () => {
  await loadCabins()
  // Initialize passenger counts
  cabins.value.forEach(cabin => {
    if (!passengerCounts.value[cabin.id]) {
      passengerCounts.value[cabin.id] = { adult: 2, child: 0, infant: 0 }
    }
  })
})

const loadCabins = async () => {
  loading.value = true
  try {
    const [cabinsResponse, pricesResponse, inventoryResponse] = await Promise.all([
      $fetch(`/api/cabins`, {
        params: { voyage_id: props.voyage.id, status: 'available' }
      }),
      $fetch(`/api/prices`, {
        params: { voyage_id: props.voyage.id }
      }),
      $fetch(`/api/inventory`, {
        params: { voyage_id: props.voyage.id }
      })
    ])

    const cabinData = cabinsResponse.data || []
    const prices = pricesResponse.data || []
    const inventory = inventoryResponse.data || []

    // Merge data
    cabins.value = cabinData.map((cabin: any) => ({
      ...cabin,
      price: prices.find((p: any) => p.cabin_type_id === cabin.cabin_type_id),
      inventory: inventory.find((i: any) => i.cabin_type_id === cabin.cabin_type_id)
    }))
  } catch (err) {
    console.error('Failed to load cabins:', err)
  } finally {
    loading.value = false
  }
}

const getCount = (cabinId: string, type: 'adult' | 'child' | 'infant') => {
  return passengerCounts.value[cabinId]?.[type] || 0
}

const increment = (cabinId: string, type: 'adult' | 'child' | 'infant') => {
  if (!passengerCounts.value[cabinId]) {
    passengerCounts.value[cabinId] = { adult: 2, child: 0, infant: 0 }
  }
  passengerCounts.value[cabinId][type]++
}

const decrement = (cabinId: string, type: 'adult' | 'child' | 'infant') => {
  if (passengerCounts.value[cabinId] && passengerCounts.value[cabinId][type] > 0) {
    passengerCounts.value[cabinId][type]--
  }
}

const isSelected = (cabinId: string) => {
  return selectedCabins.value.some(c => c.id === cabinId)
}

const toggleCabin = (cabin: any) => {
  const index = selectedCabins.value.findIndex(c => c.id === cabin.id)
  if (index >= 0) {
    selectedCabins.value.splice(index, 1)
  } else {
    const counts = passengerCounts.value[cabin.id] || { adult: 2, child: 0, infant: 0 }
    selectedCabins.value.push({
      ...cabin,
      adultCount: counts.adult,
      childCount: counts.child,
      infantCount: counts.infant
    })
  }
}

const removeCabin = (index: number) => {
  selectedCabins.value.splice(index, 1)
}

const estimatedTotal = computed(() => {
  return selectedCabins.value.reduce((total, cabin) => {
    const price = cabin.price
    if (!price) return total
    
    const adultTotal = price.adult_price * cabin.adultCount
    const childTotal = (price.child_price || 0) * cabin.childCount
    const infantTotal = (price.infant_price || 0) * cabin.infantCount
    const fees = (price.port_fee + price.service_fee) * (cabin.adultCount + cabin.childCount)
    
    return total + adultTotal + childTotal + infantTotal + fees
  }, 0)
})

const onNext = () => {
  emit('select', selectedCabins.value)
  emit('next')
}

const formatPrice = (price: number) => {
  return price?.toLocaleString('zh-CN') || '0'
}
</script>

<style scoped>
.select-cabin {
  @apply w-full;
}

.cabins-container {
  @apply space-y-6;
}

.cabin-section {
  @apply border-2 border-gray-200 rounded-lg p-4;
}

.cabin-header {
  @apply flex justify-between items-start mb-4 pb-4 border-b;
}

.cabin-type-info {
  @apply flex-1;
}

.cabin-features {
  @apply flex flex-wrap gap-2;
}

.feature-tag {
  @apply px-2 py-1 bg-gray-100 text-xs text-gray-600 rounded;
}

.cabin-price {
  @apply text-right;
}

.price-value {
  @apply text-2xl font-bold text-orange-500;
}

.price-unit {
  @apply text-sm text-gray-500;
}

.cabin-selection {
  @apply grid grid-cols-4 gap-4 items-center;
}

.passenger-count {
  @apply flex flex-col;
}

.passenger-count label {
  @apply text-sm text-gray-500 mb-1;
}

.counter {
  @apply flex items-center border rounded-lg;
}

.counter button {
  @apply w-8 h-8 flex items-center justify-center text-gray-600 hover:bg-gray-100 disabled:opacity-30;
}

.counter span {
  @apply w-10 text-center font-medium;
}

.available-count {
  @apply text-sm text-gray-500;
}

.btn-add {
  @apply px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors;
}

.btn-add.added {
  @apply bg-green-500 hover:bg-green-600;
}

.btn-add:disabled {
  @apply bg-gray-300 cursor-not-allowed;
}

.selected-summary {
  @apply mt-6 p-4 bg-gray-50 rounded-lg;
}

.selected-item {
  @apply flex justify-between items-center py-2 border-b border-gray-200 last:border-0;
}

.btn-remove {
  @apply w-6 h-6 flex items-center justify-center text-red-500 hover:bg-red-50 rounded;
}

.total-price {
  @apply flex justify-between items-center mt-4 pt-4 border-t border-gray-300;
}

.total-price .price {
  @apply text-2xl font-bold text-orange-500;
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
