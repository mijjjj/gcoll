<script setup lang="ts">
import { h, onMounted } from 'vue'
import { NButton, NDataTable, NSwitch, type DataTableColumns } from 'naive-ui'
import { FlaskConical, Plus, RefreshCw } from '@lucide/vue'
import DataToolbar from '../components/common/DataToolbar.vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import { useConsoleStore } from '../stores/console'
import type { PipelineRuleItem } from '../api/console'

const consoleStore = useConsoleStore()

const columns: DataTableColumns<PipelineRuleItem> = [
  { title: '规则名称', key: 'name', minWidth: 180, sorter: 'default' },
  {
    title: '启用',
    key: 'enabled',
    width: 88,
    render: (row) => h(NSwitch, { value: row.enabled, size: 'small' }),
  },
  { title: '条件表达式', key: 'expression', minWidth: 260, ellipsis: { tooltip: true } },
  { title: '命中次数', key: 'matched', width: 110, align: 'right', sorter: 'default' },
  { title: '目标数', key: 'targetCount', width: 90, align: 'right' },
  { title: '更新时间', key: 'updatedAt', width: 168 },
]

onMounted(() => {
  consoleStore.loadPipelineRules()
})
</script>

<template>
  <PageHeader title="规则过滤" description="维护规则列表、条件表达式和测试样本入口。">
    <template #actions>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadPipelineRules">
        <template #icon><RefreshCw :size="16" /></template>
        刷新
      </NButton>
      <NButton type="primary">
        <template #icon><Plus :size="16" /></template>
        新建规则
      </NButton>
    </template>
  </PageHeader>

  <SectionPanel>
    <DataToolbar placeholder="搜索规则名称或表达式">
      <template #actions>
        <NButton size="small" secondary>
          <template #icon><FlaskConical :size="14" /></template>
          测试样本
        </NButton>
      </template>
    </DataToolbar>
    <NDataTable size="small" :columns="columns" :data="consoleStore.pipelineRules" :bordered="false" :single-line="false" />
  </SectionPanel>
</template>
