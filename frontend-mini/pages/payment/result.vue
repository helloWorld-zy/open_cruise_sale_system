<template>
  <view class="payment-result">
    <!-- Success -->
    <view v-if="status === 'success'" class="result success">
      <image class="icon" src="/static/icons/success.png" mode="aspectFit"/>
      <text class="title">支付成功！</text>
      <text class="subtitle">您的订单已确认</text>
      
      <view class="order-info" v-if="order">
        <view class="info-row">
          <text>订单号</text>
          <text>{{ order.order_number }}</text>
        </view>
        <view class="info-row">
          <text>支付金额</text>
          <text class="price">¥{{ order.total_amount }}</text>
        </view>
      </view>

      <view class="actions">
        <button class="btn-primary" @click="goToOrder">查看订单</button>
        <button class="btn-secondary" @click="goHome">返回首页</button>
      </view>
    </view>

    <!-- Failed -->
    <view v-else-if="status === 'failed'" class="result failed">
      <image class="icon" src="/static/icons/failed.png" mode="aspectFit"/>
      <text class="title">支付失败</text>
      <text class="subtitle">请重新尝试支付</text>
      
      <view class="actions">
        <button class="btn-primary" @click="retryPayment">重新支付</button>
        <button class="btn-secondary" @click="goToOrders">查看订单</button>
      </view>
    </view>

    <!-- Loading -->
    <view v-else class="result loading">
      <view class="loading-spinner"/>
      <text class="title">正在查询支付结果...</text>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const status = ref('loading')
const order = ref(null)
const orderId = ref('')

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const { order_id, status: queryStatus } = currentPage.options
  
  orderId.value = order_id
  
  if (queryStatus) {
    status.value = queryStatus
  } else {
    checkPaymentStatus()
  }
  
  if (order_id) {
    loadOrder(order_id)
  }
})

const checkPaymentStatus = async () => {
  try {
    const res = await uni.request({
      url: `${getApp().globalData.apiBaseUrl}/orders/${orderId.value}`
    })
    
    const orderData = res.data.data
    if (orderData.status === 'paid' || orderData.status === 'confirmed') {
      status.value = 'success'
    } else {
      status.value = 'failed'
    }
  } catch (err) {
    status.value = 'failed'
  }
}

const loadOrder = async (id) => {
  try {
    const res = await uni.request({
      url: `${getApp().globalData.apiBaseUrl}/orders/${id}`
    })
    order.value = res.data.data
  } catch (err) {
    console.error('Failed to load order:', err)
  }
}

const goToOrder = () => {
  uni.redirectTo({
    url: `/pages/orders/detail?id=${orderId.value}`
  })
}

const goHome = () => {
  uni.switchTab({
    url: '/pages/index/index'
  })
}

const retryPayment = () => {
  uni.redirectTo({
    url: `/pages/booking/index?order_id=${orderId.value}`
  })
}

const goToOrders = () => {
  uni.switchTab({
    url: '/pages/orders/index'
  })
}
</script>

<style>
.payment-result {
  min-height: 100vh;
  background: #fff;
  padding: 60rpx 40rpx;
}

.result {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.icon {
  width: 200rpx;
  height: 200rpx;
  margin-bottom: 40rpx;
}

.title {
  font-size: 48rpx;
  font-weight: bold;
  margin-bottom: 20rpx;
}

.success .title {
  color: #34c759;
}

.failed .title {
  color: #ff3b30;
}

.subtitle {
  font-size: 30rpx;
  color: #666;
  margin-bottom: 60rpx;
}

.order-info {
  width: 100%;
  background: #f5f5f5;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 60rpx;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 20rpx 0;
  font-size: 30rpx;
}

.info-row:not(:last-child) {
  border-bottom: 1rpx solid #e5e5e5;
}

.price {
  color: #ff6b6b;
  font-weight: bold;
  font-size: 36rpx;
}

.actions {
  width: 100%;
}

.btn-primary, .btn-secondary {
  width: 100%;
  height: 96rpx;
  border-radius: 48rpx;
  font-size: 32rpx;
  margin-bottom: 20rpx;
}

.btn-primary {
  background: #007AFF;
  color: #fff;
}

.btn-secondary {
  background: #f5f5f5;
  color: #666;
}

.loading-spinner {
  width: 80rpx;
  height: 80rpx;
  border: 6rpx solid #007AFF;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 40rpx;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
