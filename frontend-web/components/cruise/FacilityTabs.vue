<template>
  <div class="space-y-6">
    <!-- Category Tabs -->
    <div class="border-b border-gray-200">
      <nav class="flex space-x-8" aria-label="Tabs">
        <button
          v-for="category in categories"
          :key="category.id"
          @click="activeCategory = category.id"
          class="py-4 px-1 inline-flex items-center border-b-2 font-medium text-sm transition-colors"
          :class="activeCategory === category.id
            ? 'border-blue-600 text-blue-600'
            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
          "
        >
          <UIcon
            v-if="category.icon"
            :name="category.icon"
            class="w-5 h-5 mr-2"
          />
          {{ category.name }}
          <span
            class="ml-2 py-0.5 px-2 rounded-full text-xs"
            :class="activeCategory === category.id
              ? 'bg-blue-100 text-blue-600'
              : 'bg-gray-100 text-gray-600'
            "
          >
            {{ getFacilitiesByCategory(category.id).length }}
          </span>
        </button>
      </nav>
    </div>

    <!-- Facilities Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="facility in currentFacilities"
        :key="facility.id"
        class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
      >
        <div class="flex items-start space-x-4">
          <div class="flex-shrink-0">
            <div class="w-12 h-12 rounded-lg bg-blue-50 flex items-center justify-center">
              <UIcon name="i-heroicons-sparkles" class="w-6 h-6 text-blue-600" />
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between">
              <h4 class="text-base font-semibold text-gray-900 truncate">
                {{ facility.name }}
              </h4>
              <UBadge
                v-if="!facility.isFree"
                color="amber"
                variant="soft"
                size="sm"
              >
                收费
              </UBadge>
            </div>
            
            <p v-if="facility.deckNumber" class="text-sm text-gray-500 mt-1">
              <UIcon name="i-heroicons-map-pin" class="w-4 h-4 inline mr-1" />
              {{ facility.deckNumber }} 层甲板
            </p>
            
            <p v-if="facility.openTime" class="text-sm text-gray-500 mt-1">
              <UIcon name="i-heroicons-clock" class="w-4 h-4 inline mr-1" />
              {{ facility.openTime }}
            </p>
            
            <p v-if="facility.price && !facility.isFree" class="text-sm text-amber-600 mt-1">
              <UIcon name="i-heroicons-currency-yen" class="w-4 h-4 inline mr-1" />
              {{ facility.price }} 元
            </p>
            
            <p v-if="facility.description" class="text-sm text-gray-600 mt-2 line-clamp-2">
              {{ facility.description }}
            </p>

            <!-- Suitable Tags -->
            <div v-if="facility.suitableTags && facility.suitableTags.length > 0" class="mt-3 flex flex-wrap gap-1">
              <UBadge
                v-for="tag in facility.suitableTags"
                :key="tag"
                color="blue"
                variant="soft"
                size="xs"
              >
                {{ tag }}
              </UBadge>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="currentFacilities.length === 0" class="text-center py-8">
      <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-300 mx-auto mb-3" />
      <p class="text-gray-500">该分类下暂无设施</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  facilities: any[]
}>()

const activeCategory = ref<string>('')

// Group facilities by category
const categories = computed(() => {
  const categoryMap = new Map()
  
  props.facilities.forEach(facility => {
    const category = facility.category || { id: 'uncategorized', name: '其他' }
    if (!categoryMap.has(category.id)) {
      categoryMap.set(category.id, {
        id: category.id,
        name: category.name,
        icon: category.icon,
        sortWeight: category.sortWeight || 0
      })
    }
  })

  // Add "All" category if we have uncategorized items
  if (categoryMap.has('uncategorized')) {
    categoryMap.set('uncategorized', {
      id: 'uncategorized',
      name: '其他设施',
      icon: 'i-heroicons-ellipsis-horizontal',
      sortWeight: 999
    })
  }

  // Sort by sortWeight
  return Array.from(categoryMap.values()).sort((a, b) => a.sortWeight - b.sortWeight)
})

// Set default active category
watch(categories, (newCategories) => {
  if (newCategories.length > 0 && !activeCategory.value) {
    activeCategory.value = newCategories[0].id
  }
}, { immediate: true })

const getFacilitiesByCategory = (categoryId: string) => {
  return props.facilities.filter(facility => {
    const facilityCategoryId = facility.category?.id || 'uncategorized'
    return facilityCategoryId === categoryId
  })
}

const currentFacilities = computed(() => {
  return getFacilitiesByCategory(activeCategory.value)
})
</script>
