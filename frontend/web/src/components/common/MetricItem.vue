<script setup lang="ts">
withDefaults(
  defineProps<{
    label: string
    value: string | number
    hint?: string
    tone?: 'primary' | 'success' | 'warning' | 'error' | 'neutral'
  }>(),
  {
    tone: 'primary',
  },
)
</script>

<template>
  <div class="metric-item">
    <div :class="['metric-item__icon', `metric-item__icon--${tone}`]">
      <slot name="icon" />
    </div>
    <div class="metric-item__body">
      <span>{{ label }}</span>
      <strong>{{ value }}</strong>
      <small v-if="hint">{{ hint }}</small>
    </div>
  </div>
</template>
