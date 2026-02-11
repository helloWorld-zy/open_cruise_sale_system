<template>
  <div class="space-y-6">
    <!-- Main Image -->
    <div class="relative aspect-video rounded-lg overflow-hidden">
      <img
        :src="currentImage || '/placeholder-cruise.jpg'"
        :alt="title"
        class="w-full h-full object-cover"
      />
    </div>
    
    <!-- Thumbnails -->
    <div v-if="images.length > 1" class="flex gap-2 overflow-x-auto pb-2">
      <button
        v-for="(image, index) in images"
        :key="index"
        @click="currentImage = image"
        class="flex-shrink-0 w-20 h-20 rounded-lg overflow-hidden border-2 transition-colors"
        :class="currentImage === image ? 'border-blue-600' : 'border-transparent'"
      >
        <img :src="image" :alt="`${title} ${index + 1}`" class="w-full h-full object-cover" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  images: string[]
  title: string
}>()

const currentImage = ref(props.images[0] || '')

watch(() => props.images, (newImages) => {
  if (newImages.length > 0 && !newImages.includes(currentImage.value)) {
    currentImage.value = newImages[0]
  }
}, { immediate: true })
</script>
