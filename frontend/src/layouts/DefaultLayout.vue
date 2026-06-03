<script setup>
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'

import TopNav from '@/components/layout/TopNav.vue'
import MobileDrawer from '@/components/layout/MobileDrawer.vue'
import AuthDialog from '@/components/auth/AuthDialog.vue'
import AppToast from '@/components/common/AppToast.vue'
import AppConfirmDialog from '@/components/common/AppConfirmDialog.vue'
import { useAuthStore } from '@/stores/authStore'
import { useArticleStore } from '@/stores/articleStore'

const authStore = useAuthStore()
const articleStore = useArticleStore()

onMounted(async () => {
  await authStore.initializeAuth()
  await articleStore.loadCategories()
})
</script>

<template>
  <div class="app-shell">
    <TopNav />

    <main class="content-area">
      <RouterView />
    </main>

    <MobileDrawer />
    <AuthDialog />
    <AppConfirmDialog />
    <AppToast />
  </div>
</template>

<style scoped>
.app-shell {
  min-height: 100vh;
}

.content-area {
  width: min(1380px, calc(100vw - 32px));
  margin: 0 auto;
  padding: 24px 0 64px;
  min-width: 0;
}

@media (max-width: 1080px) {
  .content-area {
    width: min(100vw - 24px, 1180px);
  }
}

@media (max-width: 900px) {
  .content-area {
    width: min(100vw - 24px, 1180px);
  }
}
</style>
