<script setup lang="ts">
import { computed, h, onMounted } from 'vue'
import { NButton, NDataTable, NDescriptions, NDescriptionsItem, NProgress, NTag, type DataTableColumns } from 'naive-ui'
import { Activity, Database, Plug, RefreshCw, Server, TriangleAlert, Wifi, Zap } from '@lucide/vue'
import DataToolbar from '../components/common/DataToolbar.vue'
import MetricItem from '../components/common/MetricItem.vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { RuntimeEvent, TaskSummary } from '../api/console'

const consoleStore = useConsoleStore()

const metricIcons = {
  runtime: Server,
  devices: Activity,
  points: Zap,
  plugins: Plug,
}

const runtimeStatus = computed(() => {
  if (consoleStore.error) return '离线'
  return consoleStore.overview?.runtime.status === 'running' ? '运行中' : '待连接'
})

const runtimeKind = computed(() => {
  if (consoleStore.error) return 'error'
  return consoleStore.overview?.runtime.status === 'running' ? 'success' : 'warning'
})

const taskColumns: DataTableColumns<TaskSummary> = [
  { title: '任务名称', key: 'name', sorter: 'default', minWidth: 168 },
  { title: '设备', key: 'deviceName', minWidth: 132 },
  { title: '南向插件', key: 'southPluginName', minWidth: 130 },
  { title: '点位数', key: 'pointCount', sorter: 'default', align: 'right', width: 88 },
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
  { title: '采集速率', key: 'rate', align: 'right', width: 112 },
  { title: '规则命中率', key: 'ruleHitRate', align: 'right', width: 112 },
  { title: '最后采集时间', key: 'lastCollectedAt', minWidth: 156 },
]

const eventColumns: DataTableColumns<RuntimeEvent> = [
  { title: '时间', key: 'time', width: 168 },
  {
    title: '级别',
    key: 'level',
    width: 88,
    render: (row) => h(NTag, { size: 'small', bordered: false, type: row.level === 'WARN' ? 'warning' : 'default' }, { default: () => row.level }),
  },
  { title: '来源', key: 'source', width: 124 },
  { title: '内容', key: 'message', ellipsis: { tooltip: true } },
]

onMounted(() => {
  consoleStore.loadOverview()
})
</script>

<template>
  <PageHeader title="工作台" description="运行状态、采集任务、插件进程和最近告警。">
    <template #actions>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadOverview">
        <template #icon>
          <RefreshCw :size="16" />
        </template>
        刷新
      </NButton>
      <NButton type="primary">
        <template #icon>
          <Activity :size="16" />
        </template>
        创建采集任务
      </NButton>
    </template>
  </PageHeader>

  <div class="metric-strip">
    <MetricItem
      v-for="metric in consoleStore.overview?.metrics ?? []"
      :key="metric.key"
      :label="metric.label"
      :value="metric.value"
      :hint="metric.hint"
      :tone="metric.tone"
    >
      <template #icon>
        <component :is="metricIcons[metric.key as keyof typeof metricIcons]" :size="19" />
      </template>
    </MetricItem>
  </div>

  <SectionPanel title="运行时健康">
    <template #actions>
      <StatusBadge :label="runtimeStatus" :kind="runtimeKind" />
    </template>
    <NDescriptions label-placement="left" :column="4" size="small">
      <NDescriptionsItem label="服务">
        {{ consoleStore.overview?.runtime.service ?? 'gcoll-server' }}
      </NDescriptionsItem>
      <NDescriptionsItem label="模式">
        {{ consoleStore.overview?.runtime.mode ?? 'server' }}
      </NDescriptionsItem>
      <NDescriptionsItem label="版本">
        {{ consoleStore.overview?.runtime.version ?? '0.1.0-dev' }}
      </NDescriptionsItem>
      <NDescriptionsItem label="检查时间">
        {{ consoleStore.overview?.runtime.checkedAt ?? '待连接' }}
      </NDescriptionsItem>
      <NDescriptionsItem label="数据库">
        <Database :size="14" /> {{ consoleStore.overview?.runtime.database ?? '待连接' }}
      </NDescriptionsItem>
      <NDescriptionsItem label="HTTP API">
        {{ consoleStore.overview?.runtime.apiBase ?? '待连接' }}
      </NDescriptionsItem>
      <NDescriptionsItem label="插件进程">
        {{ consoleStore.overview ? `${consoleStore.overview.pluginSummary.running}/${consoleStore.overview.pluginSummary.total}` : '待连接' }}
      </NDescriptionsItem>
      <NDescriptionsItem label="网络">
        <Wifi :size="14" /> {{ consoleStore.overview?.network.detail ?? '待连接' }}
      </NDescriptionsItem>
    </NDescriptions>
    <p v-if="consoleStore.error" class="error-text">{{ consoleStore.error }}</p>
  </SectionPanel>

  <div class="workbench-grid">
    <SectionPanel title="采集任务" subtitle="默认展示生命周期、吞吐和最后采集时间。">
      <template #actions>
        <NTag size="small" :bordered="false">MVP 闭环</NTag>
      </template>
      <DataToolbar placeholder="搜索任务、设备或插件">
        <template #actions>
          <NButton size="small" secondary :loading="consoleStore.loading" @click="consoleStore.loadOverview">
            <template #icon>
              <RefreshCw :size="14" />
            </template>
            刷新
          </NButton>
        </template>
      </DataToolbar>
      <NDataTable size="small" :columns="taskColumns" :data="consoleStore.tasks" :bordered="false" :single-line="false" />
    </SectionPanel>

    <SectionPanel title="数据面状态" subtitle="内存缓冲、规则过滤和北向转发的轻量状态。">
      <div class="pipeline-list">
        <div>
          <span>内存缓冲</span>
          <NProgress type="line" :percentage="consoleStore.overview?.dataPlane.queueUsagePercent ?? 0" />
        </div>
        <div>
          <span>规则命中率</span>
          <NProgress type="line" status="success" :percentage="consoleStore.overview?.dataPlane.ruleHitPercent ?? 0" />
        </div>
        <div>
          <span>北向转发</span>
          <NProgress type="line" status="warning" :percentage="consoleStore.overview?.dataPlane.forwardPercent ?? 0" />
        </div>
        <dl class="runtime-mini">
          <div><dt>实时吞吐</dt><dd>{{ consoleStore.overview?.dataPlane.throughput ?? '待连接' }}</dd></div>
          <div><dt>端到端延迟</dt><dd>{{ consoleStore.overview?.dataPlane.latency ?? '待连接' }}</dd></div>
          <div><dt>背压状态</dt><dd>{{ consoleStore.overview?.dataPlane.backpressure ?? '待连接' }}</dd></div>
        </dl>
      </div>
    </SectionPanel>
  </div>

  <SectionPanel title="最近事件">
    <template #actions>
      <StatusBadge label="含 1 条警告" kind="warning">
        <template #icon>
          <TriangleAlert :size="14" />
        </template>
      </StatusBadge>
    </template>
    <NDataTable size="small" :columns="eventColumns" :data="consoleStore.logs" :bordered="false" :pagination="false" />
  </SectionPanel>
</template>
