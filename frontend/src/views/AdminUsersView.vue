<script setup>
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import PaginationBar from '@/components/article/PaginationBar.vue'
import { useAdminStore } from '@/stores/adminStore'

const adminStore = useAdminStore()
const { userList, userPagination, userFilters, userLoading } = storeToRefs(adminStore)

onMounted(() => {
  adminStore.loadUsers()
})
</script>

<template>
  <section class="page-section admin-users-view">
    <div class="section-heading">
      <div>
        <h1>用户管理</h1>
        <p>分页查看普通用户与管理员账号。</p>
      </div>
    </div>

    <div class="filter-grid">
      <input v-model="userFilters.username" class="field" placeholder="用户名" />
      <input v-model="userFilters.nickname" class="field" placeholder="昵称" />
      <select v-model="userFilters.status" class="select-field">
        <option value="">全部状态</option>
        <option value="1">正常</option>
        <option value="2">禁用</option>
      </select>
      <button type="button" class="primary-btn" @click="adminStore.loadUsers({ page: 1 })">筛选</button>
    </div>

    <div v-if="userLoading" class="empty-box">用户列表加载中...</div>
    <div v-else-if="!userList.length" class="empty-box">暂无用户</div>

    <div v-else class="table-wrap">
      <table class="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>用户名</th>
            <th>昵称</th>
            <th>邮箱</th>
            <th>角色</th>
            <th>状态</th>
            <th>更新时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in userList" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.username }}</td>
            <td>{{ user.nickname }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.role === 0 ? '管理员' : '普通用户' }}</td>
            <td>{{ user.status === 1 ? '正常' : '禁用' }}</td>
            <td>{{ user.updatedAt }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <PaginationBar :pagination="userPagination" @change="adminStore.loadUsers({ page: $event })" />
  </section>
</template>

<style scoped>
.admin-users-view {
  padding: 24px;
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}

.table-wrap {
  overflow: auto;
}

.admin-table {
  width: 100%;
  border-collapse: collapse;
  background: rgba(255, 255, 255, 0.94);
  border-radius: 20px;
  overflow: hidden;
}

.admin-table th,
.admin-table td {
  padding: 14px 12px;
  border-bottom: 1px solid var(--color-border);
  text-align: left;
}

.admin-table thead {
  background: rgba(93, 163, 158, 0.08);
}

@media (max-width: 900px) {
  .filter-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 640px) {
  .filter-grid {
    grid-template-columns: 1fr;
  }
}
</style>
