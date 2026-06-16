<script setup lang="ts">
defineProps<{
  title?: string
  subtitle?: string
  compact?: boolean
}>()
</script>

<template>
  <section :class="['section-panel', { 'section-panel--compact': compact }]">
    <div v-if="title || subtitle || $slots.actions" class="section-panel__header">
      <div>
        <h2 v-if="title">{{ title }}</h2>
        <p v-if="subtitle">{{ subtitle }}</p>
      </div>
      <div v-if="$slots.actions" class="section-panel__actions">
        <slot name="actions" />
      </div>
    </div>
    <slot />
  </section>
</template>
