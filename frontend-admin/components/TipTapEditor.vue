<template>
  <div class="border border-gray-300 rounded-lg overflow-hidden">
    <!-- Toolbar -->
    <div class="border-b border-gray-300 bg-gray-50 p-2 flex flex-wrap gap-1">
      <UButton
        size="xs"
        :color="editor?.isActive('bold') ? 'primary' : 'gray'"
        variant="ghost"
        icon="i-heroicons-bold"
        @click="editor?.chain().focus().toggleBold().run()"
      />
      <UButton
        size="xs"
        :color="editor?.isActive('italic') ? 'primary' : 'gray'"
        variant="ghost"
        icon="i-heroicons-italic"
        @click="editor?.chain().focus().toggleItalic().run()"
      />
      <UButton
        size="xs"
        :color="editor?.isActive('underline') ? 'primary' : 'gray'"
        variant="ghost"
        icon="i-heroicons-underline"
        @click="editor?.chain().focus().toggleUnderline().run()"
      />
      <div class="w-px h-6 bg-gray-300 mx-1" />
      <UButton
        size="xs"
        :color="editor?.isActive('heading', { level: 1 }) ? 'primary' : 'gray'"
        variant="ghost"
        @click="editor?.chain().focus().toggleHeading({ level: 1 }).run()"
      >
        H1
      </UButton>
      <UButton
        size="xs"
        :color="editor?.isActive('heading', { level: 2 }) ? 'primary' : 'gray'"
        variant="ghost"
        @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()"
      >
        H2
      </UButton>
      <UButton
        size="xs"
        :color="editor?.isActive('heading', { level: 3 }) ? 'primary' : 'gray'"
        variant="ghost"
        @click="editor?.chain().focus().toggleHeading({ level: 3 }).run()"
      >
        H3
      </UButton>
      <div class="w-px h-6 bg-gray-300 mx-1" />
      <UButton
        size="xs"
        :color="editor?.isActive('bulletList') ? 'primary' : 'gray'"
        variant="ghost"
        icon="i-heroicons-list-bullet"
        @click="editor?.chain().focus().toggleBulletList().run()"
      />
      <UButton
        size="xs"
        :color="editor?.isActive('orderedList') ? 'primary' : 'gray'"
        variant="ghost"
        icon="i-heroicons-numbered-list"
        @click="editor?.chain().focus().toggleOrderedList().run()"
      />
      <div class="w-px h-6 bg-gray-300 mx-1" />
      <UButton
        size="xs"
        color="gray"
        variant="ghost"
        icon="i-heroicons-arrow-uturn-left"
        @click="editor?.chain().focus().undo().run()"
        :disabled="!editor?.can().undo()"
      />
      <UButton
        size="xs"
        color="gray"
        variant="ghost"
        icon="i-heroicons-arrow-uturn-right"
        @click="editor?.chain().focus().redo().run()"
        :disabled="!editor?.can().redo()"
      />
    </div>

    <!-- Editor Content -->
    <EditorContent
      :editor="editor"
      class="prose prose-sm max-w-none p-4 min-h-[200px] focus:outline-none"
    />
  </div>
</template>

<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import Heading from '@tiptap/extension-heading'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const editor = useEditor({
  extensions: [
    StarterKit,
    Underline,
    Heading.configure({
      levels: [1, 2, 3]
    })
  ],
  content: props.modelValue,
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  }
})

// Watch for external changes
watch(() => props.modelValue, (newValue) => {
  if (editor.value && editor.value.getHTML() !== newValue) {
    editor.value.commands.setContent(newValue)
  }
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style>
.ProseMirror:focus {
  outline: none;
}

.ProseMirror p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: #adb5bd;
  pointer-events: none;
  height: 0;
}
</style>
