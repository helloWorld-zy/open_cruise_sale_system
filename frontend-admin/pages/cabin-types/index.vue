<template>
  <div>
    <!-- Page Header -->
    <div class="md:flex md:items-center md:justify-between mb-8">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900">
          舱房类型管理
        </h2>
      </div>
      <div class="mt-4 flex md:ml-4 md:mt-0">
        <UButton
          color="primary"
          icon="i-heroicons-plus"
          @click="openCreateModal"
        >
          新增舱房类型
        </UButton>
      </div>
    </div>

    <!-- Cruise Selector -->
    <div class="bg-white shadow rounded-lg p-4 mb-6">
      <UFormGroup label="选择邮轮">
        <USelect
          v-model="selectedCruiseId"
          :options="cruiseOptions"
          placeholder="选择邮轮查看舱房类型"
          @change="fetchCabinTypes"
        />
      </UFormGroup>
    </div>

    <!-- Cabin Type List -->
    <div v-if="selectedCruiseId" class="bg-white shadow rounded-lg">
      <UTable
        :rows="cabinTypes"
        :columns="columns"
        :loading="loading"
      >
        <template #name-data="{ row }">
          <div class="flex items-center">
            <img
              v-if="row.images && row.images.length > 0"
              :src="row.images[0]"
              class="h-10 w-10 rounded-lg object-cover mr-3"
            />
            <div>
              <div class="font-medium text-gray-900">{{ row.name }}</div>
              <div class="text-sm text-gray-500">{{ row.code }}</div>
            </div>
          </div>
        </template>

        <template #specs-data="{ row }">
          <div class="text-sm text-gray-500">
            <div>{{ row.minAreaSqm }}-{{ row.maxAreaSqm }}㎡</div>
            <div>可住 {{ row.standardGuests }} 人</div>
          </div>
        </template>

        <template #status-data="{ row }">
          <UBadge
            :color="row.status === 'active' ? 'green' : 'gray'"
          >
            {{ row.status === 'active' ? '启用' : '禁用' }}
          </UBadge>
        </template>

        <template #actions-data="{ row }">
          <UButton
            color="gray"
            variant="ghost"
            icon="i-heroicons-pencil-square"
            @click="openEditModal(row)"
          />
          <UButton
            color="gray"
            variant="ghost"
            icon="i-heroicons-trash"
            @click="confirmDelete(row)"
          />
        </template>
      </UTable>
    </div>

    <div v-else class="text-center py-12 bg-white shadow rounded-lg">
      <UIcon name="i-heroicons-cube" class="mx-auto h-12 w-12 text-gray-400" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">请先选择邮轮</h3>
      <p class="mt-1 text-sm text-gray-500">选择邮轮后可查看和管理该邮轮的舱房类型</p>
    </div>

    <!-- Create/Edit Modal -->
    <UModal v-model="modalOpen">
      <UCard>
        <template #header>
          <h3 class="text-lg font-medium text-gray-900">
            {{ isEdit ? '编辑舱房类型' : '新增舱房类型' }}
          </h3>
        </template>
        
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <UFormGroup label="名称" required>
            <UInput v-model="form.name" placeholder="输入舱房类型名称" />
          </UFormGroup>

          <UFormGroup label="编号" required>
            <UInput v-model="form.code" placeholder="输入唯一编号" />
          </UFormGroup>

          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="最小面积 (㎡)">
              <UInput v-model.number="form.minAreaSqm" type="number" step="0.1" />
            </UFormGroup>

            <UFormGroup label="最大面积 (㎡)">
              <UInput v-model.number="form.maxAreaSqm" type="number" step="0.1" />
            </UFormGroup>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="标准入住人数">
              <UInput v-model.number="form.standardGuests" type="number" />
            </UFormGroup>

            <UFormGroup label="最大入住人数">
              <UInput v-model.number="form.maxGuests" type="number" />
            </UFormGroup>
          </div>

          <UFormGroup label="床型">
            <UInput v-model="form.bedTypes" placeholder="如：双床/大床" />
          </UFormGroup>

          <UFormGroup label="状态">
            <USelect
              v-model="form.status"
              :options="[
                { label: '启用', value: 'active' },
                { label: '禁用', value: 'inactive' }
              ]"
            />
          </UFormGroup>

          <UFormGroup label="描述">
            <UTextarea v-model="form.description" rows="3" />
          </UFormGroup>
        </form>

        <template #footer>
          <div class="flex justify-end gap-3">
            <UButton color="gray" variant="ghost" @click="modalOpen = false">取消</UButton>
            <UButton color="primary" :loading="saving" @click="handleSubmit">
              {{ isEdit ? '保存' : '创建' }}
            </UButton>
          </div>
        </template>
      </UCard>
    </UModal>

    <!-- Delete Confirmation Modal -->
    <UModal v-model="deleteModalOpen">
      <UCard>
        <template #header>
          <h3 class="text-lg font-medium text-gray-900">确认删除</h3>
        </template>
        <p class="text-gray-500">确定要删除舱房类型 "{{ selectedCabinType?.name }}" 吗？</p>
        <template #footer>
          <div class="flex justify-end gap-3">
            <UButton color="gray" variant="ghost" @click="deleteModalOpen = false">取消</UButton>
            <UButton color="red" :loading="deleting" @click="handleDelete">删除</UButton>
          </div>
        </template>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
const { $api } = useNuxtApp()
const route = useRoute()

definePageMeta({
  layout: 'admin'
})

const columns = [
  { key: 'name', label: '舱房类型' },
  { key: 'specs', label: '规格' },
  { key: 'status', label: '状态' },
  { key: 'actions', label: '操作' }
]

const cruiseOptions = ref([])
const selectedCruiseId = ref(route.query.cruise_id as string || '')
const cabinTypes = ref([])
const loading = ref(false)

const modalOpen = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const form = reactive({
  id: '',
  name: '',
  code: '',
  minAreaSqm: 0,
  maxAreaSqm: 0,
  standardGuests: 2,
  maxGuests: 4,
  bedTypes: '',
  description: '',
  status: 'active'
})

const deleteModalOpen = ref(false)
const selectedCabinType = ref<any>(null)
const deleting = ref(false)

const fetchCruises = async () => {
  try {
    const response = await $api('/cruises', {
      query: { status: 'active', page_size: 100 }
    })
    cruiseOptions.value = (response.data || []).map((c: any) => ({
      label: c.nameCn,
      value: c.id
    }))
  } catch (error) {
    console.error('Failed to fetch cruises:', error)
  }
}

const fetchCabinTypes = async () => {
  if (!selectedCruiseId.value) return
  
  try {
    loading.value = true
    const response = await $api('/admin/cabin-types', {
      query: { cruise_id: selectedCruiseId.value }
    })
    cabinTypes.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch cabin types:', error)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEdit.value = false
  Object.assign(form, {
    id: '',
    name: '',
    code: '',
    minAreaSqm: 0,
    maxAreaSqm: 0,
    standardGuests: 2,
    maxGuests: 4,
    bedTypes: '',
    description: '',
    status: 'active'
  })
  modalOpen.value = true
}

const openEditModal = (cabinType: any) => {
  isEdit.value = true
  Object.assign(form, cabinType)
  modalOpen.value = true
}

const handleSubmit = async () => {
  try {
    saving.value = true
    const payload = {
      cruise_id: selectedCruiseId.value,
      name: form.name,
      code: form.code,
      min_area_sqm: form.minAreaSqm,
      max_area_sqm: form.maxAreaSqm,
      standard_guests: form.standardGuests,
      max_guests: form.maxGuests,
      bed_types: form.bedTypes,
      description: form.description,
      status: form.status
    }
    
    if (isEdit.value) {
      await $api(`/admin/cabin-types/${form.id}`, {
        method: 'PUT',
        body: payload
      })
    } else {
      await $api('/admin/cabin-types', {
        method: 'POST',
        body: payload
      })
    }
    
    modalOpen.value = false
    fetchCabinTypes()
  } catch (error) {
    console.error('Failed to save cabin type:', error)
  } finally {
    saving.value = false
  }
}

const confirmDelete = (cabinType: any) => {
  selectedCabinType.value = cabinType
  deleteModalOpen.value = true
}

const handleDelete = async () => {
  if (!selectedCabinType.value) return
  
  try {
    deleting.value = true
    await $api(`/admin/cabin-types/${selectedCabinType.value.id}`, {
      method: 'DELETE'
    })
    deleteModalOpen.value = false
    fetchCabinTypes()
  } catch (error) {
    console.error('Failed to delete cabin type:', error)
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchCruises()
  if (selectedCruiseId.value) {
    fetchCabinTypes()
  }
})
</script>
