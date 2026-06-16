<script setup lang="ts">
import { Search } from '@lucide/vue'
import { NInput } from 'naive-ui'

const keyword = defineModel<string>('keyword', { default: '' })

withDefaults(
  defineProps<{
    placeholder?: string
  }>(),
  {
    placeholder: '搜索关键词',
  },
)
</script>

<template>
  <div class="data-toolbar">
    <NInput v-model:value="keyword" clearable :placeholder="placeholder" class="data-toolbar__search">
      <template #prefix>
        <Search :size="15" />
      </template>
    </NInput>
    <div class="data-toolbar__filters">
      <slot />
    </div>
    <div v-if="$slots.actions" class="data-toolbar__actions">
      <slot name="actions" />
    </div>
  </div>
</template>
