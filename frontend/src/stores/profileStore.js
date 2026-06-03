import { reactive, ref } from 'vue'
import { defineStore } from 'pinia'

import { changePassword, getProfile, updateProfile, uploadAvatar } from '@/api/user'
import { resolveAssetUrl } from '@/api/http'
import { useAuthStore } from '@/stores/authStore'
import { useUiStore } from '@/stores/uiStore'

export const useProfileStore = defineStore('profile', () => {
  const authStore = useAuthStore()
  const uiStore = useUiStore()

  const profileLoading = ref(false)
  const savingProfile = ref(false)
  const uploadingAvatar = ref(false)
  const savingPassword = ref(false)

  const profileForm = reactive({
    nickname: '',
    email: '',
    bio: '',
    avatar: '',
  })

  const passwordForm = reactive({
    old_password: '',
    new_password: '',
  })

  function syncForm(user) {
    profileForm.nickname = user?.nickname || ''
    profileForm.email = user?.email || ''
    profileForm.bio = user?.bio || ''
    profileForm.avatar = resolveAssetUrl(user?.avatar || '')
  }

  // loadProfile 加载个人资料，默认优先复用登录初始化时已经拿到的用户信息
  async function loadProfile(options = {}) {
    if (!options.force && authStore.currentUser) {
      syncForm(authStore.currentUser)
      return
    }

    profileLoading.value = true

    try {
      const data = await getProfile()
      authStore.applyUser(data)
      syncForm(authStore.currentUser)
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      profileLoading.value = false
    }
  }

  async function saveProfile() {
    savingProfile.value = true

    try {
      await updateProfile({
        nickname: profileForm.nickname,
        email: profileForm.email,
        bio: profileForm.bio,
      })
      await loadProfile({ force: true })
      uiStore.showToast('资料已更新', 'success')
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      savingProfile.value = false
    }
  }

  async function uploadAvatarFile(file) {
    uploadingAvatar.value = true

    try {
      const data = await uploadAvatar(file)
      profileForm.avatar = resolveAssetUrl(data?.avatar || data?.url || '')
      if (authStore.currentUser) {
        authStore.applyUser({
          ...authStore.currentUser,
          avatar: profileForm.avatar,
        })
      }
      uiStore.showToast('头像上传成功', 'success')
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      uploadingAvatar.value = false
    }
  }

  async function savePassword() {
    if (passwordForm.new_password.length < 6) {
      uiStore.showToast('新密码至少 6 位', 'error')
      return
    }

    savingPassword.value = true

    try {
      await changePassword(passwordForm)
      passwordForm.old_password = ''
      passwordForm.new_password = ''
      uiStore.showToast('密码已更新', 'success')
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      savingPassword.value = false
    }
  }

  return {
    profileLoading,
    savingProfile,
    uploadingAvatar,
    savingPassword,
    profileForm,
    passwordForm,
    syncForm,
    loadProfile,
    saveProfile,
    uploadAvatarFile,
    savePassword,
  }
})
