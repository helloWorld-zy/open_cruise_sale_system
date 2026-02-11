<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Hero Section -->
    <section class="relative bg-gradient-to-r from-blue-600 to-blue-800 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
        <div class="text-center">
          <h1 class="text-4xl md:text-6xl font-bold mb-6">
            探索豪华邮轮之旅
          </h1>
          <p class="text-xl md:text-2xl mb-8 text-blue-100">
            精选全球顶级邮轮航线，开启您的海上假期
          </p>
          <UButton
            size="xl"
            color="white"
            variant="solid"
            to="/cruises"
          >
            浏览邮轮
            <template #trailing>
              <UIcon name="i-heroicons-arrow-right" />
            </template>
          </UButton>
        </div>
      </div>
    </section>

    <!-- Featured Cruises -->
    <section class="py-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center mb-8">
          <h2 class="text-3xl font-bold text-gray-900">热门邮轮</h2>
          <NuxtLink to="/cruises" class="text-blue-600 hover:text-blue-800 font-medium">
            查看全部 →
          </NuxtLink>
        </div>
        
        <div v-if="loading" class="flex justify-center py-12">
          <Loading message="加载中..." />
        </div>
        
        <div v-else-if="error" class="py-12">
          <Error :message="error" @retry="fetchCruises" />
        </div>
        
        <div v-else-if="cruises.length === 0" class="py-12">
          <Empty title="暂无邮轮" description="敬请期待更多邮轮航线" />
        </div>
        
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <CruiseCard
            v-for="cruise in cruises.slice(0, 6)"
            :key="cruise.id"
            :cruise="cruise"
          />
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="py-16 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 class="text-3xl font-bold text-gray-900 text-center mb-12">为什么选择我们</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="text-center">
            <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <UIcon name="i-heroicons-ship" class="w-8 h-8 text-blue-600" />
            </div>
            <h3 class="text-xl font-semibold mb-2">精选邮轮</h3>
            <p class="text-gray-600">严选全球顶级邮轮，提供最优质的邮轮体验</p>
          </div>
          <div class="text-center">
            <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <UIcon name="i-heroicons-currency-dollar" class="w-8 h-8 text-blue-600" />
            </div>
            <h3 class="text-xl font-semibold mb-2">价格透明</h3>
            <p class="text-gray-600">明码标价，无隐藏费用，预订更放心</p>
          </div>
          <div class="text-center">
            <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <UIcon name="i-heroicons-shield-check" class="w-8 h-8 text-blue-600" />
            </div>
            <h3 class="text-xl font-semibold mb-2">安全保障</h3>
            <p class="text-gray-600">专业客服团队，全程保障您的出行安全</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { $api } = useNuxtApp()

const cruises = ref([])
const loading = ref(true)
const error = ref('')

const fetchCruises = async () => {
  try {
    loading.value = true
    error.value = ''
    const response = await $api('/cruises')
    cruises.value = response.data || []
  } catch (err: any) {
    error.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCruises()
})
</script>
