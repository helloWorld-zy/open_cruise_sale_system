<template>
  <div>
    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center py-20">
      <Loading message="加载邮轮详情..." />
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="py-20">
      <Error :message="error" @retry="fetchCruise" />
    </div>

    <!-- Cruise Detail -->
    <div v-else-if="cruise" class="min-h-screen bg-gray-50">
      <!-- Hero Section with Gallery -->
      <div class="bg-white">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <!-- Image Gallery -->
            <ImageGallery
              :images="cruise.coverImages || []"
              :title="cruise.nameCn"
            />

            <!-- Basic Info -->
            <div class="space-y-6">
              <div>
                <UBadge
                  :color="cruise.status === 'active' ? 'green' : 'gray'"
                  class="mb-3"
                >
                  {{ statusText }}
                </UBadge>
                <h1 class="text-3xl font-bold text-gray-900">{{ cruise.nameCn }}</h1>
                <p v-if="cruise.nameEn" class="text-lg text-gray-500 mt-1">
                  {{ cruise.nameEn }}
                </p>
              </div>

              <!-- Ship Specs -->
              <div class="grid grid-cols-2 gap-4">
                <div class="bg-gray-50 rounded-lg p-4">
                  <div class="text-sm text-gray-500">总吨位</div>
                  <div class="text-lg font-semibold text-gray-900">
                    {{ cruise.grossTonnage?.toLocaleString() || '-' }} 吨
                  </div>
                </div>
                <div class="bg-gray-50 rounded-lg p-4">
                  <div class="text-sm text-gray-500">载客量</div>
                  <div class="text-lg font-semibold text-gray-900">
                    {{ cruise.passengerCapacity?.toLocaleString() || '-' }} 人
                  </div>
                </div>
                <div class="bg-gray-50 rounded-lg p-4">
                  <div class="text-sm text-gray-500">甲板层数</div>
                  <div class="text-lg font-semibold text-gray-900">
                    {{ cruise.deckCount || '-' }} 层
                  </div>
                </div>
                <div class="bg-gray-50 rounded-lg p-4">
                  <div class="text-sm text-gray-500">船员人数</div>
                  <div class="text-lg font-semibold text-gray-900">
                    {{ cruise.crewCount?.toLocaleString() || '-' }} 人
                  </div>
                </div>
              </div>

              <!-- Company Info -->
              <div v-if="cruise.company" class="flex items-center space-x-3 pt-4 border-t border-gray-200">
                <div class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
                  <UIcon name="i-heroicons-building-office" class="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <div class="text-sm text-gray-500">邮轮公司</div>
                  <div class="font-medium text-gray-900">{{ cruise.company.name }}</div>
                </div>
              </div>

              <!-- Action Buttons -->
              <div class="flex gap-4 pt-4">
                <UButton color="primary" size="lg" block @click="scrollToCabins">
                  查看舱房
                </UButton>
                <UButton color="gray" variant="outline" size="lg" block @click="goBack">
                  返回列表
                </UButton>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Details Tabs -->
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <UTabs :items="tabItems" class="w-full">
          <template #cabins>
            <div id="cabins-section" class="py-6">
              <h2 class="text-2xl font-bold text-gray-900 mb-6">舱房类型</h2>
              <CabinTypeAccordion
                v-if="cruise.cabinTypes && cruise.cabinTypes.length > 0"
                :cabin-types="cruise.cabinTypes"
              />
              <Empty v-else title="暂无舱房信息" description="该邮轮暂时没有可预订的舱房" />
            </div>
          </template>

          <template #facilities>
            <div class="py-6">
              <h2 class="text-2xl font-bold text-gray-900 mb-6">船上设施</h2>
              <FacilityTabs
                v-if="cruise.facilities && cruise.facilities.length > 0"
                :facilities="cruise.facilities"
              />
              <Empty v-else title="暂无设施信息" description="该邮轮暂时没有设施信息" />
            </div>
          </template>

          <template #specs>
            <div class="py-6">
              <h2 class="text-2xl font-bold text-gray-900 mb-6">邮轮参数</h2>
              <div class="bg-white rounded-lg shadow-sm">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-0">
                  <div class="px-6 py-4 border-b border-gray-200 md:border-r">
                    <span class="text-gray-500">邮轮名称</span>
                    <span class="float-right font-medium">{{ cruise.nameCn }}</span>
                  </div>
                  <div class="px-6 py-4 border-b border-gray-200">
                    <span class="text-gray-500">英文名称</span>
                    <span class="float-right font-medium">{{ cruise.nameEn || '-' }}</span>
                  </div>
                  <div class="px-6 py-4 border-b border-gray-200 md:border-r">
                    <span class="text-gray-500">邮轮编号</span>
                    <span class="float-right font-medium">{{ cruise.code }}</span>
                  </div>
                  <div class="px-6 py-4 border-b border-gray-200">
                    <span class="text-gray-500">总吨位</span>
                    <span class="float-right font-medium">{{ cruise.grossTonnage?.toLocaleString() || '-' }} 吨</span>
                  </div>
                  <div class="px-6 py-4 border-b border-gray-200 md:border-r">
                    <span class="text-gray-500">载客量</span>
                    <span class="float-right font-medium">{{ cruise.passengerCapacity?.toLocaleString() || '-' }} 人</span>
                  </div>
                  <div class="px-6 py-4 border-b border-gray-200">
                    <span class="text-gray-500">船员数</span>
                    <span class="float-right font-medium">{{ cruise.crewCount?.toLocaleString() || '-' }} 人</span>
                  </div>
                  <div class="px-6 py-4 border-b border-gray-200 md:border-r">
                    <span class="text-gray-500">建造年份</span>
                    <span class="float-right font-medium">{{ cruise.builtYear || '-' }} 年</span>
                  </div>
                  <div class="px-6 py-4 border-b border-gray-200">
                    <span class="text-gray-500">翻新年份</span>
                    <span class="float-right font-medium">{{ cruise.renovatedYear || '-' }} 年</span>
                  </div>
                  <div class="px-6 py-4 border-b border-gray-200 md:border-r">
                    <span class="text-gray-500">长度</span>
                    <span class="float-right font-medium">{{ cruise.lengthMeters || '-' }} 米</span>
                  </div>
                  <div class="px-6 py-4">
                    <span class="text-gray-500">宽度</span>
                    <span class="float-right font-medium">{{ cruise.widthMeters || '-' }} 米</span>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </UTabs>
      </div>
    </div>

    <!-- Not Found -->
    <div v-else class="py-20">
      <Empty title="邮轮不存在" description="抱歉，您访问的邮轮信息不存在或已被删除">
        <template #action>
          <UButton color="primary" to="/cruises">
            浏览其他邮轮
          </UButton>
        </template>
      </Empty>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()

const cruise = ref<any>(null)
const loading = ref(true)
const error = ref('')

const tabItems = [
  {
    label: '舱房类型',
    slot: 'cabins'
  },
  {
    label: '船上设施',
    slot: 'facilities'
  },
  {
    label: '邮轮参数',
    slot: 'specs'
  }
]

const statusText = computed(() => {
  const statusMap: Record<string, string> = {
    active: '热售中',
    inactive: '已下架',
    maintenance: '维护中'
  }
  return statusMap[cruise.value?.status] || cruise.value?.status
})

const fetchCruise = async () => {
  try {
    loading.value = true
    error.value = ''
    const id = route.params.id as string
    const response = await $api(`/cruises/${id}`)
    cruise.value = response
  } catch (err: any) {
    error.value = err.message || '加载失败'
    if (err.statusCode === 404) {
      cruise.value = null
    }
  } finally {
    loading.value = false
  }
}

const scrollToCabins = () => {
  const element = document.getElementById('cabins-section')
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  fetchCruise()
})
</script>
