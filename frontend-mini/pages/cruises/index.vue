<template>
  <view class="cruise-list-page">
    <!-- Search Header -->
    <view class="search-header">
      <view class="search-box">
        <text class="iconfont icon-search"></text>
        <input 
          v-model="keyword"
          type="text" 
          placeholder="搜索邮轮名称"
          @confirm="searchCruises"
        />
      </view>
    </view>

    <!-- Filter Tabs -->
    <view class="filter-tabs">
      <scroll-view scroll-x class="filter-scroll">
        <view 
          v-for="filter in filters" 
          :key="filter.value"
          class="filter-item"
          :class="{ active: currentFilter === filter.value }"
          @click="setFilter(filter.value)"
        >
          {{ filter.label }}
        </view>
      </scroll-view>
    </view>

    <!-- Cruise List -->
    <scroll-view 
      scroll-y 
      class="cruise-list"
      @scrolltolower="loadMore"
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
    >
      <view v-if="loading && cruises.length === 0" class="loading-state">
        <text class="loading-text">加载中...</text>
      </view>

      <view v-else-if="cruises.length === 0" class="empty-state">
        <image class="empty-icon" src="/static/images/empty.png" mode="aspectFit"/>
        <text class="empty-text">暂无邮轮信息</text>
      </view>

      <view v-else class="cruise-grid">
        <CruiseCard
          v-for="cruise in cruises"
          :key="cruise.id"
          :cruise="cruise"
          @click="goToDetail(cruise.id)"
        />
      </view>

      <!-- Load More -->
      <view v-if="cruises.length > 0" class="load-more">
        <text v-if="loadingMore">加载更多...</text>
        <text v-else-if="noMore">没有更多了</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import CruiseCard from '@/components/CruiseCard.vue'

const keyword = ref('')
const currentFilter = ref('all')
const cruises = ref([])
const loading = ref(false)
const loadingMore = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 10

const filters = [
  { label: '全部', value: 'all' },
  { label: '热售中', value: 'active' },
  { label: '大型邮轮', value: 'large' },
  { label: '中型邮轮', value: 'medium' },
  { label: '小型邮轮', value: 'small' }
]

onMounted(() => {
  loadCruises()
})

const loadCruises = async (isRefresh = false) => {
  if (loading.value) return
  
  loading.value = true
  
  try {
    // TODO: Replace with actual API call
    // const response = await uni.request({
    //   url: 'https://api.example.com/cruises',
    //   data: {
    //     page: isRefresh ? 1 : page.value,
    //     pageSize,
    //     keyword: keyword.value,
    //     filter: currentFilter.value
    //   }
    // })
    
    // Mock data
    const mockData = [
      {
        id: '1',
        nameCn: '海洋光谱号',
        nameEn: 'Spectrum of the Seas',
        grossTonnage: 169379,
        passengerCapacity: 5622,
        crewCount: 1551,
        deckCount: 16,
        coverImages: ['/static/images/cruise-1.jpg'],
        status: 'active'
      },
      {
        id: '2',
        nameCn: '荣耀号',
        nameEn: 'MSC Bellissima',
        grossTonnage: 171598,
        passengerCapacity: 5655,
        crewCount: 1504,
        deckCount: 19,
        coverImages: ['/static/images/cruise-2.jpg'],
        status: 'active'
      }
    ]
    
    if (isRefresh) {
      cruises.value = mockData
      page.value = 1
    } else {
      cruises.value = [...cruises.value, ...mockData]
    }
    
    noMore.value = mockData.length < pageSize
  } catch (error) {
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
    refreshing.value = false
    loadingMore.value = false
  }
}

const searchCruises = () => {
  page.value = 1
  cruises.value = []
  loadCruises(true)
}

const setFilter = (value: string) => {
  currentFilter.value = value
  page.value = 1
  cruises.value = []
  loadCruises(true)
}

const loadMore = () => {
  if (noMore.value || loadingMore.value) return
  loadingMore.value = true
  page.value++
  loadCruises()
}

const onRefresh = () => {
  refreshing.value = true
  loadCruises(true)
}

const goToDetail = (id: string) => {
  uni.navigateTo({
    url: `/pages/cruises/detail?id=${id}`
  })
}
</script>

<style scoped>
.cruise-list-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.search-header {
  background: #fff;
  padding: 20rpx 30rpx;
}

.search-box {
  display: flex;
  align-items: center;
  background: #f3f4f6;
  border-radius: 32rpx;
  padding: 16rpx 24rpx;
}

.search-box input {
  flex: 1;
  margin-left: 16rpx;
  font-size: 28rpx;
}

.filter-tabs {
  background: #fff;
  border-bottom: 1rpx solid #e5e7eb;
}

.filter-scroll {
  white-space: nowrap;
  padding: 20rpx 30rpx;
}

.filter-item {
  display: inline-block;
  padding: 12rpx 24rpx;
  margin-right: 16rpx;
  font-size: 28rpx;
  color: #6b7280;
  border-radius: 28rpx;
  background: #f3f4f6;
  transition: all 0.3s;
}

.filter-item.active {
  color: #fff;
  background: #3b82f6;
}

.cruise-list {
  height: calc(100vh - 200rpx);
}

.cruise-grid {
  padding: 20rpx;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 0;
}

.empty-icon {
  width: 200rpx;
  height: 200rpx;
  margin-bottom: 20rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #9ca3af;
}

.load-more {
  text-align: center;
  padding: 30rpx;
  font-size: 26rpx;
  color: #9ca3af;
}
</style>
