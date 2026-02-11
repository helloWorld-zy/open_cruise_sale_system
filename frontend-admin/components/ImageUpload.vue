<template>
  <div class="space-y-4">
    <!-- Drop Zone -->
    <div
      class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-primary-500 transition-colors cursor-pointer"
      :class="{ 'border-primary-500 bg-primary-50': isDragging }"
      @dragenter.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @dragover.prevent
      @drop.prevent="handleDrop"
      @click="triggerFileInput"
    >
      <input
        ref="fileInput"
        type="file"
        :multiple="multiple"
        :accept="accept"
        class="hidden"
        @change="handleFileChange"
      />
      <UIcon name="i-heroicons-cloud-arrow-up" class="mx-auto h-12 w-12 text-gray-400" />
      <div class="mt-4">
        <p class="text-sm text-gray-600">
          <span class="font-medium text-primary-600">点击上传</span> 或拖拽文件到此处
        </p>
        <p class="text-xs text-gray-500 mt-1">
          支持 PNG, JPG, GIF, WEBP 格式，最大 10MB
        </p>
      </div>
    </div>

    <!-- Preview Grid -->
    <div v-if="modelValue.length > 0" class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div
        v-for="(image, index) in modelValue"
        :key="index"
        class="relative group aspect-square"
      >
        <img
          :src="image"
          class="w-full h-full object-cover rounded-lg"
          alt="Preview"
        />
        <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-50 transition-all rounded-lg flex items-center justify-center opacity-0 group-hover:opacity-100">
          <button
            type="button"
            class="p-2 bg-red-600 text-white rounded-full hover:bg-red-700"
            @click="removeImage(index)"
          >
            <UIcon name="i-heroicons-trash" class="w-5 h-5" />
          </button>
        </div>
        <div v-if="uploading[index]" class="absolute inset-0 bg-black bg-opacity-50 rounded-lg flex items-center justify-center">
          <UIcon name="i-heroicons-arrow-path" class="w-8 h-8 text-white animate-spin" />
        </div>
      </div>
    </div>

    <!-- Upload Progress -->
    <div v-if="uploadProgress > 0 && uploadProgress < 100" class="w-full bg-gray-200 rounded-full h-2.5">
      <div
        class="bg-primary-600 h-2.5 rounded-full transition-all"
        :style="{ width: uploadProgress + '%' }"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: string[]
  multiple?: boolean
  maxFiles?: number
  accept?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
  'upload': [urls: string[]]
}>()

const { $api } = useNuxtApp()

const fileInput = ref<HTMLInputElement>()
const isDragging = ref(false)
const uploading = ref<boolean[]>([])
const uploadProgress = ref(0)

const accept = props.accept || 'image/png,image/jpeg,image/gif,image/webp'

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    uploadFiles(Array.from(target.files))
  }
}

const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  if (event.dataTransfer?.files) {
    uploadFiles(Array.from(event.dataTransfer.files))
  }
}

const uploadFiles = async (files: File[]) => {
  // Validate max files
  if (props.maxFiles && modelValue.value.length + files.length > props.maxFiles) {
    alert(`最多只能上传 ${props.maxFiles} 张图片`)
    return
  }

  // Validate file types
  const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
  const invalidFiles = files.filter(file => !validTypes.includes(file.type))
  if (invalidFiles.length > 0) {
    alert('只支持 PNG, JPG, GIF, WEBP 格式的图片')
    return
  }

  // Validate file sizes (10MB max)
  const maxSize = 10 * 1024 * 1024
  const oversizedFiles = files.filter(file => file.size > maxSize)
  if (oversizedFiles.length > 0) {
    alert('图片大小不能超过 10MB')
    return
  }

  // Initialize uploading state
  const startIndex = modelValue.value.length
  files.forEach((_, index) => {
    uploading.value[startIndex + index] = true
  })

  try {
    const formData = new FormData()
    files.forEach(file => {
      formData.append('images', file)
    })

    const response = await $api('/admin/cruises/upload', {
      method: 'POST',
      body: formData
    })

    // Update model value with new URLs
    const newUrls = Array.isArray(response) ? response : [response]
    const updatedUrls = [...modelValue.value, ...newUrls]
    emit('update:modelValue', updatedUrls)
    emit('upload', newUrls)
  } catch (error) {
    console.error('Upload failed:', error)
    alert('上传失败，请重试')
  } finally {
    // Clear uploading state
    files.forEach((_, index) => {
      uploading.value[startIndex + index] = false
    })
    uploadProgress.value = 0
  }
}

const removeImage = (index: number) => {
  const updatedUrls = [...modelValue.value]
  updatedUrls.splice(index, 1)
  emit('update:modelValue', updatedUrls)
}
</script>
