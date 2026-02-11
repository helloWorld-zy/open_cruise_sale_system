<template>
  <div>
    <!-- Page Header -->
    <div class="md:flex md:items-center md:justify-between mb-8">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900">
          {{ isEdit ? '编辑邮轮' : '新增邮轮' }}
        </h2>
      </div>
      <div class="mt-4 flex md:ml-4 md:mt-0">
        <UButton color="gray" variant="ghost" @click="$router.back()">
          返回
        </UButton>
      </div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <!-- Basic Information -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-medium text-gray-900">基本信息</h3>
          </template>
          
          <div class="space-y-4">
            <UFormGroup label="邮轮公司" required>
              <USelect
                v-model="form.companyId"
                :options="companyOptions"
                placeholder="选择邮轮公司"
              />
            </UFormGroup>

            <UFormGroup label="中文名称" required>
              <UInput v-model="form.nameCn" placeholder="输入邮轮中文名称" />
            </UFormGroup>

            <UFormGroup label="英文名称">
              <UInput v-model="form.nameEn" placeholder="输入邮轮英文名称" />
            </UFormGroup>

            <UFormGroup label="邮轮编号" required>
              <UInput v-model="form.code" placeholder="输入唯一编号" :disabled="isEdit" />
            </UFormGroup>

            <UFormGroup label="状态">
              <USelect
                v-model="form.status"
                :options="statusOptions"
              />
            </UFormGroup>
          </div>
        </UCard>

        <!-- Ship Specifications -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-medium text-gray-900">邮轮参数</h3>
          </template>
          
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="总吨位">
                <UInput v-model.number="form.grossTonnage" type="number" placeholder="吨" />
              </UFormGroup>

              <UFormGroup label="载客量">
                <UInput v-model.number="form.passengerCapacity" type="number" placeholder="人" />
              </UFormGroup>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="船员人数">
                <UInput v-model.number="form.crewCount" type="number" placeholder="人" />
              </UFormGroup>

              <UFormGroup label="甲板层数">
                <UInput v-model.number="form.deckCount" type="number" placeholder="层" />
              </UFormGroup>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="建造年份">
                <UInput v-model.number="form.builtYear" type="number" placeholder="年" />
              </UFormGroup>

              <UFormGroup label="翻新年份">
                <UInput v-model.number="form.renovatedYear" type="number" placeholder="年" />
              </UFormGroup>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="长度（米）">
                <UInput v-model.number="form.lengthMeters" type="number" step="0.1" placeholder="米" />
              </UFormGroup>

              <UFormGroup label="宽度（米）">
                <UInput v-model.number="form.widthMeters" type="number" step="0.1" placeholder="米" />
              </UFormGroup>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Images -->
      <UCard>
        <template #header>
          <h3 class="text-lg font-medium text-gray-900">邮轮图片</h3>
        </template>
        <ImageUpload
          v-model="form.coverImages"
          :multiple="true"
          :max-files="10"
          @upload="handleImageUpload"
        />
      </UCard>

      <!-- Submit Buttons -->
      <div class="flex justify-end gap-4">
        <UButton color="gray" variant="ghost" @click="$router.back()">
          取消
        </UButton>
        <UButton color="primary" type="submit" :loading="saving">
          {{ isEdit ? '保存修改' : '创建邮轮' }}
        </UButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
const { $api } = useNuxtApp()
const route = useRoute()
const router = useRouter()

const isEdit = computed(() => route.params.id !== 'new')

const form = reactive({
  companyId: '',
  nameCn: '',
  nameEn: '',
  code: '',
  grossTonnage: 0,
  passengerCapacity: 0,
  crewCount: 0,
  deckCount: 0,
  builtYear: 0,
  renovatedYear: 0,
  lengthMeters: 0,
  widthMeters: 0,
  status: 'active',
  coverImages: [] as string[]
})

const companyOptions = ref([
  { label: '皇家加勒比', value: '1' },
  { label: 'MSC地中海', value: '2' },
  { label: '歌诗达', value: '3' }
])

const statusOptions = [
  { label: '上架', value: 'active' },
  { label: '下架', value: 'inactive' },
  { label: '维护', value: 'maintenance' }
]

const saving = ref(false)

const fetchCruise = async () => {
  if (!isEdit.value) return
  
  try {
    const response = await $api(`/admin/cruises/${route.params.id}`)
    Object.assign(form, {
      companyId: response.companyId,
      nameCn: response.nameCn,
      nameEn: response.nameEn,
      code: response.code,
      grossTonnage: response.grossTonnage,
      passengerCapacity: response.passengerCapacity,
      crewCount: response.crewCount,
      deckCount: response.deckCount,
      builtYear: response.builtYear,
      renovatedYear: response.renovatedYear,
      lengthMeters: response.lengthMeters,
      widthMeters: response.widthMeters,
      status: response.status,
      coverImages: response.coverImages || []
    })
  } catch (error) {
    console.error('Failed to fetch cruise:', error)
  }
}

const handleImageUpload = (urls: string[]) => {
  form.coverImages = [...form.coverImages, ...urls]
}

const handleSubmit = async () => {
  try {
    saving.value = true
    
    const payload = {
      company_id: form.companyId,
      name_cn: form.nameCn,
      name_en: form.nameEn,
      code: form.code,
      gross_tonnage: form.grossTonnage,
      passenger_capacity: form.passengerCapacity,
      crew_count: form.crewCount,
      deck_count: form.deckCount,
      built_year: form.builtYear,
      renovated_year: form.renovatedYear,
      length_meters: form.lengthMeters,
      width_meters: form.widthMeters,
      status: form.status,
      cover_images: form.coverImages
    }
    
    if (isEdit.value) {
      await $api(`/admin/cruises/${route.params.id}`, {
        method: 'PUT',
        body: payload
      })
    } else {
      await $api('/admin/cruises', {
        method: 'POST',
        body: payload
      })
    }
    
    router.push('/cruises')
  } catch (error) {
    console.error('Failed to save cruise:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchCruise()
})
</script>
