<script setup>
import { computed, reactive, watch } from 'vue'

import AdminStatusBadge from '@/components/admin/AdminStatusBadge.vue'

const props = defineProps({
  detail: {
    type: Object,
    required: true,
  },
  categories: {
    type: Array,
    default: () => [],
  },
  actionLoading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['approve', 'reject', 'ban', 'update-category'])

const form = reactive({
  reason: '',
  categoryIds: [],
})

// canReview 判断当前文章是否仍处于待审核状态，只有待审核文章才能通过或驳回
const canReview = computed(() => props.detail?.article?.statusCode === 1)

// canBan 判断当前文章是否还能执行封禁操作，避免已封禁文章重复显示封禁按钮
const canBan = computed(() => props.detail?.article?.statusCode !== 4)

watch(
  () => props.detail,
  (value) => {
    form.reason = value?.rejectReason || ''
    form.categoryIds = value?.article?.categoryIds?.slice() || []
  },
  { immediate: true },
)
</script>

<template>
  <section class="page-section panel">
    <div class="section-heading">
      <div>
        <h2>审核操作</h2>
        <p>处理文章审核、封禁和分类调整。</p>
      </div>
      <AdminStatusBadge :status="detail.article.statusCode" />
    </div>

    <div class="panel-grid">
      <label class="panel-block">
        <span>审核备注 / 驳回原因</span>
        <textarea v-model="form.reason" class="textarea-field" rows="5"></textarea>
      </label>

      <label class="panel-block">
        <span>文章分类</span>
        <select v-model="form.categoryIds" class="select-field" multiple size="6">
          <option v-for="category in categories" :key="category.id" :value="category.id">
            {{ category.name }}
          </option>
        </select>
      </label>
    </div>

    <div class="panel-actions">
      <button v-if="canReview" type="button" class="primary-btn" :disabled="actionLoading" @click="emit('approve')">
        {{ actionLoading ? '处理中...' : '通过审核' }}
      </button>
      <button type="button" class="secondary-btn" :disabled="actionLoading" @click="emit('update-category', form.categoryIds)">
        {{ actionLoading ? '处理中...' : '更新分类' }}
      </button>
      <button v-if="canReview" type="button" class="secondary-btn" :disabled="actionLoading" @click="emit('reject', form.reason)">
        {{ actionLoading ? '处理中...' : '驳回文章' }}
      </button>
      <button v-if="canBan" type="button" class="danger-btn" :disabled="actionLoading" @click="emit('ban', form.reason)">
        {{ actionLoading ? '处理中...' : '封禁文章' }}
      </button>
    </div>
  </section>
</template>

<style scoped>
.panel {
  padding: 24px;
}

.panel-grid {
  display: grid;
  gap: 18px;
}

.panel-block {
  display: grid;
  gap: 10px;
}

.panel-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 20px;
}
</style>
