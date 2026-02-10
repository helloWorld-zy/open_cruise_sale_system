<template>
  <div class="flex items-center justify-center h-64">
    <div class="text-center">
      <UIcon name="i-heroicons-arrow-path" class="w-8 h-8 animate-spin text-primary-600 mx-auto" />
      <p class="mt-2 text-sm text-gray-600">{{ message }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  message?: string
}>(), {
  message: 'Loading...'
})
</script>
