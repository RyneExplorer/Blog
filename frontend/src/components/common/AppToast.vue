<script setup>
import { storeToRefs } from 'pinia'

import { useUiStore } from '@/stores/uiStore'

const uiStore = useUiStore()
const { toast } = storeToRefs(uiStore)
</script>

<template>
  <transition name="toast-fade">
    <div v-if="toast.visible" class="app-toast" :class="`app-toast--${toast.type}`">
      {{ toast.text }}
    </div>
  </transition>
</template>

<style scoped>
.app-toast {
  position: fixed;
  right: 20px;
  bottom: 20px;
  z-index: 140;
  min-width: 220px;
  max-width: 360px;
  padding: 14px 16px;
  border-radius: 16px;
  color: #fff;
  box-shadow: var(--shadow-card);
}

.app-toast--info,
.app-toast--success {
  background: linear-gradient(135deg, #204645, #5da39e);
}

.app-toast--error {
  background: linear-gradient(135deg, #7c3428, #aa5748);
}

.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: all 0.2s ease;
}

.toast-fade-enter-from,
.toast-fade-leave-to {
  opacity: 0;
  transform: translateY(12px);
}
</style>
