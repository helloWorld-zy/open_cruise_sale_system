<template>
  <view class="home-container">
    <!-- Hero Banner -->
    <view class="hero-section">
      <image 
        class="hero-bg" 
        src="/static/images/cruise-hero.jpg" 
        mode="aspectFill"
      />
      <view class="hero-overlay">
        <text class="hero-title">探索豪华邮轮</text>
        <text class="hero-subtitle">开启您的海上度假之旅</text>
        <button class="hero-btn" @click="goToCruises">立即探索</button>
      </view>
    </view>

    <!-- Featured Cruises -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">热门邮轮</text>
        <text class="section-more" @click="goToCruises">查看更多 ></text>
      </view>
      
      <scroll-view scroll-x class="featured-scroll">
        <view 
          v-for="cruise in featuredCruises" 
          :key="cruise.id"
          class="featured-item"
          @click="goToDetail(cruise.id)"
        >
          <image 
            class="featured-image" 
            :src="cruise.coverImage || '/static/images/placeholder.jpg'" 
            mode="aspectFill"
          />
          <view class="featured-info">
            <text class="featured-name">{{ cruise.nameCn }}</text>
            <text class="featured-capacity">载客 {{ cruise.passengerCapacity }} 人</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- Quick Actions -->
    <view class="section">
      <view class="quick-actions">
        <view class="action-item" @click="goToCruises">
          <view class="action-icon cruise-icon">
            <text class="iconfont icon-cruise"></text>
          </view>
          <text class="action-text">邮轮列表</text>
        </view>
        <view class="action-item" @click="goToOrders">
          <view class="action-icon order-icon">
            <text class="iconfont icon-order"></text>
          </view>
          <text class="action-text">我的订单</text>
        </view>
        <view class="action-item" @click="goToProfile">
          <view class="action-icon profile-icon">
            <text class="iconfont icon-user"></text>
          </view>
          <text class="action-text">个人中心</text>
        </view>
        <view class="action-item" @click="goToService">
          <view class="action-icon service-icon">
            <text class="iconfont icon-service"></text>
          </view>
          <text class="action-text">客服咨询</text>
        </view>
      </view>
    </view>

    <!-- Service Introduction -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">为什么选择我们</text>
      </view>
      <view class="service-grid">
        <view class="service-item">
          <view class="service-icon-wrapper blue">
            <text class="iconfont icon-shield"></text>
          </view>
          <text class="service-title">安全可靠</text>
          <text class="service-desc">严选正规邮轮公司</text>
        </view>
        <view class="service-item">
          <view class="service-icon-wrapper green">
            <text class="iconfont tag"></text>
          </view>
          <text class="service-title">价格透明</text>
          <text class="service-desc">无隐藏费用</text>
        </view>
        <view class="service-item">
          <view class="service-icon-wrapper orange">
            <text class="iconfont icon-support"></text>
          </view>
          <text class="service-title">专属客服</text>
          <text class="service-desc">7x24小时服务</text>
        </view>
        <view class="service-item">
          <view class="service-icon-wrapper purple">
            <text class="iconfont icon-gift"></text>
          </view>
          <text class="service-title">优惠多多</text>
          <text class="service-desc">会员专属折扣</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const featuredCruises = ref([])

onMounted(() => {
  // Load featured cruises from API
  loadFeaturedCruises()
})

const loadFeaturedCruises = async () => {
  // TODO: Replace with actual API call
  featuredCruises.value = [
    {
      id: '1',
      nameCn: '海洋光谱号',
      passengerCapacity: 5000,
      coverImage: '/static/images/cruise-1.jpg'
    },
    {
      id: '2', 
      nameCn: '荣耀号',
      passengerCapacity: 4500,
      coverImage: '/static/images/cruise-2.jpg'
    },
    {
      id: '3',
      nameCn: '爱达魔都号', 
      passengerCapacity: 4000,
      coverImage: '/static/images/cruise-3.jpg'
    }
  ]
}

const goToCruises = () => {
  uni.navigateTo({
    url: '/pages/cruises/index'
  })
}

const goToDetail = (id: string) => {
  uni.navigateTo({
    url: `/pages/cruises/detail?id=${id}`
  })
}

const goToOrders = () => {
  uni.navigateTo({
    url: '/pages/orders/index'
  })
}

const goToProfile = () => {
  uni.switchTab({
    url: '/pages/profile/index'
  })
}

const goToService = () => {
  uni.navigateTo({
    url: '/pages/service/index'
  })
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.hero-section {
  position: relative;
  height: 400rpx;
}

.hero-bg {
  width: 100%;
  height: 100%;
}

.hero-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40rpx;
}

.hero-title {
  font-size: 48rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 16rpx;
}

.hero-subtitle {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 32rpx;
}

.hero-btn {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: #fff;
  padding: 20rpx 60rpx;
  border-radius: 40rpx;
  font-size: 28rpx;
  border: none;
}

.section {
  background: #fff;
  margin-top: 20rpx;
  padding: 30rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #1f2937;
}

.section-more {
  font-size: 26rpx;
  color: #3b82f6;
}

.featured-scroll {
  white-space: nowrap;
}

.featured-item {
  display: inline-block;
  width: 280rpx;
  margin-right: 20rpx;
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
}

.featured-image {
  width: 100%;
  height: 200rpx;
}

.featured-info {
  padding: 16rpx;
}

.featured-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #1f2937;
  display: block;
  margin-bottom: 8rpx;
}

.featured-capacity {
  font-size: 24rpx;
  color: #6b7280;
}

.quick-actions {
  display: flex;
  justify-content: space-around;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.action-icon {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12rpx;
}

.cruise-icon {
  background: #dbeafe;
}

.order-icon {
  background: #dcfce7;
}

.profile-icon {
  background: #fce7f3;
}

.service-icon {
  background: #fef3c7;
}

.action-text {
  font-size: 26rpx;
  color: #4b5563;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20rpx;
}

.service-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.service-icon-wrapper {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12rpx;
}

.blue {
  background: #dbeafe;
}

.green {
  background: #dcfce7;
}

.orange {
  background: #ffedd5;
}

.purple {
  background: #f3e8ff;
}

.service-title {
  font-size: 26rpx;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 4rpx;
}

.service-desc {
  font-size: 22rpx;
  color: #9ca3af;
}
</style>
