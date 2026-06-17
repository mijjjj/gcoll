<script setup lang="ts">
import { onBeforeUnmount } from 'vue'
import { useMessage } from 'naive-ui'
import { setHttpErrorNotifier } from '../../api/http'

const message = useMessage()

setHttpErrorNotifier((content) => {
  message.error(content, {
    duration: 5000,
    keepAliveOnHover: true,
  })
})

onBeforeUnmount(() => {
  setHttpErrorNotifier(null)
})
</script>

<template></template>
