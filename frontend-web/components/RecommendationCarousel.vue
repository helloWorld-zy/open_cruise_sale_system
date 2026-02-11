<template>
  <div class="recommendation-carousel">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-xl font-bold text-gray-900">{{ title }}</h2>
      <NuxtLink 
        v-if="showMoreLink" 
        :to="moreLink" 
        class="text-sm text-blue-600 hover:text-blue-700 flex items-center gap-1"
      >
        查看更多
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
        </svg>
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="w-8 h-8 border-2 border-gray-200 border-t-blue-500 rounded-full animate-spin mr-3"></div>
      <span class="text-gray-500">加载推荐...</span>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-8">
      <p class="text-gray-500">加载推荐失败</p>
      <button @click="fetchRecommendations" class="mt-2 text-blue-600 hover:text-blue-700">
        重试
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="recommendations.length === 0" class="text-center py-8 text-gray-500">
      暂无推荐
    </div>

    <!-- Carousel -->
    <div v-else class="relative">
      <!-- Scroll Container -->
      <div 
        ref="scrollContainer"
        class="flex gap-4 overflow-x-auto scrollbar-hide scroll-smooth pb-2"
        @scroll="handleScroll"
      >
        <div 
          v-for="rec in recommendations" 
          :key="rec.voyage_id"
          class="flex-shrink-0 w-72 bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow cursor-pointer"
          @click="navigateToVoyage(rec.voyage_id)"
        >
          <!-- Image -->
          <div class="relative h-40 bg-gray-100">
            <img 
              v-if="rec.cruise_image" 
              :src="rec.cruise_image" 
              :alt="rec.cruise_name"
              class="w-full h-full object-cover"
            >
            <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
              <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
            </div>

            <!-- Reason Badge -->
            <div 
              v-if="rec.reason_type"
              class="absolute top-2 left-2 px-2 py-1 text-xs font-medium rounded-full"
              :class="getReasonBadgeClass(rec.reason_type)"
            >
              {{ rec.reason }}
            </div>

            <!-- Price Change Badge -->
            <div 
              v-if="rec.price_change && rec.price_change < 0"
              class="absolute top-2 right-2 px-2 py-1 bg-red-500 text-white text-xs font-medium rounded-full"
            >
              ↓ {{ Math.abs(rec.price_change_pct).toFixed(0) }}%
            </div>
          </div>

          <!-- Content -->
          <div class="p-4">
            <h3 class="font-semibold text-gray-900 mb-1 line-clamp-1">{{ rec.cruise_name }}</h3>
            <p class="text-sm text-gray-500 mb-2">{{ rec.route_name }}</p>

            <!-- Dates -->
            <div class="flex items-center gap-2 text-sm text-gray-600 mb-3">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
              <span>{{ formatDate(rec.departure_date) }}</span>
              <span class="text-gray-400">({{ rec.duration_days }}天)</span>
            </div>

            <!-- Price -->
            <div class="flex items-baseline justify-between">
              <div>
                <span class="text-lg font-bold text-red-600">
                  ¥{{ formatPrice(rec.min_price) }}
                </span>
                <span class="text-xs text-gray-400 ml-1">起/人</span>
              </div>
              <div v-if="rec.available_cabins > 0" class="text-xs text-gray-500">
                剩{{ rec.available_cabins }}舱
              </div>
            </div>

            <!-- Match Factors -->
            <div v-if="rec.match_factors && rec.match_factors.length > 0" class="mt-3 flex flex-wrap gap-1">
              <span 
                v-for="factor in rec.match_factors.slice(0, 2)" 
                :key="factor.factor"
                class="px-2 py-0.5 bg-blue-50 text-blue-600 text-xs rounded"
              >
                {{ translateFactor(factor.factor) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Navigation Arrows -->
      <button 
        v-if="canScrollLeft"
        @click="scroll('left')"
        class="absolute left-0 top-1/2 -translate-y-1/2 -translate-x-2 w-8 h-8 bg-white rounded-full shadow-lg border border-gray-200 flex items-center justify-center hover:bg-gray-50 z-10"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
        </svg>
      </button>

      <button 
        v-if="canScrollRight"
        @click="scroll('right')"
        class="absolute right-0 top-1/2 -translate-y-1/2 translate-x-2 w-8 h-8 bg-white rounded-full shadow-lg border border-gray-200 flex items-center justify-center hover:bg-gray-50 z-10"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface MatchFactor {
  factor: string
  score: number
}

interface Recommendation {
  voyage_id: string
  cruise_id: string
  cruise_name: string
  cruise_image?: string
  voyage_number: string
  departure_date: string
  arrival_date: string
  route_name: string
  min_price: number
  max_price: number
  currency: string
  duration_days: number
  reason: string
  reason_type: string
  score: number
  match_factors?: MatchFactor[]
  tags?: string[]
  available_cabins: number
  discount_percent?: number
  price_change?: number
  price_change_pct?: number
}

const props = defineProps<{
  title?: string
  type?: 'popular' | 'personalized' | 'similar' | 'last_minute'
  voyageId?: string
  limit?: number
  showMoreLink?: boolean
  moreLink?: string
}>()

const title = computed(() => props.title || getDefaultTitle())
const recommendations = ref<Recommendation[]>([])
const loading = ref(false)
const error = ref(false)
const scrollContainer = ref<HTMLElement | null>(null)
const canScrollLeft = ref(false)
const canScrollRight = ref(true)

function getDefaultTitle(): string {
  switch (props.type) {
    case 'popular':
      return '热门推荐'
    case 'personalized':
      return '为您推荐'
    case 'similar':
      return '相似航次'
    case 'last_minute':
      return '临期特惠'
    default:
      return '推荐航次'
  }
}

function getReasonBadgeClass(reasonType: string): string {
  switch (reasonType) {
    case 'popular':
      return 'bg-orange-100 text-orange-700'
    case 'similar':
      return 'bg-blue-100 text-blue-700'
    case 'last_minute':
      return 'bg-red-100 text-red-700'
    case 'price_drop':
      return 'bg-green-100 text-green-700'
    default:
      return 'bg-gray-100 text-gray-700'
  }
}

function translateFactor(factor: string): string {
  const translations: Record<string, string> = {
    preferred_cruise: '关注邮轮',
    preferred_route: '偏好航线',
    preferred_cabin_type: '喜欢房型',
    price_match: '价格合适',
    preferred_month: '偏好月份',
    popularity: '热门',
    last_minute: '临期'
  }
  return translations[factor] || factor
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日`
}

function formatPrice(price: number): string {
  if (price >= 10000) {
    return (price / 10000).toFixed(1) + '万'
  }
  return price.toLocaleString()
}

function navigateToVoyage(voyageId: string) {
  navigateTo(`/voyages/${voyageId}`)
}

function handleScroll() {
  if (!scrollContainer.value) return
  
  const container = scrollContainer.value
  canScrollLeft.value = container.scrollLeft > 0
  canScrollRight.value = container.scrollLeft < container.scrollWidth - container.clientWidth - 10
}

function scroll(direction: 'left' | 'right') {
  if (!scrollContainer.value) return
  
  const scrollAmount = 300
  scrollContainer.value.scrollBy({
    left: direction === 'left' ? -scrollAmount : scrollAmount,
    behavior: 'smooth'
  })
}

async function fetchRecommendations() {
  loading.value = true
  error.value = false

  try {
    const { $api } = useNuxtApp()
    let endpoint = '/analytics/recommendations'
    
    switch (props.type) {
      case 'popular':
        endpoint = '/analytics/recommendations/popular'
        break
      case 'personalized':
        endpoint = '/analytics/recommendations/personalized'
        break
      case 'similar':
        if (props.voyageId) {
          endpoint = `/analytics/recommendations/similar/${props.voyageId}`
        }
        break
      case 'last_minute':
        endpoint = '/analytics/recommendations/last-minute'
        break
    }

    const response = await $api.get(endpoint, {
      params: {
        limit: props.limit || 10
      }
    })

    recommendations.value = response.data.data || []
  } catch (err) {
    console.error('Failed to fetch recommendations:', err)
    error.value = true
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchRecommendations()
  handleScroll()
})
</script>

<style scoped>
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
