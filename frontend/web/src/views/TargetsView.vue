<script setup lang="ts">
import { h, onMounted } from 'vue'
import { NButton, NDataTable, type DataTableColumns } from 'naive-ui'
import { Plus, RefreshCw, Send } from '@lucide/vue'
import DataToolbar from '../components/common/DataToolbar.vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { ForwardTargetItem } from '../api/console'

const consoleStore = useConsoleStore()

const columns: DataTableColumns<ForwardTargetItem> = [
  { title: '目标名称', key: 'name', minWidth: 180, sorter: 'default' },
  { title: '北向插件', key: 'pluginName', minWidth: 150 },
  {
    title: '状态',
    key: 'status',
    width: 110,
    render: (row) => h(StatusBadge, { label: row.status === 'running' ? '运行中' : '已停止', kind: row.status === 'running' ? 'success' : 'neutral' }),
  },
  { title: '端点', key: 'endpoint', minWidth: 260, ellipsis: { tooltip: true } },
  { title: '最近错误', key: 'lastError', minWidth: 150, ellipsis: { tooltip: true } },
  { title: '更新时间', key: 'updatedAt', width: 168 },
]

onMounted(() => {
  consoleStore.loadTargets()
})
</script>

<template>
  <PageHeader title="北向转发" description="管理转发目标、路由状态和最近错误。">
    <template #actions>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadTargets">
        <template #icon><RefreshCw :size="16" /></template>
        刷新
      </NButton>
      <NButton type="primary">
        <template #icon><Plus :size="16" /></template>
        新建目标
      </NButton>
    </template>
  </PageHeader>

  <SectionPanel>
    <DataToolbar placeholder="搜索目标、插件或端点">
      <template #actions>
        <NButton size="small" secondary>
          <template #icon><Send :size="14" /></template>
          测试转发
        </NButton>
      </template>
    </DataToolbar>
    <NDataTable size="small" :columns="columns" :data="consoleStore.targets" :bordered="false" :single-line="false" />
  </SectionPanel>
</template>
