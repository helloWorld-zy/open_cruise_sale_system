<template>
  <view class="booking-page">
    <!-- Progress Steps -->
    <view class="progress-bar">
      <view 
        v-for="(step, index) in steps" 
        :key="index"
        class="step"
        :class="{ active: currentStep >= index, current: currentStep === index }"
      >
        <view class="step-dot">{{ index + 1 }}</view>
        <text class="step-text">{{ step }}</text>
      </view>
    </view>

    <!-- Step 1: Select Voyage -->
    <view v-if="currentStep === 0" class="step-content">
      <view class="section-title">选择出发日期</view>
      <scroll-view scroll-y class="voyage-list">
        <view 
          v-for="voyage in voyages" 
          :key="voyage.id"
          class="voyage-card"
          :class="{ selected: selectedVoyage?.id === voyage.id }"
          @click="selectVoyage(voyage)"
        >
          <view class="voyage-date">
            <text class="day">{{ formatDay(voyage.departure_date) }}</text>
            <text class="month">{{ formatMonth(voyage.departure_date) }}</text>
          </view>
          <view class="voyage-info">
            <text class="route">{{ voyage.route?.name }}</text>
            <text class="duration">{{ voyage.route?.duration_days }}天航程</text>
            <text class="time">{{ voyage.departure_date }} 出发</text>
          </view>
          <view class="voyage-price" v-if="voyage.min_price">
            <text class="price">¥{{ voyage.min_price }}</text>
            <text class="unit">起/人</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- Step 2: Select Cabin -->
    <view v-if="currentStep === 1" class="step-content">
      <view class="section-title">选择舱房类型</view>
      <scroll-view scroll-y class="cabin-list">
        <view 
          v-for="cabin in cabins" 
          :key="cabin.id"
          class="cabin-card"
        >
          <view class="cabin-header">
            <view class="cabin-type">
              <text class="name">{{ cabin.cabin_type?.name_cn }}</text>
              <text class="name-en">{{ cabin.cabin_type?.name_en }}</text>
            </view>
            <view class="cabin-price" v-if="cabin.price">
              <text class="price">¥{{ cabin.price.adult_price }}</text>
              <text class="unit">/人</text>
            </view>
          </view>
          
          <view class="passenger-selector">
            <view class="selector-row">
              <text>成人(12岁+)</text>
              <view class="counter">
                <button @click="decrement(cabin.id, 'adult')" :disabled="getCount(cabin.id, 'adult') <= 1">-</button>
                <text>{{ getCount(cabin.id, 'adult') }}</text>
                <button @click="increment(cabin.id, 'adult')">+</button>
              </view>
            </view>
            <view class="selector-row">
              <text>儿童(2-11岁)</text>
              <view class="counter">
                <button @click="decrement(cabin.id, 'child')" :disabled="getCount(cabin.id, 'child') <= 0">-</button>
                <text>{{ getCount(cabin.id, 'child') }}</text>
                <button @click="increment(cabin.id, 'child')">+</button>
              </view>
            </view>
          </view>

          <button 
            class="select-btn"
            :class="{ selected: isSelected(cabin.id) }"
            @click="toggleCabin(cabin)"
          >
            {{ isSelected(cabin.id) ? '已选择' : '选择此舱房' }}
          </button>
        </view>
      </scroll-view>
    </view>

    <!-- Step 3: Passenger Info -->
    <view v-if="currentStep === 2" class="step-content">
      <view class="section-title">填写乘客信息</view>
      <scroll-view scroll-y class="passenger-list">
        <view class="contact-section">
          <view class="input-group">
            <text class="label required">联系人姓名</text>
            <input v-model="contact.name" placeholder="请输入联系人姓名" />
          </view>
          <view class="input-group">
            <text class="label required">手机号码</text>
            <input v-model="contact.phone" type="number" placeholder="请输入手机号码" />
          </view>
        </view>

        <view 
          v-for="(passenger, index) in passengers" 
          :key="index"
          class="passenger-card"
        >
          <view class="passenger-header">
            <text class="passenger-title">乘客 {{ index + 1 }}</text>
            <text class="passenger-type">{{ passenger.type === 'adult' ? '成人' : '儿童' }}</text>
          </view>
          <view class="input-group">
            <text class="label required">姓(拼音)</text>
            <input v-model="passenger.surname" placeholder="如: ZHANG" />
          </view>
          <view class="input-group">
            <text class="label required">名(拼音)</text>
            <input v-model="passenger.givenName" placeholder="如: SAN" />
          </view>
          <view class="input-group">
            <text class="label required">中文姓名</text>
            <input v-model="passenger.name" placeholder="如: 张三" />
          </view>
          <view class="input-group">
            <text class="label required">性别</text>
            <picker mode="selector" :range="['男', '女']" @change="(e) => passenger.gender = e.detail.value === '0' ? 'male' : 'female'">
              <view class="picker">
                {{ passenger.gender === 'male' ? '男' : passenger.gender === 'female' ? '女' : '请选择' }}
              </view>
            </picker>
          </view>
          <view class="input-group">
            <text class="label required">出生日期</text>
            <picker mode="date" @change="(e) => passenger.birthDate = e.detail.value">
              <view class="picker">{{ passenger.birthDate || '请选择' }}</view>
            </picker>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- Step 4: Confirm Order -->
    <view v-if="currentStep === 3" class="step-content">
      <view class="section-title">确认订单</view>
      <scroll-view scroll-y class="confirm-content">
        <view class="info-section">
          <text class="section-label">航次信息</text>
          <view class="info-row">
            <text>出发日期: {{ selectedVoyage?.departure_date }}</text>
          </view>
          <view class="info-row">
            <text>航线: {{ selectedVoyage?.route?.name }}</text>
          </view>
        </view>

        <view class="info-section">
          <text class="section-label">舱房信息</text>
          <view v-for="(cabin, index) in selectedCabins" :key="index" class="info-row">
            <text>{{ cabin.cabin_type?.name_cn }} - {{ cabin.adultCount }}成人</text>
            <text v-if="cabin.childCount">{{ cabin.childCount }}儿童</text>
          </view>
        </view>

        <view class="price-summary">
          <view class="price-row">
            <text>订单总额</text>
            <text class="total-price">¥{{ totalPrice }}</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- Bottom Actions -->
    <view class="bottom-actions">
      <button v-if="currentStep > 0" class="btn-prev" @click="prevStep">上一步</button>
      <button 
        class="btn-next" 
        :disabled="!canProceed"
        @click="nextStep"
      >
        {{ currentStep === 3 ? '立即支付' : '下一步' }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const props = defineProps({
  cruiseId: String
})

const currentStep = ref(0)
const steps = ['选择航次', '选择舱房', '填写信息', '确认订单']

// Data
const voyages = ref([])
const cabins = ref([])
const selectedVoyage = ref(null)
const selectedCabins = ref([])
const passengerCounts = ref({})
const passengers = ref([])

const contact = ref({
  name: '',
  phone: ''
})

onMounted(() => {
  loadVoyages()
})

const loadVoyages = async () => {
  try {
    const res = await uni.request({
      url: `${getApp().globalData.apiBaseUrl}/voyages`,
      data: { cruise_id: props.cruiseId, booking_status: 'open' }
    })
    voyages.value = res.data.data || []
  } catch (err) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  }
}

const selectVoyage = (voyage) => {
  selectedVoyage.value = voyage
  loadCabins()
}

const loadCabins = async () => {
  try {
    const [cabinsRes, pricesRes] = await Promise.all([
      uni.request({
        url: `${getApp().globalData.apiBaseUrl}/cabins`,
        data: { voyage_id: selectedVoyage.value.id }
      }),
      uni.request({
        url: `${getApp().globalData.apiBaseUrl}/prices`,
        data: { voyage_id: selectedVoyage.value.id }
      })
    ])

    const cabinData = cabinsRes.data.data || []
    const prices = pricesRes.data.data || []

    cabins.value = cabinData.map(cabin => ({
      ...cabin,
      price: prices.find(p => p.cabin_type_id === cabin.cabin_type_id)
    }))
  } catch (err) {
    uni.showToast({ title: '加载舱房失败', icon: 'none' })
  }
}

const getCount = (cabinId, type) => {
  return passengerCounts.value[cabinId]?.[type] || (type === 'adult' ? 2 : 0)
}

const increment = (cabinId, type) => {
  if (!passengerCounts.value[cabinId]) {
    passengerCounts.value[cabinId] = { adult: 2, child: 0 }
  }
  passengerCounts.value[cabinId][type]++
}

const decrement = (cabinId, type) => {
  if (passengerCounts.value[cabinId] && passengerCounts.value[cabinId][type] > (type === 'adult' ? 1 : 0)) {
    passengerCounts.value[cabinId][type]--
  }
}

const isSelected = (cabinId) => {
  return selectedCabins.value.some(c => c.id === cabinId)
}

const toggleCabin = (cabin) => {
  const index = selectedCabins.value.findIndex(c => c.id === cabin.id)
  if (index >= 0) {
    selectedCabins.value.splice(index, 1)
  } else {
    const counts = passengerCounts.value[cabin.id] || { adult: 2, child: 0 }
    selectedCabins.value.push({
      ...cabin,
      adultCount: counts.adult,
      childCount: counts.child
    })
  }
}

const initPassengers = () => {
  passengers.value = []
  selectedCabins.value.forEach(cabin => {
    for (let i = 0; i < cabin.adultCount; i++) {
      passengers.value.push({
        type: 'adult',
        surname: '',
        givenName: '',
        name: '',
        gender: '',
        birthDate: ''
      })
    }
    for (let i = 0; i < cabin.childCount; i++) {
      passengers.value.push({
        type: 'child',
        surname: '',
        givenName: '',
        name: '',
        gender: '',
        birthDate: ''
      })
    }
  })
}

const createOrder = async () => {
  try {
    const res = await uni.request({
      url: `${getApp().globalData.apiBaseUrl}/orders`,
      method: 'POST',
      data: {
        cruise_id: props.cruiseId,
        voyage_id: selectedVoyage.value.id,
        items: selectedCabins.value.map(cabin => ({
          cabin_id: cabin.id,
          cabin_type_id: cabin.cabin_type_id,
          adult_count: cabin.adultCount,
          child_count: cabin.childCount
        })),
        passengers: passengers.value.map(p => ({
          name: p.name,
          surname: p.surname,
          given_name: p.givenName,
          gender: p.gender,
          birth_date: p.birthDate,
          passenger_type: p.type
        })),
        contact_name: contact.value.name,
        contact_phone: contact.value.phone
      }
    })

    if (res.data.data) {
      // Navigate to WeChat payment
      const orderId = res.data.data.id
      requestWechatPayment(orderId)
    }
  } catch (err) {
    uni.showToast({ title: '创建订单失败', icon: 'none' })
  }
}

const requestWechatPayment = (orderId) => {
  uni.requestPayment({
    provider: 'wxpay',
    orderInfo: {
      orderId: orderId
    },
    success: () => {
      uni.redirectTo({
        url: `/pages/payment/result?order_id=${orderId}&status=success`
      })
    },
    fail: () => {
      uni.redirectTo({
        url: `/pages/payment/result?order_id=${orderId}&status=failed`
      })
    }
  })
}

const nextStep = () => {
  if (currentStep.value === 2) {
    initPassengers()
  }
  if (currentStep.value === 3) {
    createOrder()
    return
  }
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const canProceed = computed(() => {
  switch (currentStep.value) {
    case 0:
      return selectedVoyage.value !== null
    case 1:
      return selectedCabins.value.length > 0
    case 2:
      return contact.value.name && contact.value.phone && 
             passengers.value.every(p => p.name && p.surname && p.givenName && p.gender && p.birthDate)
    case 3:
      return true
    default:
      return false
  }
})

const totalPrice = computed(() => {
  return selectedCabins.value.reduce((total, cabin) => {
    if (!cabin.price) return total
    const adultTotal = cabin.price.adult_price * cabin.adultCount
    const childTotal = (cabin.price.child_price || 0) * cabin.childCount
    const fees = (cabin.price.port_fee + cabin.price.service_fee) * (cabin.adultCount + cabin.childCount)
    return total + adultTotal + childTotal + fees
  }, 0)
})

const formatDay = (dateStr) => {
  return new Date(dateStr).getDate()
}

const formatMonth = (dateStr) => {
  return (new Date(dateStr).getMonth() + 1) + '月'
}
</script>

<style>
.booking-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 120rpx;
}

.progress-bar {
  display: flex;
  justify-content: space-around;
  padding: 30rpx;
  background: #fff;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.step-dot {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background: #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  color: #fff;
}

.step.active .step-dot {
  background: #007AFF;
}

.step.current .step-dot {
  background: #007AFF;
  transform: scale(1.1);
}

.step-text {
  font-size: 24rpx;
  margin-top: 10rpx;
  color: #999;
}

.step.active .step-text {
  color: #007AFF;
}

.step-content {
  padding: 20rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  padding: 20rpx 0;
}

.voyage-list, .cabin-list, .passenger-list, .confirm-content {
  max-height: calc(100vh - 300rpx);
}

.voyage-card, .cabin-card, .passenger-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  display: flex;
  align-items: center;
}

.voyage-card.selected {
  border: 2rpx solid #007AFF;
}

.voyage-date {
  width: 100rpx;
  text-align: center;
  background: #f0f7ff;
  border-radius: 12rpx;
  padding: 20rpx;
}

.voyage-date .day {
  font-size: 48rpx;
  font-weight: bold;
  color: #007AFF;
}

.voyage-date .month {
  font-size: 24rpx;
  color: #007AFF;
}

.voyage-info {
  flex: 1;
  margin-left: 20rpx;
}

.voyage-info .route {
  font-size: 32rpx;
  font-weight: bold;
  display: block;
}

.voyage-info .duration, .voyage-info .time {
  font-size: 26rpx;
  color: #666;
  display: block;
  margin-top: 10rpx;
}

.voyage-price {
  text-align: right;
}

.voyage-price .price {
  font-size: 40rpx;
  color: #ff6b6b;
  font-weight: bold;
}

.voyage-price .unit {
  font-size: 24rpx;
  color: #999;
}

.cabin-card {
  flex-direction: column;
  align-items: stretch;
}

.cabin-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.cabin-type .name {
  font-size: 32rpx;
  font-weight: bold;
  display: block;
}

.cabin-type .name-en {
  font-size: 24rpx;
  color: #999;
}

.cabin-price .price {
  font-size: 40rpx;
  color: #ff6b6b;
  font-weight: bold;
}

.passenger-selector {
  margin: 20rpx 0;
}

.selector-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #eee;
}

.counter {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.counter button {
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
  border: 1rpx solid #ddd;
  background: #fff;
  font-size: 32rpx;
}

.counter button[disabled] {
  opacity: 0.3;
}

.select-btn {
  width: 100%;
  height: 80rpx;
  background: #007AFF;
  color: #fff;
  border-radius: 40rpx;
  font-size: 30rpx;
  margin-top: 20rpx;
}

.select-btn.selected {
  background: #34c759;
}

.contact-section, .passenger-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.passenger-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.passenger-title {
  font-size: 32rpx;
  font-weight: bold;
}

.passenger-type {
  font-size: 26rpx;
  color: #007AFF;
  background: #f0f7ff;
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
}

.input-group {
  margin-bottom: 20rpx;
}

.input-group .label {
  font-size: 28rpx;
  color: #333;
  margin-bottom: 10rpx;
  display: block;
}

.input-group .label.required::after {
  content: ' *';
  color: #ff6b6b;
}

.input-group input, .picker {
  height: 80rpx;
  background: #f5f5f5;
  border-radius: 12rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
}

.info-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.section-label {
  font-size: 28rpx;
  font-weight: bold;
  margin-bottom: 20rpx;
  display: block;
}

.info-row {
  font-size: 28rpx;
  color: #666;
  padding: 10rpx 0;
}

.price-summary {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
}

.price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.total-price {
  font-size: 48rpx;
  color: #ff6b6b;
  font-weight: bold;
}

.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  padding: 20rpx 30rpx;
  display: flex;
  gap: 20rpx;
  box-shadow: 0 -2rpx 10rpx rgba(0,0,0,0.1);
}

.btn-prev, .btn-next {
  flex: 1;
  height: 88rpx;
  border-radius: 44rpx;
  font-size: 32rpx;
}

.btn-prev {
  background: #f5f5f5;
  color: #666;
}

.btn-next {
  background: #007AFF;
  color: #fff;
}

.btn-next[disabled] {
  opacity: 0.5;
}
</style>
