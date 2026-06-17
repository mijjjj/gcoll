<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue'
import { NAlert, NButton, NDataTable, NForm, NFormItem, NInput, NModal, NSelect, NTag, NSwitch, useMessage, type DataTableColumns } from 'naive-ui'
import { CirclePlus, Play, PlugZap, RefreshCw, Save } from '@lucide/vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { DeviceItem, PointItem, RuntimeEvent } from '../api/console'

const consoleStore = useConsoleStore()
const message = useMessage()

const showDeviceModal = ref(false)
const showPointModal = ref(false)
const configText = ref('{}')

const deviceForm = reactive({
  name: '',
  code: '',
  groupId: '',
  pluginId: '',
  description: '',
})

const pointForm = reactive({
  name: '',
  description: '',
  address: '',
  valueType: 'float',
  unit: '',
  metadataText: '{}',
})

const selectedDevice = computed(() => consoleStore.selectedDevice)
const configPage = computed(() => consoleStore.devicePluginConfigPage)
const pluginOptions = computed(() => consoleStore.plugins.filter((plugin) => plugin.type === 'southbound').map((plugin) => ({ label: plugin.name, value: plugin.id })))
const groupOptions = computed(() => consoleStore.deviceGroups.map((group) => ({ label: group.name, value: group.id })))
const valueTypeOptions = [
  { label: 'bool', value: 'bool' },
  { label: 'int', value: 'int' },
  { label: 'float', value: 'float' },
  { label: 'string', value: 'string' },
  { label: 'bytes', value: 'bytes' },
  { label: 'datetime', value: 'datetime' },
  { label: 'json', value: 'json' },
]

watch(
  () => configPage.value?.config,
  (config) => {
    configText.value = JSON.stringify(config ?? {}, null, 2)
  },
  { immediate: true },
)

const deviceColumns: DataTableColumns<DeviceItem> = [
  { title: '设备名称', key: 'name', minWidth: 160 },
  { title: '编码', key: 'code', minWidth: 130 },
  { title: '插件', key: 'pluginName', minWidth: 160 },
  {
    title: '状态',
    key: 'status',
    width: 110,
    render: (row) => h(StatusBadge, { label: statusLabel(row.status), kind: row.status === 'online' ? 'success' : row.status === 'error' ? 'error' : 'neutral' }),
  },
  {
    title: '启用',
    key: 'enabled',
    width: 90,
    render: (row) => h(NTag, { size: 'small', type: row.enabled ? 'success' : 'default', bordered: false }, { default: () => (row.enabled ? '已启用' : '未启用') }),
  },
  { title: '点位', key: 'pointCount', width: 86 },
  { title: '最近在线', key: 'lastSeenAt', minWidth: 150 },
  {
    title: '操作',
    key: 'actions',
    width: 96,
    render: (row) =>
      h(
        NButton,
        {
          size: 'small',
          text: true,
          type: row.id === consoleStore.selectedDeviceId ? 'primary' : 'default',
          onClick: () => consoleStore.selectDevice(row.id),
        },
        { default: () => '查看' },
      ),
  },
]

const pointColumns: DataTableColumns<PointItem> = [
  { title: '点位名称', key: 'name', minWidth: 160 },
  { title: '地址', key: 'address', minWidth: 160 },
  { title: '值类型', key: 'valueType', width: 100 },
  { title: '单位', key: 'unit', width: 90 },
  {
    title: '状态',
    key: 'enabled',
    width: 100,
    render: (row) => h(NTag, { size: 'small', type: row.enabled ? 'success' : 'default', bordered: false }, { default: () => (row.enabled ? '启用' : '停用') }),
  },
  {
    title: 'metadata',
    key: 'metadata',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render: (row) => JSON.stringify(row.metadata ?? {}),
  },
]

const eventColumns: DataTableColumns<RuntimeEvent> = [
  { title: '时间', key: 'time', minWidth: 160 },
  { title: '级别', key: 'level', width: 90 },
  { title: '来源', key: 'source', width: 110 },
  { title: '消息', key: 'message', minWidth: 260 },
]

onMounted(async () => {
  await Promise.all([consoleStore.loadPlugins(), consoleStore.loadDevices(), consoleStore.loadTasks()])
})

function statusLabel(status: string) {
  const labels: Record<string, string> = {
    online: '在线',
    offline: '离线',
    disabled: '禁用',
    error: '异常',
  }
  return labels[status] ?? status
}

function openDeviceModal() {
  deviceForm.name = ''
  deviceForm.code = ''
  deviceForm.groupId = consoleStore.deviceGroups[0]?.id ?? ''
  deviceForm.pluginId = pluginOptions.value[0]?.value ?? ''
  deviceForm.description = ''
  showDeviceModal.value = true
}

async function createDevice() {
  await consoleStore.createDevice({ ...deviceForm })
  showDeviceModal.value = false
  message.success('设备已创建')
}

function openPointModal() {
  pointForm.name = ''
  pointForm.description = ''
  pointForm.address = ''
  pointForm.valueType = 'float'
  pointForm.unit = ''
  pointForm.metadataText = '{}'
  showPointModal.value = true
}

async function createPoint() {
  const metadata = parseJson(pointForm.metadataText)
  await consoleStore.createPoint({
    name: pointForm.name,
    description: pointForm.description,
    address: pointForm.address,
    valueType: pointForm.valueType,
    unit: pointForm.unit,
    metadata,
  })
  showPointModal.value = false
  message.success('点位已创建')
}

async function saveConfig() {
  const config = parseJson(configText.value)
  await consoleStore.saveDevicePluginConfig(config)
  message.success('设备插件配置已保存')
}

async function testConnection() {
  const result = await consoleStore.testDevicePluginConnection()
  if (result?.success) {
    message.success(result.message)
  } else if (result?.message) {
    message.warning(result.message)
  }
}

async function startTask() {
  await consoleStore.startSelectedDeviceTask()
  message.success('采集任务已启动')
}

function parseJson(value: string): Record<string, unknown> {
  try {
    const parsed = JSON.parse(value || '{}')
    if (!parsed || Array.isArray(parsed) || typeof parsed !== 'object') {
      throw new Error('JSON 必须是对象')
    }
    return parsed as Record<string, unknown>
  } catch (error) {
    throw new Error(error instanceof Error ? error.message : 'JSON 格式无效')
  }
}
</script>

<template>
  <PageHeader title="设备与点位" description="按设备绑定南向插件，保存设备插件配置和通用点位表。">
    <template #actions>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadDevices">
        <template #icon><RefreshCw :size="16" /></template>
        刷新
      </NButton>
      <NButton type="primary" :disabled="!pluginOptions.length || !groupOptions.length" @click="openDeviceModal">
        <template #icon><CirclePlus :size="16" /></template>
        新建设备
      </NButton>
    </template>
  </PageHeader>

  <div class="device-workspace">
    <SectionPanel title="设备列表" subtitle="设备只绑定插件 ID，协议细节由插件解释。">
      <NDataTable size="small" :columns="deviceColumns" :data="consoleStore.devices" :bordered="false" :single-line="false" />
    </SectionPanel>

    <SectionPanel title="插件配置" :subtitle="selectedDevice ? `${selectedDevice.name} / ${selectedDevice.pluginName}` : '请选择设备'">
      <template #actions>
        <NButton size="small" secondary :disabled="!selectedDevice" @click="testConnection">
          <template #icon><PlugZap :size="14" /></template>
          测试
        </NButton>
        <NButton size="small" type="primary" :disabled="!selectedDevice" @click="saveConfig">
          <template #icon><Save :size="14" /></template>
          保存
        </NButton>
      </template>

      <NAlert v-if="configPage && !configPage.configured" type="warning" :bordered="false">
        当前设备还没有已保存的插件运行配置，保存配置后才能启动采集。
      </NAlert>
      <NAlert v-else-if="!selectedDevice" type="info" :bordered="false">请选择左侧设备。</NAlert>

      <NInput v-model:value="configText" type="textarea" :autosize="{ minRows: 8, maxRows: 16 }" placeholder="{}" />

      <div v-if="configPage?.warnings?.length" class="warning-list">
        <NTag v-for="warning in configPage.warnings" :key="warning" size="small" type="warning" bordered>{{ warning }}</NTag>
      </div>
    </SectionPanel>
  </div>

  <SectionPanel title="通用点位表" subtitle="地址和 metadata 由当前设备绑定的插件解释。">
    <template #actions>
      <NButton size="small" secondary :disabled="!selectedDevice" @click="openPointModal">
        <template #icon><CirclePlus :size="14" /></template>
        新增点位
      </NButton>
      <NButton size="small" type="primary" :disabled="!selectedDevice || !configPage?.configured || !consoleStore.points.length" @click="startTask">
        <template #icon><Play :size="14" /></template>
        启动采集
      </NButton>
    </template>
    <NDataTable size="small" :columns="pointColumns" :data="consoleStore.points" :bordered="false" :single-line="false" />
  </SectionPanel>

  <SectionPanel title="最近事件" subtitle="插件调试摘要和控制面日志统一进入运行事件。">
    <NDataTable size="small" :columns="eventColumns" :data="configPage?.recentEvents ?? []" :bordered="false" :single-line="false" />
  </SectionPanel>

  <NModal v-model:show="showDeviceModal" preset="dialog" title="新建设备" positive-text="创建" negative-text="取消" @positive-click="createDevice">
    <NForm label-placement="top">
      <NFormItem label="设备名称"><NInput v-model:value="deviceForm.name" /></NFormItem>
      <NFormItem label="设备编码"><NInput v-model:value="deviceForm.code" /></NFormItem>
      <NFormItem label="设备分组"><NSelect v-model:value="deviceForm.groupId" :options="groupOptions" /></NFormItem>
      <NFormItem label="南向插件"><NSelect v-model:value="deviceForm.pluginId" :options="pluginOptions" /></NFormItem>
      <NFormItem label="说明"><NInput v-model:value="deviceForm.description" type="textarea" /></NFormItem>
    </NForm>
  </NModal>

  <NModal v-model:show="showPointModal" preset="dialog" title="新增点位" positive-text="新增" negative-text="取消" @positive-click="createPoint">
    <NForm label-placement="top">
      <NFormItem label="点位名称"><NInput v-model:value="pointForm.name" /></NFormItem>
      <NFormItem label="地址"><NInput v-model:value="pointForm.address" /></NFormItem>
      <NFormItem label="值类型"><NSelect v-model:value="pointForm.valueType" :options="valueTypeOptions" /></NFormItem>
      <NFormItem label="单位"><NInput v-model:value="pointForm.unit" /></NFormItem>
      <NFormItem label="说明"><NInput v-model:value="pointForm.description" /></NFormItem>
      <NFormItem label="metadata JSON"><NInput v-model:value="pointForm.metadataText" type="textarea" :autosize="{ minRows: 4, maxRows: 8 }" /></NFormItem>
      <NFormItem label="启用"><NSwitch :value="true" disabled /></NFormItem>
    </NForm>
  </NModal>
</template>
