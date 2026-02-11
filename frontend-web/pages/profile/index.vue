<template>
  <div class="profile-page">
    <div class="container">
      <h1 class="text-2xl font-bold mb-6">个人中心</h1>

      <!-- User Info Card -->
      <div class="profile-card">
        <div class="avatar-section">
          <img :src="user?.avatar_url || '/default-avatar.png'" class="avatar" />
          <div class="user-info">
            <h2 class="name">{{ user?.nickname || user?.phone || '未设置昵称' }}</h2>
            <p class="phone">{{ user?.phone }}</p>
          </div>
        </div>
        <div class="stats">
          <div class="stat">
            <span class="value">{{ orderStats?.total_orders || 0 }}</span>
            <span class="label">全部订单</span>
          </div>
          <div class="stat">
            <span class="value">{{ passengerCount }}</span>
            <span class="label">常用乘客</span>
          </div>
        </div>
      </div>

      <!-- Menu -->
      <div class="menu-list">
        <NuxtLink to="/orders" class="menu-item">
          <span>我的订单</span>
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
          </svg>
        </NuxtLink>
        <NuxtLink to="/profile/passengers" class="menu-item">
          <span>常用乘客管理</span>
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
          </svg>
        </NuxtLink>
        <button class="menu-item text-red-500" @click="logout">
          <span>退出登录</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const user = ref<any>(null)
const orderStats = ref<any>(null)
const passengers = ref<any[]>([])

const passengerCount = computed(() => passengers.value.length)

onMounted(async () => {
  const [profileRes, statsRes, passengersRes] = await Promise.all([
    $fetch('/api/user/profile'),
    $fetch('/api/orders/statistics'),
    $fetch('/api/user/passengers')
  ])
  
  user.value = profileRes.data
  orderStats.value = statsRes.data
  passengers.value = passengersRes.data || []
})

const logout = () => {
  if (confirm('确定要退出登录吗？')) {
    // Clear auth tokens
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    navigateTo('/')
  }
}
</script>

<style scoped>
.profile-page {
  @apply min-h-screen bg-gray-50 py-8;
}

.container {
  @apply max-w-2xl mx-auto px-4;
}

.profile-card {
  @apply bg-white rounded-2xl p-6 shadow-sm mb-6;
}

.avatar-section {
  @apply flex items-center gap-4 mb-4;
}

.avatar {
  @apply w-16 h-16 rounded-full object-cover;
}

.user-info .name {
  @apply font-bold text-lg;
}

.user-info .phone {
  @apply text-gray-500 text-sm;
}

.stats {
  @apply flex gap-8 pt-4 border-t;
}

.stat {
  @apply flex flex-col;
}

.stat .value {
  @apply text-xl font-bold text-blue-500;
}

.stat .label {
  @apply text-sm text-gray-500;
}

.menu-list {
  @apply bg-white rounded-2xl shadow-sm overflow-hidden;
}

.menu-item {
  @apply flex justify-between items-center px-6 py-4 border-b last:border-0 hover:bg-gray-50;
}

.menu-item span {
  @apply font-medium;
}
</style>
