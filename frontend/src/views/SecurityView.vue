<script setup>
import { storeToRefs } from 'pinia'

import { useProfileStore } from '@/stores/profileStore'

const profileStore = useProfileStore()
const { passwordForm, savingPassword } = storeToRefs(profileStore)
</script>

<template>
  <section class="page-section security-view">
    <div class="section-heading">
      <div>
        <h1>账号安全</h1>
        <p>修改当前账号密码。</p>
      </div>
    </div>

    <div class="security-grid">
      <label class="field-block">
        <span>旧密码</span>
        <input v-model="passwordForm.old_password" type="password" class="field" />
      </label>
      <label class="field-block">
        <span>新密码</span>
        <input v-model="passwordForm.new_password" type="password" class="field" />
      </label>
    </div>

    <div class="actions">
      <button type="button" class="primary-btn" :disabled="savingPassword" @click="profileStore.savePassword">
        {{ savingPassword ? '更新中...' : '更新密码' }}
      </button>
    </div>
  </section>
</template>

<style scoped>
.security-view {
  padding: 24px;
}

.security-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.field-block {
  display: grid;
  gap: 8px;
}

.actions {
  margin-top: 20px;
}

@media (max-width: 760px) {
  .security-grid {
    grid-template-columns: 1fr;
  }
}
</style>
