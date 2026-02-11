<template>
  <view class="cruise-card" @click="handleClick">
    <image 
      class="card-image" 
      :src="cruise.coverImages?.[0] || '/static/images/placeholder.jpg'" 
      mode="aspectFill"
    />
    <view class="card-content">
      <view class="card-header">
        <text class="cruise-name">{{ cruise.nameCn }}</text>
        <view class="status-badge" :class="cruise.status">
          {{ statusText }}
        </view>
      </view>
      
      <text v-if="cruise.nameEn" class="cruise-name-en">{{ cruise.nameEn }}</text>
      
      <view class="cruise-specs">
        <view class="spec-item">
          <text class="iconfont icon-scale"></text>
          <text>{{ formatNumber(cruise.grossTonnage) }} 吨</text>
        </view>
        <view class="spec-item">
          <text class="iconfont icon-users"></text>
          <text>{{ formatNumber(cruise.passengerCapacity) }} 人</text>
        </view>
        <view class="spec-item">
          <text class="iconfont icon-deck"></text>
          <text>{{ cruise.deckCount }} 层甲板</text>
        </view>
      </view>
      
      <view class="card-footer">
        <view class="company-info" v-if="cruise.company">
          <text class="company-name">{{ cruise.company.name }}</text>
        </view>
        <view class="view-btn">
          <text>查看详情</text>
          <text class="iconfont icon-arrow-right"></text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
const props = defineProps<{
  cruise: {
    id: string
    nameCn: string
    nameEn?: string
    grossTonnage?: number
    passengerCapacity?: number
    deckCount?: number
    status: string
    coverImages?: string[]
    company?: {
      name: string
    }
  }
}>()

const emit = defineEmits<{
  click: [id: string]
}>()

const statusText = computed(() => {
  const map: Record<string, string> = {
    active: '热售中',
    inactive: '已下架',
    maintenance: '维护中'
  }
  return map[props.cruise.status] || props.cruise.status
})

const formatNumber = (num?: number) => {
  if (!num) return '-'
  return num.toLocaleString()
}

const handleClick = () => {
  emit('click', props.cruise.id)
}
</script>

<style scoped>
.cruise-card {
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.08);
}

.card-image {
  width: 100%;
  height: 300rpx;
}

.card-content {
  padding: 24rpx;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8rpx;
}

.cruise-name {
  font-size: 32rpx;
  font-weight: bold;
  color: #1f2937;
  flex: 1;
  margin-right: 16rpx;
}

.status-badge {
  padding: 6rpx 12rpx;
  border-radius: 8rpx;
  font-size: 20rpx;
  flex-shrink: 0;
}

.status-badge.active {
  background: #dcfce7;
  color: #166534;
}

.status-badge.inactive {
  background: #f3f4f6;
  color: #6b7280;
}

.status-badge.maintenance {
  background: #fef3c7;
  color: #92400e;
}

.cruise-name-en {
  font-size: 24rpx;
  color: #9ca3af;
  margin-bottom: 16rpx;
}

.cruise-specs {
  display: flex;
  gap: 24rpx;
  margin-bottom: 20rpx;
}

.spec-item {
  display: flex;
  align-items: center;
  font-size: 24rpx;
  color: #6b7280;
}

.spec-item text:first-child {
  margin-right: 8rpx;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16rpx;
  border-top: 1rpx solid #f3f4f6;
}

.company-name {
  font-size: 24rpx;
  color: #6b7280;
}

.view-btn {
  display: flex;
  align-items: center;
  color: #3b82f6;
  font-size: 26rpx;
}

.view-btn text:last-child {
  margin-left: 8rpx;
}
</style>
