<script setup>
import '@wangeditor/editor/dist/css/style.css'

import { onBeforeUnmount, ref, shallowRef, watch } from 'vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'

import { uploadContentImage } from '@/api/article'
import { useUiStore } from '@/stores/uiStore'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])

const uiStore = useUiStore()
const editorRef = shallowRef(null)
const htmlValue = ref(props.modelValue || '')
const imageUploading = ref(false)

const toolbarConfig = {
  toolbarKeys: [
    'headerSelect',
    '|',
    'bold',
    'italic',
    'underline',
    'through',
    'color',
    'bgColor',
    'insertLink',
    'uploadImage',
    '|',
    'bulletedList',
    'numberedList',
    '|',
    'group-justify',
    'blockquote',
    'divider',
    'codeBlock',
    'insertTable',
  ],
}

const editorConfig = {
  placeholder: '开始写正文内容，支持粘贴截图、拖拽图片和所见即所得排版。',
  autoFocus: true,
  scroll: true,
  MENU_CONF: {
    uploadImage: {
      maxFileSize: 5 * 1024 * 1024,
      maxNumberOfFiles: 10,
      allowedFileTypes: ['image/jpeg', 'image/png', 'image/webp', 'image/gif'],
      async customUpload(file, insertFn) {
        imageUploading.value = true

        try {
          const data = await uploadContentImage(file)
          const url = data?.content_image || data?.url || ''

          if (!url) {
            throw new Error('正文图片上传失败，请重试')
          }

          insertFn(url, file.name || 'article-image', url)
          uiStore.showToast('正文图片上传成功', 'success')
        } catch (error) {
          uiStore.showToast(error.message || '正文图片上传失败', 'error')
          throw error
        } finally {
          imageUploading.value = false
        }
      },
    },
  },
}

watch(
  () => props.modelValue,
  (value) => {
    if (value !== htmlValue.value) {
      htmlValue.value = value || ''
    }
  },
)

watch(htmlValue, (value) => {
  emit('update:modelValue', value || '')
})

function handleCreated(editor) {
  editorRef.value = Object.seal(editor)
}

function handleAlert(message, type) {
  uiStore.showToast(String(message || ''), type === 'error' ? 'error' : 'info')
}

function handleBodyClick(event) {
  const target = event.target
  if (!(target instanceof HTMLElement)) {
    return
  }

  const anchor = target.closest('a[href]')
  if (!(anchor instanceof HTMLAnchorElement)) {
    return
  }

  event.preventDefault()
  window.open(anchor.href, '_blank', 'noopener,noreferrer')
}

onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor) {
    editor.destroy()
  }
})
</script>

<template>
  <div class="rich-editor" :class="{ 'is-uploading': imageUploading }">
    <div class="rich-editor__toolbar">
      <Toolbar :editor="editorRef" :default-config="toolbarConfig" mode="default" />
    </div>

    <div class="rich-editor__body" @click="handleBodyClick">
      <Editor
        v-model="htmlValue"
        :default-config="editorConfig"
        mode="default"
        @on-created="handleCreated"
        @custom-alert="handleAlert"
      />
    </div>

    <div class="rich-editor__foot">
      <span>支持标题、颜色、表格、代码块、分割线、链接与图片。</span>
      <span>{{ imageUploading ? '图片上传中...' : '' }}</span>
    </div>
  </div>
</template>

<style scoped>
.rich-editor {
  border: 1px solid rgba(33, 41, 52, 0.1);
  border-radius: 14px;
  overflow: hidden;
  background: #fff;
  box-shadow: none;
}

.rich-editor.is-uploading {
  border-color: rgba(93, 163, 158, 0.32);
}

.rich-editor__toolbar {
  border-bottom: 1px solid rgba(33, 41, 52, 0.08);
  background: #fafbfc;
}

.rich-editor__body {
  min-height: 370px;
  background: #fff;
}

.rich-editor__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  border-top: 1px solid rgba(33, 41, 52, 0.06);
  color: var(--color-muted);
  font-size: 12px;
}

:deep(.w-e-bar) {
  padding: 8px 10px;
  background: transparent;
}

:deep(.w-e-toolbar) {
  gap: 2px 4px;
}

:deep(.w-e-bar-item button) {
  min-width: 32px;
  border-radius: 8px;
  transition:
    background-color 0.18s ease,
    color 0.18s ease;
}

:deep(.w-e-text-container) {
  background: transparent;
}

:deep(.w-e-text-container [data-slate-editor]) {
  min-height: 280px;
  padding: 18px 20px 22px;
  font-size: 15px;
  line-height: 1.8;
  color: #24323d;
}

:deep(.w-e-text-placeholder) {
  top: 18px;
  left: 20px;
  font-style: normal;
}

:deep(.w-e-text-container [data-slate-editor] p) {
  margin: 0 0 16px;
}

:deep(.w-e-text-container [data-slate-editor] a) {
  color: #1f6fff;
  cursor: pointer;
  text-decoration: none;
  transition:
    color 0.18s ease,
    opacity 0.18s ease;
}

:deep(.w-e-text-container [data-slate-editor] a:hover) {
  color: #0f56d8;
  opacity: 0.92;
}

:deep(.w-e-text-container [data-slate-editor] h1),
:deep(.w-e-text-container [data-slate-editor] h2),
:deep(.w-e-text-container [data-slate-editor] h3) {
  margin: 22px 0 12px;
  color: #1d2f37;
  line-height: 1.35;
}

:deep(.w-e-text-container [data-slate-editor] blockquote) {
  margin: 18px 0;
  border-left: 4px solid rgba(93, 163, 158, 0.72);
  background: rgba(93, 163, 158, 0.08);
  color: #4c616e;
}

:deep(.w-e-text-container [data-slate-editor] pre > code) {
  padding: 14px 16px;
  border-radius: 10px;
  background: #1f2c34;
  color: #eff5f7;
}

:deep(.w-e-text-container [data-slate-editor] .w-e-image-container) {
  display: block;
  max-width: min(100%, 860px);
  margin: 22px auto;
}

:deep(.w-e-text-container [data-slate-editor] .w-e-image-container img) {
  width: 100%;
  border-radius: 10px;
  box-shadow: none;
}

:deep(.w-e-text-container [data-slate-editor] table) {
  width: 100%;
  margin: 18px 0;
  border-collapse: collapse;
  overflow: hidden;
  border-radius: 14px;
}

:deep(.w-e-text-container [data-slate-editor] th),
:deep(.w-e-text-container [data-slate-editor] td) {
  border: 1px solid rgba(33, 41, 52, 0.08);
  padding: 10px 12px;
}

@media (max-width: 768px) {
  .rich-editor__body {
    min-height: 240px;
  }

  .rich-editor__foot {
    flex-direction: column;
    align-items: flex-start;
  }

  :deep(.w-e-text-container [data-slate-editor]) {
    min-height: 240px;
    padding: 16px 16px 20px;
    font-size: 15px;
  }

  :deep(.w-e-text-placeholder) {
    top: 16px;
    left: 16px;
  }
}
</style>
