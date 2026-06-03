<script setup>
import CommentItem from '@/components/comment/CommentItem.vue'

defineProps({
  comments: {
    type: Array,
    default: () => [],
  },
  articleId: {
    type: Number,
    required: true,
  },
})
</script>

<template>
  <div class="comment-tree">
    <CommentItem v-for="node in comments" :key="node.comment.id" :node="node" :article-id="articleId" :depth="0" />
  </div>
</template>

<style scoped>
.comment-tree {
  display: grid;
  gap: 0;
  margin-top: 18px;
}
</style>
