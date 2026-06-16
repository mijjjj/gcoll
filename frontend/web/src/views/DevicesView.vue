<script setup lang="ts">
import { computed, h, onMounted, ref } from 'vue'
import {
  NButton,
  NDataTable,
  NDescriptions,
  NDescriptionsItem,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NRadioButton,
  NRadioGroup,
  NSwitch,
  NTag,
  type DataTableColumns,
} from 'naive-ui'
import { Bug, ChevronDown, Download, MoreHorizontal, Network, Play, Plus, RefreshCw, Search, Upload } from '@lucide/vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { ModbusTcpDebugLog, ModbusTcpPoint, ModbusTcpReadBlock, PointItem } from '../api/console'

const consoleStore = useConsoleStore()
const keyword = ref('')

const selectedDevice = computed(() => consoleStore.selectedDevice)
const protocolPage = computed(() => consoleStore.modbusTcpDeviceConfigPage)
const protocolConfig = computed(() => protocolPage.value?.config)

const filteredDevices = computed(() => {
  const text = keyword.value.trim().toLowerCase()
  if (!text) return consoleStore.devices
  return consoleStore.devices.filter((device) => `${device.name} ${device.code} ${device.pluginName}`.toLowerCase().includes(text))
})

const groupedDevices = computed(() =>
  consoleStore.deviceGroups.map((group) => ({
    ...group,
    devices: filteredDevices.value.filter((device) => device.groupId === group.id),
  })),
)

function statusLabel(status: string) {
  return status === 'online' ? '在线' : status === 'running' ? '运行中' : status === 'offline' ? '离线' : '已停止'
}

function statusKind(status: string) {
  if (status === 'online' || status === 'running') return 'success'
  if (status === 'offline') return 'error'
  return 'neutral'
}

const areaLabel = (area: string) => {
  const labels: Record<string, string> = {
    coil: '线圈',
    discrete_input: '离散输入',
    holding_register: '保持寄存器',
    input_register: '输入寄存器',
  }
  return labels[area] ?? area
}

const pointColumns: DataTableColumns<PointItem> = [
  { type: 'selection', width: 44 },
  { title: '点位名称', key: 'name', width: 150 },
  { title: '描述', key: 'description', minWidth: 150 },
  { title: '地址', key: 'address', width: 160 },
  { title: '值类型', key: 'valueType', width: 110 },
  { title: '单位', key: 'unit', width: 88 },
  {
    title: '状态',
    key: 'enabled',
    width: 100,
    render: (row) => h(StatusBadge, { label: row.enabled ? '已启用' : '已停用', kind: row.enabled ? 'success' : 'neutral' }),
  },
  {
    title: '标签',
    key: 'tags',
    minWidth: 180,
    render: (row) =>
      h(
        'div',
        { class: 'tag-list' },
        Object.entries(row.tags).map(([key, value]) => h(NTag, { size: 'small', bordered: false }, { default: () => `${key}: ${value}` })),
      ),
  },
]

const modbusPointColumns: DataTableColumns<ModbusTcpPoint> = [
  { type: 'selection', width: 44 },
  { title: '点位名称', key: 'name', width: 140 },
  {
    title: '区域',
    key: 'area',
    width: 118,
    render: (row) => h(NTag, { size: 'small', bordered: false }, { default: () => areaLabel(row.area) }),
  },
  { title: '地址', key: 'address', width: 84 },
  { title: '长度', key: 'quantity', width: 74 },
  { title: '值类型', key: 'valueType', width: 96 },
  {
    title: '模式',
    key: 'mode',
    width: 84,
    render: (row) => h(NTag, { size: 'small', type: row.mode === 'write' ? 'warning' : 'info', bordered: false }, { default: () => (row.mode === 'write' ? '写入' : '读取') }),
  },
  { title: '当前值', key: 'current', minWidth: 116 },
  {
    title: '质量',
    key: 'quality',
    width: 86,
    render: (row) => h(StatusBadge, { label: row.quality === 'good' ? '良好' : '不确定', kind: row.quality === 'good' ? 'success' : 'warning' }),
  },
  { title: '最近读取', key: 'lastReadAt', width: 174 },
]

const planColumns: DataTableColumns<ModbusTcpReadBlock> = [
  { title: '区域', key: 'area', width: 116, render: (row) => areaLabel(row.area) },
  { title: '起始地址', key: 'start', width: 88 },
  { title: '读取长度', key: 'quantity', width: 88 },
  { title: '最近耗时', key: 'latencyMs', width: 88, render: (row) => `${row.latencyMs} ms` },
]

const logColumns: DataTableColumns<ModbusTcpDebugLog> = [
  { title: '时间', key: 'time', width: 172 },
  { title: '级别', key: 'level', width: 72 },
  { title: '消息', key: 'message', minWidth: 140 },
  { title: '区域', key: 'area', width: 108, render: (row) => areaLabel(row.area) },
  { title: '耗时', key: 'costMs', width: 74, render: (row) => (row.costMs ? `${row.costMs} ms` : '-') },
]

onMounted(() => {
  consoleStore.loadDevices()
})
</script>

<template>
  <div class="device-console">
    <div class="device-top">
      <aside class="device-browser">
        <div class="panel-titlebar">
          <strong>设备列表</strong>
          <div class="mini-actions">
            <NButton quaternary circle size="small" aria-label="新增设备">
              <Plus :size="16" />
            </NButton>
            <NButton quaternary circle size="small" aria-label="刷新设备" :loading="consoleStore.loading" @click="consoleStore.loadDevices">
              <RefreshCw :size="15" />
            </NButton>
            <NButton quaternary circle size="small" aria-label="更多操作">
              <MoreHorizontal :size="16" />
            </NButton>
          </div>
        </div>

        <NInput v-model:value="keyword" size="small" placeholder="搜索设备名称、编号或插件" class="device-search">
          <template #prefix>
            <Search :size="15" />
          </template>
        </NInput>

        <div class="device-groups">
          <section v-for="group in groupedDevices" :key="group.id" class="device-group">
            <div class="group-title">
              <ChevronDown :size="14" />
              <span>{{ group.name }}（{{ group.devices.length }}）</span>
            </div>
            <div class="device-list">
              <button
                v-for="device in group.devices"
                :key="device.id"
                type="button"
                :class="['device-card', { 'device-card--active': consoleStore.selectedDeviceId === device.id }]"
                @click="consoleStore.selectDevice(device.id)"
              >
                <span :class="['device-dot', `device-dot--${device.status === 'online' ? 'online' : 'offline'}`]" />
                <span class="device-card__main">
                  <strong>{{ device.name }}</strong>
                  <small>{{ device.pluginName }}</small>
                </span>
                <span :class="['device-card__status', { 'is-online': device.status === 'online' }]">
                  {{ statusLabel(device.status) }}
                </span>
              </button>
            </div>
          </section>
        </div>

        <div class="device-count">共 {{ consoleStore.devices.length }} 个设备</div>
      </aside>

      <section class="device-detail-panel">
        <div class="detail-titlebar">
          <div class="detail-tab">
            <span>设备详情 - {{ selectedDevice?.name ?? '待选择' }}</span>
          </div>
        </div>

        <div v-if="selectedDevice" class="detail-body detail-body--single">
          <section class="config-form-shell">
            <NDescriptions label-placement="left" :column="2" size="small" bordered>
              <NDescriptionsItem label="设备编号">{{ selectedDevice.code }}</NDescriptionsItem>
              <NDescriptionsItem label="状态">
                <StatusBadge :label="statusLabel(selectedDevice.status)" :kind="statusKind(selectedDevice.status)" />
              </NDescriptionsItem>
              <NDescriptionsItem label="南向插件">{{ selectedDevice.pluginName }}</NDescriptionsItem>
              <NDescriptionsItem label="上报模式">
                {{ selectedDevice.reportMode === 'change' ? '变化上报' : '立即上报' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="点位数量">{{ selectedDevice.pointCount }}</NDescriptionsItem>
              <NDescriptionsItem label="最后在线">{{ selectedDevice.lastSeenAt }}</NDescriptionsItem>
              <NDescriptionsItem label="描述">{{ selectedDevice.description }}</NDescriptionsItem>
            </NDescriptions>

            <NForm v-if="protocolPage && protocolConfig" label-placement="left" label-width="112" class="device-protocol-form">
              <NFormItem label="主机地址">
                <NInput :value="protocolConfig.host" />
              </NFormItem>
              <NFormItem label="端口">
                <NInputNumber :value="protocolConfig.port" :min="1" :max="65535" />
              </NFormItem>
              <NFormItem label="Unit ID">
                <NInputNumber :value="protocolConfig.unitId" :min="0" :max="247" />
              </NFormItem>
              <NFormItem label="超时时间">
                <NInputNumber :value="protocolConfig.timeoutMs" :min="100" :max="60000">
                  <template #suffix>ms</template>
                </NInputNumber>
              </NFormItem>
              <NFormItem label="轮询间隔">
                <NInputNumber :value="protocolConfig.pollIntervalMs" :min="100">
                  <template #suffix>ms</template>
                </NInputNumber>
              </NFormItem>
              <NFormItem label="上报模式">
                <NRadioGroup :value="protocolConfig.reportMode">
                  <NRadioButton value="change">变化上报</NRadioButton>
                  <NRadioButton value="all">全部上报</NRadioButton>
                </NRadioGroup>
              </NFormItem>
              <NFormItem label="调试模式">
                <NSwitch :value="protocolConfig.debugEnabled" />
              </NFormItem>
            </NForm>

            <div v-if="protocolPage && protocolConfig" class="protocol-optimizer">
              <strong>协议优化参数</strong>
              <NDescriptions :column="2" size="small" label-placement="left" bordered>
                <NDescriptionsItem label="线圈批量上限">{{ protocolConfig.maxCoilBatch }}</NDescriptionsItem>
                <NDescriptionsItem label="寄存器批量上限">{{ protocolConfig.maxRegisterBatch }}</NDescriptionsItem>
                <NDescriptionsItem label="低延迟阈值">{{ protocolConfig.lowLatencyMs }} ms</NDescriptionsItem>
                <NDescriptionsItem label="高延迟阈值">{{ protocolConfig.highLatencyMs }} ms</NDescriptionsItem>
              </NDescriptions>
            </div>

            <div class="form-actions">
              <NButton>取消</NButton>
              <NButton type="primary">保存配置</NButton>
            </div>
          </section>

          <aside class="test-panel">
            <div class="test-header">
              <strong>测试结果</strong>
              <div class="test-actions">
                <NButton size="small" secondary type="primary" :disabled="!protocolPage" :loading="consoleStore.loading" @click="consoleStore.loadModbusTcpDeviceConfigPage">
                  <template #icon><Network :size="14" /></template>
                  测试连接
                </NButton>
                <NButton size="small" secondary :disabled="!protocolPage">
                  <template #icon><Bug :size="14" /></template>
                  调试模式
                </NButton>
              </div>
            </div>
            <div v-if="protocolPage" class="success-box">
              <StatusBadge label="最近测试通过" kind="success" />
              <span>{{ protocolPage.debugLogs[0]?.message ?? '等待执行测试' }}</span>
            </div>
            <dl class="test-meta">
              <div><dt>插件 ID</dt><dd>{{ selectedDevice.pluginId }}</dd></div>
              <div><dt>设备 ID</dt><dd>{{ selectedDevice.id }}</dd></div>
              <div><dt>配置作用域</dt><dd>当前设备</dd></div>
              <div><dt>权限状态</dt><dd>{{ selectedDevice.enabled ? '可编辑' : '只读' }}</dd></div>
              <div v-if="protocolConfig"><dt>目标地址</dt><dd>{{ protocolConfig.host }}:{{ protocolConfig.port }}</dd></div>
            </dl>

            <div v-if="protocolPage" class="read-plan-panel">
              <strong>读取计划</strong>
              <NDataTable size="small" :columns="planColumns" :data="protocolPage.readPlan" :bordered="false" />
            </div>

            <div v-if="protocolPage" class="read-plan-panel">
              <strong>调试日志</strong>
              <NDataTable size="small" :columns="logColumns" :data="protocolPage.debugLogs" :bordered="false" />
            </div>
          </aside>
        </div>
      </section>
    </div>

    <section class="points-panel">
      <div class="points-titlebar">
        <div class="points-title">
          <strong>{{ protocolPage?.points.length ? '协议点位配置' : '通用点位表' }}（共 {{ protocolPage?.points.length || consoleStore.points.length }} 个点）</strong>
          <div class="points-actions">
            <NButton size="small" secondary type="primary">
              <template #icon><Upload :size="14" /></template>
              导入点位
            </NButton>
            <NButton size="small" secondary type="primary">
              <template #icon><Download :size="14" /></template>
              导出点位
            </NButton>
            <NButton size="small" secondary :loading="consoleStore.loading" @click="consoleStore.loadSelectedDeviceDetails">
              <template #icon><RefreshCw :size="14" /></template>
              刷新
            </NButton>
          </div>
        </div>
        <NButton type="primary" size="small">
          <template #icon><Play :size="14" /></template>
          启动任务
        </NButton>
      </div>

      <NDataTable
        v-if="protocolPage?.points.length"
        size="small"
        :columns="modbusPointColumns"
        :data="protocolPage.points"
        :bordered="false"
        :single-line="false"
        class="points-table"
      />
      <NDataTable v-else size="small" :columns="pointColumns" :data="consoleStore.points" :bordered="false" :single-line="false" class="points-table" />
    </section>
  </div>
</template>
