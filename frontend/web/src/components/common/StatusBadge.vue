<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    label: string
    kind?: 'success' | 'warning' | 'error' | 'info' | 'neutral'
  }>(),
  {
    kind: 'neutral',
  },
)

const className = computed(() => ['status-badge', `status-badge--${props.kind}`])
</script>

<template>
  <span :class="className">
    <slot name="icon">
      <span class="status-dot" />
    </slot>
    <span>{{ label }}</span>
  </span>
</template>
