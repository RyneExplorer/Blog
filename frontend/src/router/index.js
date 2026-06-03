import { createRouter, createWebHistory } from 'vue-router'

import ArticleListView from '@/views/ArticleListView.vue'
import ArticleDetailView from '@/views/ArticleDetailView.vue'
import CategoriesView from '@/views/CategoriesView.vue'
import MyArticlesView from '@/views/MyArticlesView.vue'
import FavoritesView from '@/views/FavoritesView.vue'
import ProfileView from '@/views/ProfileView.vue'
import SecurityView from '@/views/SecurityView.vue'
import EditorView from '@/views/EditorView.vue'
import AdminReviewsView from '@/views/AdminReviewsView.vue'
import AdminReviewDetailView from '@/views/AdminReviewDetailView.vue'
import AdminUsersView from '@/views/AdminUsersView.vue'
import pinia from '@/stores'
import { useAuthStore } from '@/stores/authStore'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: {
        name: 'articles',
      },
    },
    {
      path: '/articles',
      name: 'articles',
      component: ArticleListView,
    },
    {
      path: '/articles/:id',
      name: 'article-detail',
      component: ArticleDetailView,
      props: true,
    },
    {
      path: '/categories',
      name: 'categories',
      component: CategoriesView,
    },
    {
      path: '/me/articles',
      name: 'my-articles',
      component: MyArticlesView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/me/favorites',
      name: 'favorites',
      component: FavoritesView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/me/profile',
      name: 'profile',
      component: ProfileView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/me/security',
      name: 'security',
      component: SecurityView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/editor/new',
      name: 'editor-new',
      component: EditorView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/editor/:id',
      name: 'editor-edit',
      component: EditorView,
      props: true,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/admin/reviews',
      name: 'admin-reviews',
      component: AdminReviewsView,
      meta: {
        requiresAuth: true,
        requiresAdmin: true,
      },
    },
    {
      path: '/admin/reviews/:id',
      name: 'admin-review-detail',
      component: AdminReviewDetailView,
      props: true,
      meta: {
        requiresAuth: true,
        requiresAdmin: true,
      },
    },
    {
      path: '/admin/users',
      name: 'admin-users',
      component: AdminUsersView,
      meta: {
        requiresAuth: true,
        requiresAdmin: true,
      },
    },
  ],
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore(pinia)
  await authStore.ensureInitialized()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    authStore.openAuth('login', to.fullPath)
    return false
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return { name: 'articles' }
  }

  return true
})

export default router
