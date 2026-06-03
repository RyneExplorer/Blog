<script setup>
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue', 'submit'])
</script>

<template>
  <div class="comment-composer">
    <textarea
      :value="modelValue"
      class="textarea-field comment-composer__textarea"
        rows="2"
      placeholder="输入你的评论内容"
      @input="emit('update:modelValue', $event.target.value)"
    />

    <div class="comment-composer__actions">
      <button type="button" class="primary-btn" @click="emit('submit')">发布评论</button>
    </div>
  </div>
</template>

<style scoped>
.comment-composer {
  display: grid;
  gap: 14px;
}

.comment-composer__textarea {
  min-height: 74px;
  border-radius: 22px;
  padding: 14px 18px;
  resize: vertical;
}

.comment-composer__actions {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 640px) {
  .comment-composer__textarea {
    min-height: 64px;
    border-radius: 18px;
    padding: 14px 16px;
  }
}
</style>
