<script setup>
const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  modelValue: {
    type: String,
    default: '',
  },
  target: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['update:modelValue', 'close', 'submit'])
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="reply-sheet">
      <button type="button" class="reply-sheet__mask" aria-label="关闭回复窗口" @click="emit('close')"></button>

      <div class="reply-sheet__panel">
        <div class="reply-sheet__grab"></div>

        <div class="reply-sheet__head">
          <div>
            <strong>回复 @{{ target?.author?.nickname || '用户' }}</strong>
            <p>{{ target?.comment?.content || '输入你的回复内容' }}</p>
          </div>
          <button type="button" class="reply-sheet__close" @click="emit('close')">取消</button>
        </div>

        <textarea
          :value="modelValue"
          class="textarea-field reply-sheet__textarea"
          rows="4"
          placeholder="写下你的回复"
          @input="emit('update:modelValue', $event.target.value)"
        />

        <div class="reply-sheet__actions">
          <button type="button" class="primary-btn" @click="emit('submit')">发送回复</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.reply-sheet {
  position: fixed;
  inset: 0;
  z-index: 60;
}

.reply-sheet__mask {
  position: absolute;
  inset: 0;
  border: none;
  background: rgba(19, 30, 31, 0.34);
  cursor: pointer;
}

.reply-sheet__panel {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  margin: 0 auto;
  width: min(760px, calc(100vw - 24px));
  padding: 14px 18px 18px;
  background: #fffdf8;
  border-radius: 24px 24px 0 0;
  box-shadow: 0 -24px 60px rgba(25, 45, 46, 0.16);
}

.reply-sheet__grab {
  width: 54px;
  height: 5px;
  margin: 0 auto 14px;
  border-radius: 999px;
  background: rgba(93, 109, 109, 0.2);
}

.reply-sheet__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.reply-sheet__head strong {
  display: block;
  color: #1f2d3d;
  font-size: 16px;
}

.reply-sheet__head p {
  margin: 6px 0 0;
  color: #7b8798;
  font-size: 13px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.reply-sheet__close {
  border: none;
  background: transparent;
  color: var(--color-muted);
  cursor: pointer;
  font-size: 14px;
}

.reply-sheet__textarea {
  min-height: 132px;
  border-radius: 18px;
  padding: 16px 18px;
}

.reply-sheet__actions {
  margin-top: 14px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 640px) {
  .reply-sheet__panel {
    width: calc(100vw - 12px);
    padding: 12px 14px 16px;
    border-radius: 20px 20px 0 0;
  }

  .reply-sheet__head {
    gap: 10px;
  }

  .reply-sheet__head strong {
    font-size: 15px;
  }
}
</style>
