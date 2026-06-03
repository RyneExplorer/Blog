<script setup>
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import UploadField from '@/components/common/UploadField.vue'
import { useProfileStore } from '@/stores/profileStore'

const profileStore = useProfileStore()
const { profileLoading, savingProfile, uploadingAvatar, profileForm } = storeToRefs(profileStore)

onMounted(() => {
  profileStore.loadProfile()
})
</script>

<template>
  <section class="profile-view page-section">
    <div class="section-heading">
      <div>
        <h1>个人资料</h1>
        <p>更新昵称、头像、邮箱和个人简介。</p>
      </div>
    </div>

    <div v-if="profileLoading" class="empty-box">正在加载个人资料...</div>

    <div v-else class="profile-grid">
      <UploadField
        :model-value="profileForm.avatar"
        :uploading="uploadingAvatar"
        :preview-aspect="1"
        label="上传头像"
        helper-text="支持 JPG、PNG、WEBP、GIF，选择后将直接上传。"
        @select="profileStore.uploadAvatarFile"
      />

      <div class="form-grid">
        <label class="field-block">
          <span>昵称</span>
          <input v-model="profileForm.nickname" class="field" />
        </label>
        <label class="field-block">
          <span>邮箱</span>
          <input v-model="profileForm.email" class="field" />
        </label>
        <label class="field-block field-block--full">
          <span>个人简介</span>
          <textarea v-model="profileForm.bio" class="textarea-field" rows="6"></textarea>
        </label>
      </div>
    </div>

    <div class="actions">
      <button type="button" class="primary-btn" :disabled="savingProfile" @click="profileStore.saveProfile">
        {{ savingProfile ? '保存中...' : '保存资料' }}
      </button>
    </div>
  </section>
</template>

<style scoped>
.profile-view {
  padding: 24px;
}

.profile-grid {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 20px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.field-block {
  display: grid;
  gap: 8px;
}

.field-block--full {
  grid-column: 1 / -1;
}

.actions {
  margin-top: 22px;
}

@media (max-width: 900px) {
  .profile-grid,
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
