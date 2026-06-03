import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', () => {
  const toast = ref({
    text: '',
    type: 'info',
    visible: false,
  })
  const confirmDialog = ref({
    visible: false,
    title: '',
    message: '',
    confirmText: '确认',
    cancelText: '取消',
    intent: 'primary',
  })
  const mobileDrawerOpen = ref(false)
  let toastTimer = 0
  let confirmResolver = null

  function showToast(text, type = 'info') {
    toast.value = {
      text,
      type,
      visible: true,
    }

    if (toastTimer) {
      window.clearTimeout(toastTimer)
    }

    toastTimer = window.setTimeout(() => {
      toast.value.visible = false
      toastTimer = 0
    }, 2600)
  }

  function hideToast() {
    toast.value.visible = false
  }

  function openDrawer() {
    mobileDrawerOpen.value = true
  }

  function closeDrawer() {
    mobileDrawerOpen.value = false
  }

  function openConfirm(options = {}) {
    confirmDialog.value = {
      visible: true,
      title: options.title || '请确认操作',
      message: options.message || '确认继续执行当前操作吗？',
      confirmText: options.confirmText || '确认',
      cancelText: options.cancelText || '取消',
      intent: options.intent || 'primary',
    }

    return new Promise((resolve) => {
      confirmResolver = resolve
    })
  }

  function resolveConfirm(result) {
    confirmDialog.value.visible = false
    if (confirmResolver) {
      confirmResolver(Boolean(result))
      confirmResolver = null
    }
  }

  return {
    toast,
    confirmDialog,
    mobileDrawerOpen,
    showToast,
    hideToast,
    openDrawer,
    closeDrawer,
    openConfirm,
    resolveConfirm,
  }
})
