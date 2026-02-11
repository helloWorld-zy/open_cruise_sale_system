<template>
  <view class="cruise-detail-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-container">
      <text class="loading-text">加载中...</text>
    </view>

    <!-- Content -->
    <view v-else-if="cruise" class="content-container">
      <!-- Image Gallery -->
      <swiper class="image-gallery" :indicator-dots="true" :autoplay="true" :interval="5000">
        <swiper-item v-for="(image, index) in cruise.coverImages" :key="index">
          <image class="gallery-image" :src="image" mode="aspectFill" />
        </swiper-item>
        <swiper-item v-if="!cruise.coverImages || cruise.coverImages.length === 0">
          <image class="gallery-image" src="/static/images/placeholder.jpg" mode="aspectFill" />
        </swiper-item>
      </swiper>

      <!-- Basic Info Card -->
      <view class="info-card">
        <view class="info-header">
          <text class="cruise-name">{{ cruise.nameCn }}</text>
          <text v-if="cruise.nameEn" class="cruise-name-en">{{ cruise.nameEn }}</text>
          <view class="status-badge" :class="cruise.status">
            {{ statusText }}
          </view>
        </view>

        <view class="specs-grid">
          <view class="spec-item">
            <text class="spec-label">总吨位</text>
            <text class="spec-value">{{ cruise.grossTonnage?.toLocaleString() || '-' }} 吨</text>
          </view>
          <view class="spec-item">
            <text class="spec-label">载客量</text>
            <text class="spec-value">{{ cruise.passengerCapacity?.toLocaleString() || '-' }} 人</text>
          </view>
          <view class="spec-item">
            <text class="spec-label">甲板层</text>
            <text class="spec-value">{{ cruise.deckCount || '-' }} 层</text>
          </view>
          <view class="spec-item">
            <text class="spec-label">船员数</text>
            <text class="spec-value">{{ cruise.crewCount?.toLocaleString() || '-' }} 人</text>
          </view>
        </view>
      </view>

      <!-- Tabs -->
      <view class="tabs-container">
        <view class="tab-header">
          <view 
            v-for="tab in tabs" 
            :key="tab.key"
            class="tab-item"
            :class="{ active: activeTab === tab.key }"
            @click="activeTab = tab.key"
          >
            {{ tab.label }}
          </view>
        </view>

        <view class="tab-content">
          <!-- Cabin Types -->
          <view v-if="activeTab === 'cabins'" class="cabins-tab">
            <view 
              v-for="cabin in cruise.cabinTypes" 
              :key="cabin.id"
              class="cabin-item"
            >
              <image 
                class="cabin-image" 
                :src="cabin.images?.[0] || '/static/images/placeholder.jpg'"
                mode="aspectFill"
              />
              <view class="cabin-info">
                <text class="cabin-name">{{ cabin.name }}</text>
                <text class="cabin-specs">{{ cabin.minAreaSqm }}-{{ cabin.maxAreaSqm }}㎡ | 可住{{ cabin.standardGuests }}人</text>
                <view v-if="cabin.featureTags && cabin.featureTags.length > 0" class="cabin-tags">
                  <text v-for="tag in cabin.featureTags" :key="tag" class="cabin-tag">{{ tag }}</text>
                </view>
              </view>
              <view class="cabin-price">
                <text class="price-label">起价</text>
                <text class="price-value">¥{{ cabin.price || '咨询' }}</text>
              </view>
            </view>
            <view v-if="!cruise.cabinTypes || cruise.cabinTypes.length === 0" class="empty-tab">
              <text>暂无舱房信息</text>
            </view>
          </view>

          <!-- Facilities -->
          <view v-if="activeTab === 'facilities'" class="facilities-tab">
            <view 
              v-for="facility in cruise.facilities" 
              :key="facility.id"
              class="facility-item"
            >
              <view class="facility-icon">
                <text class="iconfont icon-facility"></text>
              </view>
              <view class="facility-info">
                <view class="facility-header">
                  <text class="facility-name">{{ facility.name }}</text>
                  <text v-if="!facility.isFree" class="facility-price">¥{{ facility.price }}</text>
                </view>
                <text v-if="facility.deckNumber" class="facility-location">{{ facility.deckNumber }} 层甲板</text>
                <text v-if="facility.openTime" class="facility-time">营业时间: {{ facility.openTime }}</text>
              </view>
            </view>
            <view v-if="!cruise.facilities || cruise.facilities.length === 0" class="empty-tab">
              <text>暂无设施信息</text>
            </view>
          </view>

          <!-- Specs -->
          <view v-if="activeTab === 'specs'" class="specs-tab">
            <view class="specs-list">
              <view class="spec-row">
                <text class="spec-label">邮轮名称</text>
                <text class="spec-value">{{ cruise.nameCn }}</text>
              </view>
              <view class="spec-row">
                <text class="spec-label">英文名称</text>
                <text class="spec-value">{{ cruise.nameEn || '-' }}</text>
              </view>
              <view class="spec-row">
                <text class="spec-label">邮轮编号</text>
                <text class="spec-value">{{ cruise.code }}</text>
              </view>
              <view class="spec-row">
                <text class="spec-label">建造年份</text>
                <text class="spec-value">{{ cruise.builtYear || '-' }} 年</text>
              </view>
              <view class="spec-row">
                <text class="spec-label">翻新年份</text>
                <text class="spec-value">{{ cruise.renovatedYear || '-' }} 年</text>
              </view>
              <view class="spec-row">
                <text class="spec-label">长度</text>
                <text class="spec-value">{{ cruise.lengthMeters || '-' }} 米</text>
              </view>
              <view class="spec-row">
                <text class="spec-label">宽度</text>
                <text class="spec-value">{{ cruise.widthMeters || '-' }} 米</text>
              </view>
            </view>
          </view>
        </view>
      </view>

      <!-- Bottom Action Bar -->
      <view class="action-bar">
        <view class="action-left">
          <view class="action-btn" @click="goBack">
            <text class="iconfont icon-back"></text>
            <text class="action-text">返回</text>
          </view>
          <view class="action-btn" @click="shareCruise">
            <text class="iconfont icon-share"></text>
            <text class="action-text">分享</text>
          </view>
        </view>
        <button class="booking-btn" @click="startBooking">立即预订</button>
      </view>
    </view>

    <!-- Error/Empty -->
    <view v-else class="error-container">
      <image class="error-icon" src="/static/images/empty.png" mode="aspectFit"/>
      <text class="error-text">邮轮信息不存在</text>
      <button class="back-btn" @click="goBack">返回列表</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const cruise = ref<any>(null)
const loading = ref(true)
const activeTab = ref('cabins')

const tabs = [
  { key: 'cabins', label: '舱房' },
  { key: 'facilities', label: '设施' },
  { key: 'specs', label: '参数' }
]

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const id = currentPage.options?.id
  
  if (id) {
    loadCruiseDetail(id)
  } else {
    loading.value = false
  }
})

const loadCruiseDetail = async (id: string) => {
  try {
    loading.value = true
    
    // TODO: Replace with actual API call
    // const response = await uni.request({
    //   url: `https://api.example.com/cruises/${id}`
    // })
    
    // Mock data
    cruise.value = {
      id,
      nameCn: '海洋光谱号',
      nameEn: 'Spectrum of the Seas',
      code: 'SOTS001',
      grossTonnage: 169379,
      passengerCapacity: 5622,
      crewCount: 1551,
      deckCount: 16,
      builtYear: 2019,
      lengthMeters: 347,
      widthMeters: 41,
      status: 'active',
      coverImages: ['/static/images/cruise-1.jpg'],
      cabinTypes: [
        {
          id: 'c1',
          name: '内舱房',
          minAreaSqm: 15,
          maxAreaSqm: 17,
          standardGuests: 2,
          price: 2999,
          featureTags: ['性价比', '舒适']
        },
        {
          id: 'c2',
          name: '海景房',
          minAreaSqm: 17,
          maxAreaSqm: 20,
          standardGuests: 2,
          price: 3999,
          featureTags: ['海景', '采光好']
        }
      ],
      facilities: [
        {
          id: 'f1',
          name: '皇家剧院',
          deckNumber: 3,
          isFree: true
        },
        {
          id: 'f2',
          name: '冲浪模拟器',
          deckNumber: 14,
          isFree: true
        }
      ]
    }
  } catch (error) {
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

const statusText = computed(() => {
  const map: Record<string, string> = {
    active: '热售中',
    inactive: '已下架',
    maintenance: '维护中'
  }
  return map[cruise.value?.status] || cruise.value?.status
})

const goBack = () => {
  uni.navigateBack()
}

const shareCruise = () => {
  // TODO: Implement share functionality
  uni.showShareMenu({
    withShareTicket: true,
    menus: ['shareAppMessage', 'shareTimeline']
  })
}

const startBooking = () => {
  uni.navigateTo({
    url: `/pages/booking/index?cruiseId=${cruise.value.id}`
  })
}
</script>

<style scoped>
.cruise-detail-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 120rpx;
}

.loading-container,
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

.loading-text,
.error-text {
  font-size: 28rpx;
  color: #9ca3af;
  margin-top: 20rpx;
}

.error-icon {
  width: 200rpx;
  height: 200rpx;
}

.back-btn {
  margin-top: 40rpx;
  padding: 16rpx 40rpx;
  background: #3b82f6;
  color: #fff;
  border-radius: 8rpx;
  font-size: 28rpx;
}

.image-gallery {
  height: 500rpx;
}

.gallery-image {
  width: 100%;
  height: 100%;
}

.info-card {
  background: #fff;
  margin: 20rpx;
  padding: 30rpx;
  border-radius: 16rpx;
}

.info-header {
  margin-bottom: 24rpx;
}

.cruise-name {
  font-size: 36rpx;
  font-weight: bold;
  color: #1f2937;
}

.cruise-name-en {
  font-size: 24rpx;
  color: #6b7280;
  margin-top: 8rpx;
}

.status-badge {
  display: inline-block;
  padding: 8rpx 16rpx;
  border-radius: 8rpx;
  font-size: 22rpx;
  margin-top: 16rpx;
}

.status-badge.active {
  background: #dcfce7;
  color: #166534;
}

.specs-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20rpx;
}

.spec-item {
  text-align: center;
}

.spec-label {
  display: block;
  font-size: 22rpx;
  color: #9ca3af;
  margin-bottom: 4rpx;
}

.spec-value {
  font-size: 26rpx;
  font-weight: 500;
  color: #1f2937;
}

.tabs-container {
  background: #fff;
  margin: 20rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.tab-header {
  display: flex;
  border-bottom: 1rpx solid #e5e7eb;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  font-size: 28rpx;
  color: #6b7280;
  position: relative;
}

.tab-item.active {
  color: #3b82f6;
  font-weight: 500;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 40rpx;
  height: 4rpx;
  background: #3b82f6;
  border-radius: 2rpx;
}

.tab-content {
  padding: 20rpx;
}

.cabin-item {
  display: flex;
  padding: 20rpx;
  background: #f9fafb;
  border-radius: 12rpx;
  margin-bottom: 16rpx;
}

.cabin-image {
  width: 120rpx;
  height: 120rpx;
  border-radius: 8rpx;
  margin-right: 20rpx;
}

.cabin-info {
  flex: 1;
}

.cabin-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #1f2937;
}

.cabin-specs {
  font-size: 24rpx;
  color: #6b7280;
  margin-top: 8rpx;
}

.cabin-tags {
  display: flex;
  gap: 8rpx;
  margin-top: 12rpx;
}

.cabin-tag {
  padding: 4rpx 12rpx;
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 20rpx;
  border-radius: 4rpx;
}

.cabin-price {
  text-align: right;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.price-label {
  font-size: 20rpx;
  color: #9ca3af;
}

.price-value {
  font-size: 32rpx;
  font-weight: bold;
  color: #dc2626;
}

.facility-item {
  display: flex;
  padding: 20rpx;
  border-bottom: 1rpx solid #f3f4f6;
}

.facility-icon {
  width: 80rpx;
  height: 80rpx;
  background: #f3f4f6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
}

.facility-info {
  flex: 1;
}

.facility-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.facility-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #1f2937;
}

.facility-price {
  font-size: 24rpx;
  color: #f59e0b;
}

.facility-location,
.facility-time {
  font-size: 24rpx;
  color: #6b7280;
  margin-top: 8rpx;
}

.specs-list {
  padding: 10rpx;
}

.spec-row {
  display: flex;
  justify-content: space-between;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f3f4f6;
}

.spec-row:last-child {
  border-bottom: none;
}

.empty-tab {
  text-align: center;
  padding: 60rpx 0;
  color: #9ca3af;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  border-top: 1rpx solid #e5e7eb;
  padding: 20rpx 30rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
}

.action-left {
  display: flex;
  gap: 40rpx;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.action-text {
  font-size: 20rpx;
  color: #6b7280;
  margin-top: 4rpx;
}

.booking-btn {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: #fff;
  padding: 20rpx 60rpx;
  border-radius: 40rpx;
  font-size: 30rpx;
  font-weight: 500;
  border: none;
}
</style>
