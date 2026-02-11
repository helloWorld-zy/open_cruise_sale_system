<template>
  <div class="space-y-4">
    <div
      v-for="cabin in cabinTypes"
      :key="cabin.id"
      class="border border-gray-200 rounded-lg overflow-hidden"
    >
      <button
        @click="toggleExpand(cabin.id)"
        class="w-full px-6 py-4 flex items-center justify-between bg-gray-50 hover:bg-gray-100 transition-colors"
      >
        <div class="flex items-center space-x-4">
          <img
            v-if="cabin.images && cabin.images.length > 0"
            :src="cabin.images[0]"
            :alt="cabin.name"
            class="w-16 h-16 rounded-lg object-cover"
          />
          <div v-else class="w-16 h-16 rounded-lg bg-gray-200 flex items-center justify-center">
            <UIcon name="i-heroicons-photo" class="w-6 h-6 text-gray-400" />
          </div>
          
          <div class="text-left">
            <h4 class="font-semibold text-gray-900">{{ cabin.name }}</h4>
            <p class="text-sm text-gray-500">
              {{ cabin.minAreaSqm }}-{{ cabin.maxAreaSqm }}㎡ · 可住{{ cabin.standardGuests }}人
            </p>
          </div>
        </div>
        
        <div class="flex items-center space-x-4">
          <span class="text-lg font-bold text-blue-600">¥{{ cabin.price || '咨询' }}</span>
          <UIcon
            :name="expanded[cabin.id] ? 'i-heroicons-chevron-up' : 'i-heroicons-chevron-down'"
            class="w-5 h-5 text-gray-400"
          />
        </div>
      </button>
      
      <div v-if="expanded[cabin.id]" class="px-6 py-4 border-t border-gray-200">
        <div class="prose prose-sm max-w-none text-gray-600 mb-4" v-html="cabin.description" />
        
        <div v-if="cabin.amenities && cabin.amenities.length > 0" class="mb-4">
          <h5 class="font-medium text-gray-900 mb-2">设施配置</h5>
          <div class="flex flex-wrap gap-2">
            <UBadge
              v-for="amenity in cabin.amenities"
              :key="amenity"
              color="gray"
              variant="soft"
            >
              {{ amenity }}
            </UBadge>
          </div>
        </div>
        
        <div v-if="cabin.images && cabin.images.length > 1" class="grid grid-cols-4 gap-2">
          <img
            v-for="(image, index) in cabin.images.slice(1)"
            :key="index"
            :src="image"
            :alt="`${cabin.name} ${index + 2}`"
            class="w-full aspect-square object-cover rounded-lg"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  cabinTypes: any[]
}>()

const expanded = ref<Record<string, boolean>>({})

const toggleExpand = (id: string) => {
  expanded.value[id] = !expanded.value[id]
}
</script>
