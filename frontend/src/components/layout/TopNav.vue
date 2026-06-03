<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import { useAuthStore } from '@/stores/authStore'
import { useUiStore } from '@/stores/uiStore'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const uiStore = useUiStore()
const { isLoggedIn, isAdmin, currentUser, avatarText } = storeToRefs(authStore)
const userMenuOpen = ref(false)
const userMenuRef = ref(null)

const navItems = computed(() => {
  const items = [
    { label: '文章', to: { name: 'articles' } },
    { label: '分类', to: { name: 'categories' } },
    { label: '我的文章', to: { name: 'my-articles' }, auth: true },
    { label: '我的收藏', to: { name: 'favorites' }, auth: true },
  ]

  if (isAdmin.value) {
    items.push(
      { label: '审核中心', to: { name: 'admin-reviews' } },
      { label: '用户管理', to: { name: 'admin-users' } },
    )
  }

  return items
})

function jump(item) {
  if (item.auth && !isLoggedIn.value) {
    authStore.openAuth('login', router.resolve(item.to).fullPath)
    return
  }

  userMenuOpen.value = false
  router.push(item.to)
}

function toggleUserMenu() {
  userMenuOpen.value = !userMenuOpen.value
}

function goProfile() {
  userMenuOpen.value = false
  router.push({ name: 'profile' })
}

async function logout() {
  userMenuOpen.value = false
  const confirmed = await uiStore.openConfirm({
    title: '确认退出登录？',
    message: '退出后你仍然可以浏览公开文章，但评论、收藏和写作功能会暂时不可用。',
    confirmText: '退出登录',
    cancelText: '取消',
    intent: 'danger',
  })

  if (!confirmed) {
    return
  }

  await authStore.logout()
  router.push({ name: 'articles' })
}

function handleDocumentClick(event) {
  if (!userMenuRef.value) {
    return
  }

  if (!userMenuRef.value.contains(event.target)) {
    userMenuOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleDocumentClick)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleDocumentClick)
})
</script>

<template>
  <header class="top-nav">
    <div class="top-nav__inner">
      <button type="button" class="menu-button" @click="uiStore.openDrawer">☰</button>

      <RouterLink :to="{ name: 'articles' }" class="brand">
        <span class="brand__mark">尘</span>
        <div class="brand__text">
          <strong>汇尘轩</strong>
          <span>Blog Console</span>
        </div>
      </RouterLink>

      <nav class="nav-links">
        <button
          v-for="item in navItems"
          :key="item.label"
          type="button"
          class="nav-link"
          :class="{ 'is-active': route.name === item.to.name }"
          @click="jump(item)"
        >
          {{ item.label }}
        </button>
      </nav>

      <div class="top-nav__actions">
        <button type="button" class="secondary-btn write-btn" @click="jump({ auth: true, to: { name: 'editor-new' } })">
          写文章
        </button>

        <button v-if="!isLoggedIn" type="button" class="primary-btn" @click="authStore.openAuth('login')">
          登录 / 注册
        </button>

        <div v-else ref="userMenuRef" class="user-dropdown">
          <button type="button" class="user-entry" @click.stop="toggleUserMenu">
            <span class="user-entry__avatar">
              <img v-if="currentUser?.avatar" :src="currentUser.avatar" :alt="currentUser?.nickname || currentUser?.username" />
              <template v-else>{{ avatarText }}</template>
            </span>
            <span class="user-entry__meta">
              <span class="user-entry__name">{{ currentUser?.nickname || currentUser?.username }}</span>
              <span class="user-entry__sub">{{ currentUser?.bio || currentUser?.email || '已登录用户' }}</span>
            </span>
          </button>

          <transition name="menu-fade">
            <div v-if="userMenuOpen" class="user-menu">
              <button type="button" class="user-menu__item" @click="goProfile">个人资料</button>
              <button type="button" class="user-menu__item user-menu__item--danger" @click="logout">退出登录</button>
            </div>
          </transition>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.top-nav {
  position: sticky;
  top: 0;
  z-index: 90;
  background: rgba(246, 244, 238, 0.82);
  backdrop-filter: blur(18px);
  border-bottom: 1px solid var(--color-border);
}

.top-nav__inner {
  width: min(1380px, calc(100vw - 32px));
  margin: 0 auto;
  min-height: 78px;
  display: flex;
  align-items: center;
  gap: 18px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--color-text);
}

.brand__mark {
  width: 42px;
  height: 42px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, var(--color-primary-dark), var(--color-primary));
  color: #fff;
  font-weight: 700;
}

.brand__text strong {
  display: block;
  font-size: 17px;
  font-weight: 700;
}

.brand__text span {
  color: var(--color-muted);
  font-size: 12px;
}

.nav-links {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  gap: 8px;
}

.nav-link {
  border: none;
  background: transparent;
  padding: 10px 14px;
  border-radius: 999px;
  color: var(--color-muted);
  cursor: pointer;
}

.nav-link.is-active {
  background: var(--color-primary-soft);
  color: var(--color-primary-dark);
}

.top-nav__actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.write-btn {
  min-height: 54px;
  padding: 0 24px;
  border-radius: 18px;
  font-size: 17px;
}

.menu-button {
  display: none;
  width: 42px;
  height: 42px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fff;
  cursor: pointer;
}

.user-entry {
  display: flex;
  align-items: center;
  gap: 10px;
  border: none;
  background: linear-gradient(135deg, rgb(34, 88, 87), rgb(73, 139, 136));
  border-radius: 999px;
  padding: 6px 8px 6px 6px;
  color: #f7f1e6;
  cursor: pointer;
}

.user-entry__avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  overflow: hidden;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, var(--color-primary-dark), var(--color-primary));
  color: #fff;
  font-weight: 700;
}

.user-entry__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-entry__name {
  display: block;
  max-width: 120px;
  color: #f7f1e6;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-entry__meta {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  min-width: 0;
}

.user-entry__sub {
  max-width: 140px;
  color: rgba(247, 241, 230, 0.78);
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-dropdown {
  position: relative;
}

.user-menu {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  min-width: 180px;
  padding: 8px;
  border: 1px solid var(--color-border);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: var(--shadow-soft);
}

.user-menu__item {
  width: 100%;
  border: none;
  background: transparent;
  border-radius: 12px;
  padding: 13px 14px;
  text-align: left;
  cursor: pointer;
}

.user-menu__item:hover {
  background: var(--color-primary-soft);
}

.user-menu__item--danger:hover {
  background: var(--color-danger-soft);
  color: var(--color-danger);
}

.menu-fade-enter-active,
.menu-fade-leave-active {
  transition: all 0.16s ease;
}

.menu-fade-enter-from,
.menu-fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

@media (max-width: 900px) {
  .top-nav__inner {
    width: min(100vw - 24px, 1180px);
  }

  .menu-button {
    display: inline-grid;
    place-items: center;
  }

  .nav-links,
  .write-btn {
    display: none;
  }

  .user-entry__sub {
    display: none;
  }
}
</style>
