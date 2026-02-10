<template>
  <div class="flex items-center justify-center h-64">
    <div class="text-center">
      <UIcon name="i-heroicons-exclamation-triangle" class="w-12 h-12 text-red-500 mx-auto" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">Error</h3>
      <p class="mt-1 text-sm text-gray-500">{{ message }}</p>
      <UButton
        v-if="showRetry"
        color="primary"
        variant="soft"
        class="mt-4"
        @click="$emit('retry')"
      >
        Try Again
      </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  message: string
  showRetry?: boolean
}>()

defineEmits<{
  retry: []
}>()
</script>
