<script setup>
import { computed, ref } from 'vue'

import { useUiStore } from '@/stores/uiStore'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  label: {
    type: String,
    default: '上传图片',
  },
  uploading: {
    type: Boolean,
    default: false,
  },
  previewAspect: {
    type: Number,
    default: 1,
  },
  helperText: {
    type: String,
    default: '',
  },
  maxSizeMb: {
    type: Number,
    default: 5,
  },
  compact: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['select'])

const uiStore = useUiStore()
const fileInput = ref(null)

const previewStyle = computed(() => ({
  aspectRatio: String(props.previewAspect || 1),
}))

const previewActionText = computed(() => {
  if (props.uploading) {
    return '上传中...'
  }

  return props.modelValue ? `更换${props.label}` : props.label
})

// trigger 打开本地文件选择框
function trigger() {
  if (props.uploading) {
    return
  }

  fileInput.value?.click()
}

// resetInput 清空 input 的值，确保重复选择同一张图片也能触发 change 事件
function resetInput(event) {
  if (event?.target) {
    event.target.value = ''
  }
}

// isSupportedImage 判断前端是否支持当前图片类型
function isSupportedImage(file) {
  return ['image/jpeg', 'image/png', 'image/webp', 'image/gif'].includes(file.type)
}

// onChange 校验图片类型和大小，校验通过后直接上传原始文件
function onChange(event) {
  const file = event.target.files?.[0]
  resetInput(event)

  if (!file) {
    return
  }

  if (!isSupportedImage(file)) {
    uiStore.showToast('仅支持 JPG、PNG、WEBP 或 GIF 图片', 'error')
    return
  }

  if (file.size > props.maxSizeMb * 1024 * 1024) {
    uiStore.showToast(`图片大小不能超过 ${props.maxSizeMb}MB`, 'error')
    return
  }

  emit('select', file)
}
</script>

<template>
  <div class="upload-field" :class="{ 'upload-field--compact': compact }">
    <button
      type="button"
      class="upload-preview"
      :class="{ 'upload-preview--compact': compact }"
      :style="previewStyle"
      :disabled="uploading"
      @click="trigger"
    >
      <img v-if="modelValue" :src="modelValue" :alt="label" />

      <div v-else class="upload-placeholder">
        <span>{{ label }}</span>
        <small>{{ helperText || '选择图片后将直接上传' }}</small>
      </div>

      <div v-if="modelValue" class="upload-overlay">
        <span class="upload-overlay__action">{{ previewActionText }}</span>
        <small v-if="helperText && compact" class="upload-overlay__helper">{{ helperText }}</small>
      </div>
    </button>

    <div v-if="!compact" class="upload-actions">
      <button type="button" class="secondary-btn" :disabled="uploading" @click="trigger">
        {{ uploading ? '上传中...' : label }}
      </button>
      <p v-if="helperText" class="upload-helper">{{ helperText }}</p>
    </div>

    <input ref="fileInput" type="file" accept="image/*" class="hidden-input" :disabled="uploading" @change="onChange" />
  </div>
</template>

<style scoped>
.upload-field {
  display: grid;
  gap: 14px;
}

.upload-field--compact {
  gap: 0;
}

.upload-preview {
  position: relative;
  width: 100%;
  min-height: 180px;
  border: 1px dashed var(--color-border-strong);
  border-radius: 22px;
  overflow: hidden;
  background:
    linear-gradient(135deg, rgba(93, 163, 158, 0.08), rgba(255, 255, 255, 0.9)),
    rgba(255, 255, 255, 0.84);
  cursor: pointer;
  transition:
    transform 0.18s ease,
    box-shadow 0.18s ease,
    border-color 0.18s ease;
}

.upload-preview:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(35, 47, 49, 0.08);
  border-color: rgba(93, 163, 158, 0.32);
}

.upload-preview:disabled {
  cursor: default;
  opacity: 0.72;
}

.upload-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-placeholder {
  width: 100%;
  height: 100%;
  min-height: 180px;
  display: grid;
  place-items: center;
  text-align: center;
  gap: 6px;
  color: var(--color-muted);
  padding: 24px;
}

.upload-placeholder span {
  font-size: 18px;
  font-weight: 700;
  color: var(--color-primary-dark);
}

.upload-placeholder small,
.upload-helper {
  color: var(--color-muted);
  line-height: 1.6;
}

.upload-overlay {
  position: absolute;
  inset: auto 0 0 0;
  display: grid;
  gap: 4px;
  padding: 18px 16px 14px;
  background: linear-gradient(180deg, rgba(24, 65, 66, 0), rgba(24, 65, 66, 0.72));
  color: #fff;
  text-align: left;
}

.upload-overlay__action {
  font-size: 0.92rem;
  font-weight: 700;
}

.upload-overlay__helper {
  color: rgba(255, 255, 255, 0.78);
  line-height: 1.5;
}

.upload-actions {
  display: grid;
  gap: 8px;
}

.upload-preview--compact {
  min-height: 0;
}

.hidden-input {
  display: none;
}
</style>
