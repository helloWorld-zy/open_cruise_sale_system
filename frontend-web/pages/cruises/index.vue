<template>
  <div>
    <!-- Page Header -->
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 class="text-3xl font-bold text-gray-900">邮轮列表</h1>
        <p class="mt-2 text-gray-600">探索豪华邮轮，开启精彩旅程</p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Filters -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-8">
        <div class="flex flex-wrap gap-4">
          <UInput
            v-model="filters.keyword"
            icon="i-heroicons-magnifying-glass"
            placeholder="搜索邮轮名称"
            class="w-full md:w-64"
          />
          <USelect
            v-model="filters.status"
            :options="statusOptions"
            placeholder="状态"
            class="w-full md:w-40"
          />
          <UButton color="primary" @click="applyFilters">
            搜索
          </UButton>
          <UButton color="gray" variant="ghost" @click="resetFilters">
            重置
          </UButton>
        </div>
      </div>

      <!-- Cruise List -->
      <div v-if="loading" class="flex justify-center py-12">
        <Loading message="加载中..." />
      </div>

      <div v-else-if="error" class="py-12">
        <Error :message="error" @retry="fetchCruises" />
      </div>

      <div v-else-if="cruises.length === 0" class="py-12">
        <Empty title="暂无邮轮" description="请调整筛选条件或稍后再试" />
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <CruiseCard
          v-for="cruise in cruises"
          :key="cruise.id"
          :cruise="cruise"
        />
      </div>

      <!-- Pagination -->
      <div v-if="cruises.length > 0" class="mt-8 flex justify-center">
        <UPagination
          v-model="page"
          :total="total"
          :page-count="pageSize"
          @update:model-value="fetchCruises"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { $api } = useNuxtApp()

const cruises = ref([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)

const filters = reactive({
  keyword: '',
  status: ''
})

const statusOptions = [
  { label: '全部', value: '' },
  { label: '上架', value: 'active' },
  { label: '下架', value: 'inactive' }
]

const fetchCruises = async () => {
  try {
    loading.value = true
    error.value = ''
    const response = await $api('/cruises', {
      query: {
        page: page.value,
        page_size: pageSize.value,
        keyword: filters.keyword || undefined,
        status: filters.status || undefined
      }
    })
    cruises.value = response.data || []
    total.value = response.pagination?.total || 0
  } catch (err: any) {
    error.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
}

const applyFilters = () => {
  page.value = 1
  fetchCruises()
}

const resetFilters = () => {
  filters.keyword = ''
  filters.status = ''
  page.value = 1
  fetchCruises()
}

onMounted(() => {
  fetchCruises()
})
</script>
