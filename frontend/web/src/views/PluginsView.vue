<script setup lang="ts">
import { h, onMounted } from 'vue'
import { NButton, NDataTable, NTag, type DataTableColumns } from 'naive-ui'
import { PackagePlus, RefreshCw, RotateCcw, Upload } from '@lucide/vue'
import DataToolbar from '../components/common/DataToolbar.vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { PluginItem } from '../api/console'

const consoleStore = useConsoleStore()

const columns: DataTableColumns<PluginItem> = [
  { title: '插件名称', key: 'name', minWidth: 170, sorter: 'default' },
  {
    title: '类型',
    key: 'type',
    width: 112,
    render: (row) => h(NTag, { size: 'small', type: row.type === 'southbound' ? 'info' : 'success', bordered: false }, { default: () => (row.type === 'southbound' ? '南向插件' : '北向插件') }),
  },
  { title: '版本', key: 'version', width: 100 },
  { title: '运行时', key: 'runtime', width: 110 },
  { title: '协议', key: 'protocol', width: 100 },
  {
    title: '状态',
    key: 'status',
    width: 110,
    render: (row) => h(StatusBadge, { label: row.status === 'running' ? '运行中' : '已停止', kind: row.status === 'running' ? 'success' : 'neutral' }),
  },
  {
    title: '权限',
    key: 'permissions',
    minWidth: 220,
    render: (row) =>
      h(
        'div',
        { class: 'tag-list' },
        row.permissions.map((permission) => h(NTag, { size: 'small', bordered: false }, { default: () => permission })),
      ),
  },
  { title: '更新时间', key: 'updatedAt', width: 168 },
  { title: '插件 ID', key: 'id', minWidth: 220, ellipsis: { tooltip: true } },
]

onMounted(() => {
  consoleStore.loadPlugins()
})
</script>

<template>
  <PageHeader title="插件中心" description="当前阶段聚焦本地插件导入、启停、升级、回滚和调试。">
    <template #actions>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadPlugins">
        <template #icon><RefreshCw :size="16" /></template>
        刷新
      </NButton>
      <NButton type="primary">
        <template #icon><Upload :size="16" /></template>
        导入插件
      </NButton>
    </template>
  </PageHeader>

  <SectionPanel>
    <DataToolbar placeholder="搜索插件名称、权限或协议">
      <template #actions>
        <NButton size="small" secondary>
          <template #icon><PackagePlus :size="14" /></template>
          启用
        </NButton>
        <NButton size="small" secondary>
          <template #icon><RotateCcw :size="14" /></template>
          回滚
        </NButton>
      </template>
    </DataToolbar>
    <NDataTable size="small" :columns="columns" :data="consoleStore.plugins" :bordered="false" :single-line="false" />
  </SectionPanel>
</template>
