<script setup lang="ts">
import { h, onMounted } from 'vue'
import { NButton, NDataTable, NTag, type DataTableColumns } from 'naive-ui'
import { Pause, RefreshCw } from '@lucide/vue'
import DataToolbar from '../components/common/DataToolbar.vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { PointCacheItem } from '../api/console'

const consoleStore = useConsoleStore()

const columns: DataTableColumns<PointCacheItem> = [
  { title: '点位名称', key: 'pointName', minWidth: 150, sorter: 'default' },
  { title: '设备 ID', key: 'deviceId', minWidth: 180 },
  { title: '当前值', key: 'value', width: 140, className: 'value-cell' },
  {
    title: '质量',
    key: 'quality',
    width: 110,
    render: (row) => h(StatusBadge, { label: row.quality === 'good' ? '良好' : row.quality === 'uncertain' ? '不确定' : '异常', kind: row.quality === 'good' ? 'success' : row.quality === 'uncertain' ? 'warning' : 'error' }),
  },
  {
    title: '变化',
    key: 'changed',
    width: 96,
    render: (row) => h(NTag, { size: 'small', type: row.changed ? 'info' : 'default', bordered: false }, { default: () => (row.changed ? '已变化' : '未变化') }),
  },
  { title: '更新时间', key: 'updatedAt', minWidth: 180 },
]

onMounted(() => {
  consoleStore.loadPointCache()
})
</script>

<template>
  <PageHeader title="点位缓存" description="只展示最新点位值，不提供历史采集明细查询。">
    <template #actions>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadPointCache">
        <template #icon><RefreshCw :size="16" /></template>
        刷新
      </NButton>
      <NButton type="primary" secondary>
        <template #icon><Pause :size="16" /></template>
        暂停刷新
      </NButton>
    </template>
  </PageHeader>

  <SectionPanel>
    <DataToolbar placeholder="搜索点位、设备或质量" />
    <NDataTable size="small" :columns="columns" :data="consoleStore.pointCache" :bordered="false" :single-line="false" />
  </SectionPanel>
</template>
