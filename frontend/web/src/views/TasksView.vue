<script setup lang="ts">
import { h, onMounted } from 'vue'
import { NButton, NDataTable, NTag, type DataTableColumns } from 'naive-ui'
import { Play, RefreshCw, RotateCw, Square } from '@lucide/vue'
import DataToolbar from '../components/common/DataToolbar.vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { TaskSummary } from '../api/console'

const consoleStore = useConsoleStore()

const columns: DataTableColumns<TaskSummary> = [
  { type: 'selection', width: 44 },
  { title: '任务名称', key: 'name', minWidth: 180, sorter: 'default' },
  { title: '设备', key: 'deviceName', minWidth: 140 },
  { title: '南向插件', key: 'southPluginName', minWidth: 140 },
  { title: '点位数', key: 'pointCount', width: 90, align: 'right', sorter: 'default' },
  {
    title: '上报模式',
    key: 'reportMode',
    width: 112,
    render: (row) => h(NTag, { size: 'small', type: row.reportMode === 'change' ? 'info' : 'warning', bordered: false }, { default: () => (row.reportMode === 'change' ? '变化上报' : '立即上报') }),
  },
  {
    title: '状态',
    key: 'status',
    width: 108,
    render: (row) => h(StatusBadge, { label: row.status === 'running' ? '运行中' : '已停止', kind: row.status === 'running' ? 'success' : 'neutral' }),
  },
  { title: '采集速率', key: 'rate', width: 112, align: 'right' },
  { title: '规则命中率', key: 'ruleHitRate', width: 112, align: 'right' },
  { title: '最后采集时间', key: 'lastCollectedAt', minWidth: 156 },
]

onMounted(() => {
  consoleStore.loadTasks()
})
</script>

<template>
  <PageHeader title="采集任务" description="管理采集任务生命周期、上报模式、吞吐和规则命中率。">
    <template #actions>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadTasks">
        <template #icon><RefreshCw :size="16" /></template>
        刷新
      </NButton>
      <NButton type="primary">
        <template #icon><Play :size="16" /></template>
        新建任务
      </NButton>
    </template>
  </PageHeader>

  <SectionPanel>
    <DataToolbar placeholder="搜索任务、设备或插件">
      <template #actions>
        <NButton size="small" secondary>
          <template #icon><Play :size="14" /></template>
          启动
        </NButton>
        <NButton size="small" secondary>
          <template #icon><Square :size="14" /></template>
          停止
        </NButton>
        <NButton size="small" secondary>
          <template #icon><RotateCw :size="14" /></template>
          重启
        </NButton>
      </template>
    </DataToolbar>
    <NDataTable size="small" :columns="columns" :data="consoleStore.tasks" :bordered="false" :single-line="false" />
  </SectionPanel>
</template>
