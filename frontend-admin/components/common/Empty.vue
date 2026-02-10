<template>
  <div class="flex items-center justify-center h-64">
    <div class="text-center">
      <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-400 mx-auto" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">{{ title }}</h3>
      <p class="mt-1 text-sm text-gray-500">{{ description }}</p>
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title?: string
  description?: string
}>(), {
  title: 'No items',
  description: 'Get started by creating a new item.'
})
</script>
