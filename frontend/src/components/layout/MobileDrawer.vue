<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import { useAuthStore } from '@/stores/authStore'
import { useUiStore } from '@/stores/uiStore'

const router = useRouter()
const uiStore = useUiStore()
const authStore = useAuthStore()
const { mobileDrawerOpen } = storeToRefs(uiStore)
const { isLoggedIn, isAdmin } = storeToRefs(authStore)

const links = computed(() => {
  const base = [
    { label: '文章', to: { name: 'articles' } },
    { label: '分类', to: { name: 'categories' } },
  ]

  if (isLoggedIn.value) {
    base.push(
      { label: '我的文章', to: { name: 'my-articles' } },
      { label: '我的收藏', to: { name: 'favorites' } },
      { label: '个人资料', to: { name: 'profile' } },
      { label: '写文章', to: { name: 'editor-new' } },
    )
  }

  if (isAdmin.value) {
    base.push(
      { label: '审核中心', to: { name: 'admin-reviews' } },
      { label: '用户管理', to: { name: 'admin-users' } },
    )
  }

  return base
})

function navigate(to) {
  uiStore.closeDrawer()
  router.push(to)
}
</script>

<template>
  <transition name="drawer-fade">
    <div v-if="mobileDrawerOpen" class="drawer-overlay" @click.self="uiStore.closeDrawer">
      <div class="drawer-panel">
        <div class="drawer-header">
          <strong>导航</strong>
          <button type="button" class="drawer-close" @click="uiStore.closeDrawer">×</button>
        </div>

        <div class="drawer-links">
          <button v-for="item in links" :key="item.label" type="button" class="drawer-link" @click="navigate(item.to)">
            {{ item.label }}
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.drawer-overlay {
  position: fixed;
  inset: 0;
  z-index: 120;
  background: rgba(27, 33, 32, 0.46);
}

.drawer-panel {
  width: min(320px, 92vw);
  height: 100%;
  padding: 20px;
  background: rgba(249, 247, 241, 0.98);
  box-shadow: var(--shadow-soft);
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.drawer-close {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.96);
  cursor: pointer;
}

.drawer-links {
  display: grid;
  gap: 12px;
}

.drawer-link {
  padding: 13px 14px;
  border: 1px solid var(--color-border);
  border-radius: 14px;
  background: #fff;
  text-align: left;
  cursor: pointer;
}

.drawer-fade-enter-active,
.drawer-fade-leave-active {
  transition: opacity 0.18s ease;
}

.drawer-fade-enter-from,
.drawer-fade-leave-to {
  opacity: 0;
}
</style>
