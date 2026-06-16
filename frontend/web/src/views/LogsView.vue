<script setup lang="ts">
import { h, onMounted } from 'vue'
import { NButton, NDataTable, NTag, type DataTableColumns } from 'naive-ui'
import { Download, Pause, RefreshCw, Trash2 } from '@lucide/vue'
import DataToolbar from '../components/common/DataToolbar.vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import { useConsoleStore } from '../stores/console'
import type { RuntimeEvent } from '../api/console'

const consoleStore = useConsoleStore()

const columns: DataTableColumns<RuntimeEvent> = [
  { title: '时间', key: 'time', width: 168, sorter: 'default' },
  {
    title: '级别',
    key: 'level',
    width: 88,
    render: (row) => h(NTag, { size: 'small', type: row.level === 'WARN' ? 'warning' : 'default', bordered: false }, { default: () => row.level }),
  },
  { title: '来源', key: 'source', width: 120 },
  { title: '插件', key: 'pluginId', minWidth: 180, ellipsis: { tooltip: true } },
  { title: '设备', key: 'deviceId', minWidth: 170, ellipsis: { tooltip: true } },
  { title: '任务', key: 'taskId', minWidth: 150, ellipsis: { tooltip: true } },
  { title: '内容', key: 'message', minWidth: 260, ellipsis: { tooltip: true } },
  { title: 'Trace', key: 'traceId', minWidth: 140, ellipsis: { tooltip: true } },
]

onMounted(() => {
  consoleStore.loadLogs()
})
</script>

<template>
  <PageHeader title="日志诊断" description="按来源、级别、插件、设备和任务查看运行事件。">
    <template #actions>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadLogs">
        <template #icon><RefreshCw :size="16" /></template>
        刷新
      </NButton>
      <NButton type="primary" secondary>
        <template #icon><Download :size="16" /></template>
        导出
      </NButton>
    </template>
  </PageHeader>

  <SectionPanel>
    <DataToolbar placeholder="搜索日志内容、插件、设备或 Trace">
      <template #actions>
        <NButton size="small" secondary>
          <template #icon><Pause :size="14" /></template>
          暂停
        </NButton>
        <NButton size="small" secondary>
          <template #icon><Trash2 :size="14" /></template>
          清空视图
        </NButton>
      </template>
    </DataToolbar>
    <NDataTable size="small" :columns="columns" :data="consoleStore.logs" :bordered="false" :single-line="false" />
  </SectionPanel>
</template>
