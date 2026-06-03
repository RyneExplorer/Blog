<script setup>
import { storeToRefs } from 'pinia'

import { useUiStore } from '@/stores/uiStore'

const uiStore = useUiStore()
const { confirmDialog } = storeToRefs(uiStore)

function cancel() {
  uiStore.resolveConfirm(false)
}

function confirm() {
  uiStore.resolveConfirm(true)
}
</script>

<template>
  <transition name="confirm-fade">
    <div v-if="confirmDialog.visible" class="confirm-overlay" @click.self="cancel">
      <div class="confirm-card">
        <div class="confirm-card__icon" :class="`confirm-card__icon--${confirmDialog.intent}`">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M12 8v5" />
            <circle cx="12" cy="16.5" r="0.6" fill="currentColor" stroke="none" />
            <path d="M10.3 3.9 2.9 16.7A1.5 1.5 0 0 0 4.2 19h15.6a1.5 1.5 0 0 0 1.3-2.3L13.7 3.9a1.5 1.5 0 0 0-2.6 0Z" />
          </svg>
        </div>

        <div class="confirm-card__body">
          <h2>{{ confirmDialog.title }}</h2>
          <p>{{ confirmDialog.message }}</p>
        </div>

        <div class="confirm-card__actions">
          <button type="button" class="secondary-btn confirm-card__btn" @click="cancel">
            {{ confirmDialog.cancelText }}
          </button>
          <button
            type="button"
            class="confirm-card__btn"
            :class="confirmDialog.intent === 'danger' ? 'danger-btn' : 'primary-btn'"
            @click="confirm"
          >
            {{ confirmDialog.confirmText }}
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.confirm-overlay {
  position: fixed;
  inset: 0;
  z-index: 110;
  display: grid;
  place-items: center;
  padding: 20px;
  background: rgba(43, 48, 47, 0.38);
  backdrop-filter: blur(10px);
}

.confirm-card {
  width: min(100%, 420px);
  padding: 24px;
  border: 1px solid rgba(255, 255, 255, 0.45);
  border-radius: 26px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.95), rgba(255, 251, 244, 0.92)),
    var(--color-surface);
  box-shadow: 0 24px 56px rgba(35, 47, 49, 0.16);
}

.confirm-card__icon {
  width: 48px;
  height: 48px;
  display: grid;
  place-items: center;
  border-radius: 16px;
  margin-bottom: 16px;
}

.confirm-card__icon--danger {
  color: var(--color-danger);
  background: rgba(170, 87, 72, 0.12);
}

.confirm-card__icon--primary {
  color: var(--color-primary-dark);
  background: rgba(93, 163, 158, 0.14);
}

.confirm-card__icon svg {
  width: 24px;
  height: 24px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.confirm-card__body h2 {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--color-primary-dark);
  line-height: 1.2;
}

.confirm-card__body p {
  margin-top: 10px;
  color: var(--color-muted);
  line-height: 1.7;
}

.confirm-card__actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.confirm-card__btn {
  min-width: 96px;
}

.confirm-fade-enter-active,
.confirm-fade-leave-active {
  transition: opacity 0.18s ease;
}

.confirm-fade-enter-from,
.confirm-fade-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .confirm-card {
    padding: 20px;
    border-radius: 22px;
  }

  .confirm-card__actions {
    flex-direction: column-reverse;
  }

  .confirm-card__btn {
    width: 100%;
  }
}
</style>
