<template>
  <div>
    <!-- Page Header -->
    <div class="md:flex md:items-center md:justify-between mb-8">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900">
          设施管理
        </h2>
      </div>
      <div class="mt-4 flex gap-3 md:mt-0">
        <UButton
          color="gray"
          variant="soft"
          icon="i-heroicons-folder-plus"
          @click="openCategoryModal"
        >
          管理分类
        </UButton>
        <UButton
          color="primary"
          icon="i-heroicons-plus"
          @click="openCreateModal"
        >
          新增设施
        </UButton>
      </div>
    </div>

    <!-- Cruise Selector -->
    <div class="bg-white shadow rounded-lg p-4 mb-6">
      <UFormGroup label="选择邮轮">
        <USelect
          v-model="selectedCruiseId"
          :options="cruiseOptions"
          placeholder="选择邮轮查看设施"
          @change="fetchFacilities"
        />
      </UFormGroup>
    </div>

    <!-- Facility List -->
    <div v-if="selectedCruiseId" class="bg-white shadow rounded-lg">
      <UTable
        :rows="facilities"
        :columns="columns"
        :loading="loading"
      >
        <template #name-data="{ row }">
          <div class="flex items-center">
            <div class="h-10 w-10 rounded-lg bg-blue-100 flex items-center justify-center mr-3">
              <UIcon name="i-heroicons-sparkles" class="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <div class="font-medium text-gray-900">{{ row.name }}</div>
              <div v-if="row.category" class="text-sm text-gray-500">{{ row.category.name }}</div>
            </div>
          </div>
        </template>

        <template #location-data="{ row }">
          <div class="text-sm text-gray-500">
            <div v-if="row.deckNumber">{{ row.deckNumber }} 层甲板</div>
            <div v-if="row.openTime">营业时间: {{ row.openTime }}</div>
          </div>
        </template>

        <template #pricing-data="{ row }">
          <UBadge
            :color="row.isFree ? 'green' : 'amber'"
          >
            {{ row.isFree ? '免费' : `¥${row.price}` }}
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
      <UIcon name="i-heroicons-building-storefront" class="mx-auto h-12 w-12 text-gray-400" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">请先选择邮轮</h3>
      <p class="mt-1 text-sm text-gray-500">选择邮轮后可查看和管理该邮轮的设施</p>
    </div>

    <!-- Facility Modal -->
    <UModal v-model="modalOpen">
      <UCard>
        <template #header>
          <h3 class="text-lg font-medium text-gray-900">
            {{ isEdit ? '编辑设施' : '新增设施' }}
          </h3>
        </template>
        
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <UFormGroup label="名称" required>
            <UInput v-model="form.name" placeholder="输入设施名称" />
          </UFormGroup>

          <UFormGroup label="分类">
            <USelect
              v-model="form.categoryId"
              :options="categoryOptions"
              placeholder="选择分类"
            />
          </UFormGroup>

          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="甲板层">
              <UInput v-model.number="form.deckNumber" type="number" placeholder="层" />
            </UFormGroup>

            <UFormGroup label="营业时间">
              <UInput v-model="form.openTime" placeholder="如: 09:00-22:00" />
            </UFormGroup>
          </div>

          <UFormGroup>
            <UCheckbox v-model="form.isFree" label="免费设施" />
          </UFormGroup>

          <UFormGroup v-if="!form.isFree" label="价格 (元)">
            <UInput v-model.number="form.price" type="number" step="0.01" />
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

    <!-- Category Management Modal -->
    <UModal v-model="categoryModalOpen" size="lg">
      <UCard>
        <template #header>
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-medium text-gray-900">设施分类管理</h3>
            <UButton color="primary" size="sm" icon="i-heroicons-plus" @click="openCreateCategory">
              新增分类
            </UButton>
          </div>
        </template>

        <UTable
          :rows="categories"
          :columns="[{ key: 'name', label: '分类名称' }, { key: 'icon', label: '图标' }, { key: 'actions', label: '操作' }]"
        >
          <template #actions-data="{ row }">
            <UButton color="gray" variant="ghost" icon="i-heroicons-pencil" @click="editCategory(row)" />
            <UButton color="gray" variant="ghost" icon="i-heroicons-trash" @click="deleteCategory(row)" />
          </template>
        </UTable>
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
  { key: 'name', label: '设施' },
  { key: 'location', label: '位置/时间' },
  { key: 'pricing', label: '费用' },
  { key: 'actions', label: '操作' }
]

const cruiseOptions = ref([])
const selectedCruiseId = ref(route.query.cruise_id as string || '')
const facilities = ref([])
const categories = ref([])
const categoryOptions = ref([])
const loading = ref(false)

const modalOpen = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const form = reactive({
  id: '',
  name: '',
  categoryId: '',
  deckNumber: 0,
  openTime: '',
  isFree: true,
  price: 0,
  description: ''
})

const categoryModalOpen = ref(false)

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

const fetchCategories = async () => {
  if (!selectedCruiseId.value) return
  
  try {
    const response = await $api('/admin/facility-categories', {
      query: { cruise_id: selectedCruiseId.value }
    })
    categories.value = response || []
    categoryOptions.value = categories.value.map((c: any) => ({
      label: c.name,
      value: c.id
    }))
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

const fetchFacilities = async () => {
  if (!selectedCruiseId.value) return
  
  try {
    loading.value = true
    const response = await $api('/admin/facilities', {
      query: { cruise_id: selectedCruiseId.value }
    })
    facilities.value = response.data || []
    fetchCategories()
  } catch (error) {
    console.error('Failed to fetch facilities:', error)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEdit.value = false
  Object.assign(form, {
    id: '',
    name: '',
    categoryId: '',
    deckNumber: 0,
    openTime: '',
    isFree: true,
    price: 0,
    description: ''
  })
  modalOpen.value = true
}

const openEditModal = (facility: any) => {
  isEdit.value = true
  Object.assign(form, facility)
  modalOpen.value = true
}

const handleSubmit = async () => {
  try {
    saving.value = true
    const payload = {
      cruise_id: selectedCruiseId.value,
      category_id: form.categoryId,
      name: form.name,
      deck_number: form.deckNumber,
      open_time: form.openTime,
      is_free: form.isFree,
      price: form.price,
      description: form.description
    }
    
    if (isEdit.value) {
      await $api(`/admin/facilities/${form.id}`, {
        method: 'PUT',
        body: payload
      })
    } else {
      await $api('/admin/facilities', {
        method: 'POST',
        body: payload
      })
    }
    
    modalOpen.value = false
    fetchFacilities()
  } catch (error) {
    console.error('Failed to save facility:', error)
  } finally {
    saving.value = false
  }
}

const openCategoryModal = () => {
  fetchCategories()
  categoryModalOpen.value = true
}

const openCreateCategory = () => {
  // Implement category creation
}

const editCategory = (category: any) => {
  // Implement category edit
}

const deleteCategory = async (category: any) => {
  if (!confirm(`确定要删除分类 "${category.name}" 吗？`)) return
  
  try {
    await $api(`/admin/facility-categories/${category.id}`, {
      method: 'DELETE'
    })
    fetchCategories()
  } catch (error) {
    console.error('Failed to delete category:', error)
  }
}

const confirmDelete = (facility: any) => {
  if (!confirm(`确定要删除设施 "${facility.name}" 吗？`)) return
  deleteFacility(facility.id)
}

const deleteFacility = async (id: string) => {
  try {
    await $api(`/admin/facilities/${id}`, {
      method: 'DELETE'
    })
    fetchFacilities()
  } catch (error) {
    console.error('Failed to delete facility:', error)
  }
}

onMounted(() => {
  fetchCruises()
  if (selectedCruiseId.value) {
    fetchFacilities()
  }
})
</script>
