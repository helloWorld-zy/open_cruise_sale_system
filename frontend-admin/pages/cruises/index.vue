<template>
  <div>
    <!-- Page Header -->
    <div class="md:flex md:items-center md:justify-between mb-8">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
          邮轮管理
        </h2>
      </div>
      <div class="mt-4 flex md:ml-4 md:mt-0">
        <UButton
          color="primary"
          icon="i-heroicons-plus"
          @click="navigateTo('/cruises/new')"
        >
          新增邮轮
        </UButton>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white shadow rounded-lg p-4 mb-6">
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
        <UButton color="primary" @click="fetchCruises">
          搜索
        </UButton>
        <UButton color="gray" variant="ghost" @click="resetFilters">
          重置
        </UButton>
      </div>
    </div>

    <!-- Cruise List -->
    <div class="bg-white shadow rounded-lg">
      <UTable
        :rows="cruises"
        :columns="columns"
        :loading="loading"
      >
        <template #nameCn-data="{ row }">
          <div class="flex items-center">
            <img
              v-if="row.coverImages && row.coverImages.length > 0"
              :src="row.coverImages[0]"
              class="h-10 w-10 rounded-lg object-cover mr-3"
            />
            <div>
              <div class="font-medium text-gray-900">{{ row.nameCn }}</div>
              <div class="text-sm text-gray-500">{{ row.code }}</div>
            </div>
          </div>
        </template>

        <template #specs-data="{ row }">
          <div class="text-sm text-gray-500">
            <div>{{ row.grossTonnage?.toLocaleString() }} 吨</div>
            <div>{{ row.passengerCapacity?.toLocaleString() }} 人</div>
          </div>
        </template>

        <template #status-data="{ row }">
          <UBadge
            :color="row.status === 'active' ? 'green' : row.status === 'maintenance' ? 'yellow' : 'gray'"
          >
            {{ statusText[row.status] || row.status }}
          </UBadge>
        </template>

        <template #actions-data="{ row }">
          <UDropdown
            :items="[
              { label: '编辑', icon: 'i-heroicons-pencil-square', click: () => navigateTo(`/cruises/${row.id}`) },
              { label: '查看舱房', icon: 'i-heroicons-square-2-stack', click: () => navigateTo(`/cabin-types?cruise_id=${row.id}`) },
              { label: '查看设施', icon: 'i-heroicons-building-storefront', click: () => navigateTo(`/facilities?cruise_id=${row.id}`) },
              { label: '删除', icon: 'i-heroicons-trash', click: () => confirmDelete(row) }
            ]"
          >
            <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal" />
          </UDropdown>
        </template>
      </UTable>

      <!-- Pagination -->
      <div class="px-4 py-3 border-t border-gray-200">
        <UPagination
          v-model="page"
          :total="total"
          :page-count="pageSize"
          @update:model-value="fetchCruises"
        />
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <UModal v-model="deleteModalOpen">
      <UCard>
        <template #header>
          <h3 class="text-lg font-medium text-gray-900">确认删除</h3>
        </template>
        <p class="text-gray-500">确定要删除邮轮 "{{ selectedCruise?.nameCn }}" 吗？此操作不可恢复。</p>
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

definePageMeta({
  layout: 'admin'
})

const columns = [
  { key: 'nameCn', label: '邮轮' },
  { key: 'specs', label: '规格' },
  { key: 'status', label: '状态' },
  { key: 'actions', label: '操作' }
]

const statusOptions = [
  { label: '全部', value: '' },
  { label: '上架', value: 'active' },
  { label: '下架', value: 'inactive' },
  { label: '维护', value: 'maintenance' }
]

const statusText: Record<string, string> = {
  active: '上架',
  inactive: '下架',
  maintenance: '维护'
}

const cruises = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const filters = reactive({
  keyword: '',
  status: ''
})

const deleteModalOpen = ref(false)
const selectedCruise = ref<any>(null)
const deleting = ref(false)

const fetchCruises = async () => {
  try {
    loading.value = true
    const response = await $api('/admin/cruises', {
      query: {
        page: page.value,
        page_size: pageSize.value,
        keyword: filters.keyword || undefined,
        status: filters.status || undefined
      }
    })
    cruises.value = response.data || []
    total.value = response.pagination?.total || 0
  } catch (error) {
    console.error('Failed to fetch cruises:', error)
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.keyword = ''
  filters.status = ''
  page.value = 1
  fetchCruises()
}

const confirmDelete = (cruise: any) => {
  selectedCruise.value = cruise
  deleteModalOpen.value = true
}

const handleDelete = async () => {
  if (!selectedCruise.value) return
  
  try {
    deleting.value = true
    await $api(`/admin/cruises/${selectedCruise.value.id}`, {
      method: 'DELETE'
    })
    deleteModalOpen.value = false
    fetchCruises()
  } catch (error) {
    console.error('Failed to delete cruise:', error)
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchCruises()
})
</script>
