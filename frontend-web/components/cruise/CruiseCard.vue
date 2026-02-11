<template>
  <UCard class="overflow-hidden hover:shadow-lg transition-shadow">
    <!-- Cover Image -->
    <div class="relative h-48 -mx-4 -mt-4 mb-4">
      <img
        v-if="cruise.coverImages && cruise.coverImages.length > 0"
        :src="cruise.coverImages[0]"
        :alt="cruise.nameCn"
        class="w-full h-full object-cover"
      />
      <div v-else class="w-full h-full bg-gray-200 flex items-center justify-center">
        <UIcon name="i-heroicons-photo" class="w-12 h-12 text-gray-400" />
      </div>
      
      <!-- Status Badge -->
      <UBadge
        :color="cruise.status === 'active' ? 'green' : 'gray'"
        class="absolute top-4 right-4"
      >
        {{ statusText }}
      </UBadge>
    </div>

    <!-- Content -->
    <div class="space-y-3">
      <h3 class="text-xl font-semibold text-gray-900 line-clamp-1">
        {{ cruise.nameCn }}
      </h3>
      
      <p v-if="cruise.nameEn" class="text-sm text-gray-500">
        {{ cruise.nameEn }}
      </p>

      <!-- Specs -->
      <div class="flex flex-wrap gap-4 text-sm text-gray-600">
        <span v-if="cruise.grossTonnage" class="flex items-center">
          <UIcon name="i-heroicons-scale" class="w-4 h-4 mr-1" />
          {{ cruise.grossTonnage.toLocaleString() }} 吨
        </span>
        <span v-if="cruise.passengerCapacity" class="flex items-center">
          <UIcon name="i-heroicons-users" class="w-4 h-4 mr-1" />
          {{ cruise.passengerCapacity.toLocaleString() }} 人
        </span>
        <span v-if="cruise.deckCount" class="flex items-center">
          <UIcon name="i-heroicons-squares-2x2" class="w-4 h-4 mr-1" />
          {{ cruise.deckCount }} 层甲板
        </span>
      </div>

      <!-- Actions -->
      <div class="pt-4 flex gap-3">
        <UButton
          :to="`/cruises/${cruise.id}`"
          color="primary"
          variant="solid"
          block
        >
          查看详情
        </UButton>
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import type { Cruise } from '@cruisebooking/types'

const props = defineProps<{
  cruise: Cruise
}>()

const statusText = computed(() => {
  const statusMap: Record<string, string> = {
    active: '热售中',
    inactive: '已下架',
    maintenance: '维护中'
  }
  return statusMap[props.cruise.status] || props.cruise.status
})
</script>
