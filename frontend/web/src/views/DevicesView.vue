<script setup lang="ts">
import { computed, h, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import {
  NAlert,
  NButton,
  NDataTable,
  NDropdown,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NPopover,
  NSelect,
  NSwitch,
  NTabPane,
  NTabs,
  NTag,
  NTooltip,
  useDialog,
  useMessage,
  type DataTableColumns,
  type DropdownOption,
} from 'naive-ui'
import { ChevronDown, CirclePlus, FilePlus, FolderPlus, Play, PlugZap, RefreshCw, Save } from '@lucide/vue'
import PageHeader from '../components/common/PageHeader.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { DeviceGroup, DeviceItem, PointItem, RuntimeEvent } from '../api/console'
import { getApiBasePath } from '../api/http'
import { getCurrentLanguage } from '../i18n'

interface ConfigField {
  name: string
  title: string
  description: string
  type: string
  component: string
  defaultValue: unknown
  options: Array<{ label: string; value: string }>
}

type PointField = ConfigField

interface PointAddressExample {
  label: string
  value: string
  description: string
  group: string
}

interface PointAddressExampleGroup {
  label: string
  items: PointAddressExample[]
}

type ContextTarget = { type: 'blank' } | { type: 'group'; group: DeviceGroup } | { type: 'device'; device: DeviceItem }

function buildSchemaDefaults(schema: Record<string, unknown> | null | undefined) {
  const defaults: Record<string, unknown> = {}
  const properties = (schema?.properties ?? {}) as Record<string, Record<string, unknown>>
  for (const [name, property] of Object.entries(properties)) {
    if (Object.prototype.hasOwnProperty.call(property, 'default')) {
      defaults[name] = property.default
    }
  }
  return defaults
}

function injectPageScript(html: string, js: string) {
  const bridgeScript = `<script>
window.addEventListener('load', () => {
  window.parent.postMessage({ type: 'gcoll:page-ready' }, '*')
})
<\/script>`
  const script = `${bridgeScript}${js ? `\n<script>${js}<\/script>` : ''}`
  if (html.includes('</body>')) {
    return html.replace('</body>', `${script}\n  </body>`)
  }
  if (html.includes('</html>')) {
    return html.replace('</html>', `${script}\n</html>`)
  }
  return `${html}\n${script}`
}

function buildPointAddressExamples(schema: Record<string, unknown> | null | undefined) {
  const properties = (schema?.properties ?? {}) as Record<string, Record<string, unknown>>
  const addressField = (properties.address ?? {}) as Record<string, unknown>
  const ui = (addressField.ui ?? {}) as Record<string, unknown>
  const examples = Array.isArray(ui.examples) ? ui.examples : []
  return examples
    .map((item) => {
      const example = item as Record<string, unknown>
      const value = String(example.value ?? '').trim()
      if (!value) {
        return null
      }
      return {
        label: String(example.label ?? value),
        value,
        description: String(example.description ?? ''),
        group: String(example.group ?? '常用示例'),
      } satisfies PointAddressExample
    })
    .filter((item): item is PointAddressExample => Boolean(item))
}

function groupPointAddressExamples(examples: PointAddressExample[]) {
  const groups = new Map<string, PointAddressExample[]>()
  for (const example of examples) {
    const items = groups.get(example.group) ?? []
    items.push(example)
    groups.set(example.group, items)
  }
  return Array.from(groups.entries()).map(([label, items]) => ({ label, items })) satisfies PointAddressExampleGroup[]
}

function buildPluginConfigApis(deviceId: string | undefined) {
  const runtimeApiBase = consoleStore.overview?.runtime.apiBase?.trim()
  const apiBase = runtimeApiBase ? `${runtimeApiBase.replace(/\/+$/, '')}/api/v1` : getApiBasePath()
  if (!deviceId) {
    return {
      getConfig: '',
      saveConfig: '',
      testConfig: '',
    }
  }
  const basePath = `${apiBase}/devices/${deviceId}/protocol-config`
  return {
    getConfig: basePath,
    saveConfig: basePath,
    testConfig: `${basePath}/test`,
  }
}

function toMessageData<T>(value: T): T {
  return JSON.parse(JSON.stringify(value ?? null)) as T
}

const consoleStore = useConsoleStore()
const message = useMessage()
const dialog = useDialog()

const activeTab = ref('config')
const searchText = ref('')
const showGroupModal = ref(false)
const showDeviceModal = ref(false)
const showPointModal = ref(false)
const editingPointIndex = ref(-1)
const contextVisible = ref(false)
const contextX = ref(0)
const contextY = ref(0)
const contextTarget = ref<ContextTarget>({ type: 'blank' })
const iframeRef = ref<HTMLIFrameElement | null>(null)

const configDraft = ref<Record<string, unknown>>({})
const pointDraft = ref<PointItem[]>([])

const groupForm = reactive({ name: '' })
const deviceForm = reactive({
  name: '',
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
  enabled: true,
})
const pointMetadataDraft = ref<Record<string, unknown>>({})

const selectedDevice = computed(() => consoleStore.selectedDevice)
const configPage = computed(() => consoleStore.devicePluginConfigPage)
const customConfigPage = computed(() => configPage.value?.customConfigPage)
const customPointPage = computed(() => configPage.value?.customPointPage)
const pointSchema = computed(() => (configPage.value?.pointSchema ?? {}) as Record<string, unknown>)
const useCustomConfigPage = computed(() => Boolean(customConfigPage.value?.enabled && customConfigPage.value?.html))
const useCustomPointPage = computed(() => Boolean(customPointPage.value?.enabled && customPointPage.value?.html))
const pluginOptions = computed(() => consoleStore.plugins.filter((plugin) => plugin.type === 'southbound').map((plugin) => ({ label: plugin.name, value: plugin.id })))
const groupOptions = computed(() => consoleStore.deviceGroups.map((group) => ({ label: group.name, value: group.id })))
const filteredGroups = computed(() => {
  const keyword = searchText.value.trim().toLowerCase()
  return consoleStore.deviceGroups.map((group) => {
    const devices = consoleStore.devices.filter((device) => {
      if (device.groupId !== group.id) return false
      if (!keyword) return true
      return `${device.name} ${device.pluginName}`.toLowerCase().includes(keyword)
    })
    return { ...group, devices }
  })
})
const contextOptions = computed<DropdownOption[]>(() => {
  if (contextTarget.value.type === 'group') {
    return [
      { label: '在分组下添加设备', key: 'add-device' },
      { label: '删除分组', key: 'delete-group' },
    ]
  }
  if (contextTarget.value.type === 'device') {
    return [
      { label: '删除设备', key: 'delete-device' },
    ]
  }
  return [
    { label: '添加分组', key: 'add-group' },
    { label: '添加设备', key: 'add-device' },
  ]
})

const configFields = computed<ConfigField[]>(() => {
  const schema = configPage.value?.configSchema ?? {}
  return buildSchemaFields(schema)
})

const pointFields = computed<PointField[]>(() => {
  const schema = pointSchema.value ?? {}
  const hiddenFields = new Set(['name', 'description', 'address', 'valueType', 'unit', 'enabled'])
  return buildSchemaFields(schema).filter((field) => !hiddenFields.has(field.name))
})

function buildSchemaFields(schema: Record<string, unknown> | null | undefined) {
  const properties = (schema?.properties ?? {}) as Record<string, Record<string, unknown>>
  return Object.entries(properties).map(([name, property]) => {
    const ui = (property.ui ?? {}) as Record<string, unknown>
    const enumValues = Array.isArray(property.enum) ? property.enum : []
    return {
      name,
      title: String(property.title ?? name),
      description: String(property.description ?? ''),
      type: String(property.type ?? 'string'),
      component: String(ui.component ?? property.type ?? 'text'),
      defaultValue: property.default,
      options: enumValues.map((item) => ({ label: String(item), value: String(item) })),
    }
  })
}

const configSrcdoc = computed(() => {
  const page = customConfigPage.value
  if (!page?.html) return ''
  return injectPageScript(page.html, page.js)
})

const pointSrcdoc = computed(() => {
  const page = customPointPage.value
  if (!page?.html) return ''
  return injectPageScript(page.html, page.js)
})

const pointColumns = computed<DataTableColumns<PointItem>>(() => [
  { title: '点位名称', key: 'name', minWidth: 150 },
  { title: '地址', key: 'address', minWidth: 130 },
  { title: '值类型', key: 'valueType', width: 100 },
  { title: '单位', key: 'unit', width: 80 },
  {
    title: '启用',
    key: 'enabled',
    width: 86,
    render: (row) => h(NTag, { size: 'small', type: row.enabled ? 'success' : 'default', bordered: false }, { default: () => (row.enabled ? '启用' : '停用') }),
  },
  {
    title: '操作',
    key: 'actions',
    width: 126,
    render: (_row, index) => [
      h(NButton, { size: 'small', text: true, type: 'primary', onClick: () => openPointModal(index) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', text: true, type: 'error', onClick: () => removePoint(index) }, { default: () => '删除' }),
    ],
  },
])

const pointAddressExamples = computed(() => buildPointAddressExamples(pointSchema.value))
const pointAddressExampleGroups = computed(() => groupPointAddressExamples(pointAddressExamples.value))
const pointAddressDescription = computed(() => {
  const properties = (pointSchema.value?.properties ?? {}) as Record<string, Record<string, unknown>>
  return String(properties.address?.description ?? '')
})

const eventColumns: DataTableColumns<RuntimeEvent> = [
  { title: '时间', key: 'time', minWidth: 160 },
  { title: '级别', key: 'level', width: 80 },
  { title: '来源', key: 'source', width: 100 },
  { title: '消息', key: 'message', minWidth: 240, ellipsis: { tooltip: true } },
]

watch(
  () => configPage.value,
  (page) => {
    const defaults = buildSchemaDefaults((page?.configSchema ?? {}) as Record<string, unknown>)
    configDraft.value = { ...defaults, ...(page?.config ?? {}) }
    nextTick(postConfigToIframe)
  },
  { immediate: true },
)

watch(
  () => consoleStore.points,
  (points) => {
    pointDraft.value = points.map((point) => ({ ...point, tags: { ...point.tags }, metadata: { ...(point.metadata ?? {}) } }))
  },
  { immediate: true, deep: true },
)

onMounted(async () => {
  window.addEventListener('message', handlePluginMessage)
  await Promise.all([consoleStore.loadPlugins(), consoleStore.loadDevices(), consoleStore.loadTasks()])
})

onBeforeUnmount(() => {
  window.removeEventListener('message', handlePluginMessage)
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

function deviceStatusKind(status: string) {
  if (status === 'online') return 'success'
  if (status === 'error') return 'error'
  return 'neutral'
}

function fieldNumber(name: string, defaultValue: unknown) {
  const value = configDraft.value[name] ?? defaultValue
  const number = Number(value)
  return Number.isFinite(number) ? number : null
}

function fieldString(name: string, defaultValue: unknown) {
  return String(configDraft.value[name] ?? defaultValue ?? '')
}

function fieldBool(name: string, defaultValue: unknown) {
  return Boolean(configDraft.value[name] ?? defaultValue ?? false)
}

function pointFieldNumber(name: string, defaultValue: unknown) {
  const value = pointMetadataDraft.value[name] ?? defaultValue
  const number = Number(value)
  return Number.isFinite(number) ? number : null
}

function pointFieldString(name: string, defaultValue: unknown) {
  return String(pointMetadataDraft.value[name] ?? defaultValue ?? '')
}

function pointFieldBool(name: string, defaultValue: unknown) {
  return Boolean(pointMetadataDraft.value[name] ?? defaultValue ?? false)
}

function setConfigField(name: string, value: unknown) {
  configDraft.value = { ...configDraft.value, [name]: value }
  postConfigToIframe()
}

function setPointField(name: string, value: unknown) {
  pointMetadataDraft.value = { ...pointMetadataDraft.value, [name]: value }
}

function openContextMenu(event: MouseEvent, target: ContextTarget) {
  event.preventDefault()
  contextTarget.value = target
  contextX.value = event.clientX
  contextY.value = event.clientY
  contextVisible.value = true
}

function handleContextSelect(key: string) {
  contextVisible.value = false
  if (key === 'add-group') {
    openGroupModal()
  } else if (key === 'add-device') {
    const groupId = contextTarget.value.type === 'group' ? contextTarget.value.group.id : ''
    openDeviceModal(groupId)
  } else if (key === 'delete-device' && contextTarget.value.type === 'device') {
    confirmDeleteDevice(contextTarget.value.device)
  } else if (key === 'delete-group' && contextTarget.value.type === 'group') {
    confirmDeleteGroup(contextTarget.value.group)
  }
}

function openGroupModal() {
  groupForm.name = ''
  showGroupModal.value = true
}

async function createGroup() {
  await consoleStore.createDeviceGroup({ name: groupForm.name })
  showGroupModal.value = false
  message.success('设备分组已创建')
}

function openDeviceModal(groupId = '') {
  deviceForm.name = ''
  deviceForm.groupId = groupId || consoleStore.deviceGroups[0]?.id || ''
  deviceForm.pluginId = pluginOptions.value[0]?.value ?? ''
  deviceForm.description = ''
  showDeviceModal.value = true
}

async function createDevice() {
  await consoleStore.createDevice({ ...deviceForm })
  activeTab.value = 'config'
  showDeviceModal.value = false
  message.success('设备已创建')
}

function confirmDeleteDevice(device: DeviceItem) {
  dialog.warning({
    title: '删除设备',
    content: `确认删除设备 ${device.name}？设备配置、点位和采集任务会一并删除。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await consoleStore.deleteDevice(device.id)
      message.success('设备已删除')
    },
  })
}

function confirmDeleteGroup(group: DeviceGroup) {
  dialog.warning({
    title: '删除分组',
    content: `确认删除分组 ${group.name}？分组下存在设备时后端会拒绝删除。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await consoleStore.deleteDeviceGroup(group.id)
      message.success('设备分组已删除')
    },
  })
}

function selectDevice(deviceId: string) {
  activeTab.value = 'config'
  void consoleStore.selectDevice(deviceId)
}

async function saveConfig() {
  await consoleStore.saveDevicePluginConfig(configDraft.value)
  message.success('设备插件配置已保存')
}

async function testConnection() {
  const result = await consoleStore.testDevicePluginConnection(configDraft.value)
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

function openPointModal(index = -1) {
  editingPointIndex.value = index
  const point = index >= 0 ? pointDraft.value[index] : null
  pointForm.name = point?.name ?? ''
  pointForm.description = point?.description ?? ''
  pointForm.address = point?.address ?? ''
  pointForm.valueType = point?.valueType ?? 'float'
  pointForm.unit = point?.unit ?? ''
  pointForm.enabled = point?.enabled ?? true
  const defaults = buildSchemaDefaults(pointSchema.value)
  pointMetadataDraft.value = { ...defaults, ...(point?.metadata ?? {}) }
  showPointModal.value = true
}

function savePointDraft() {
  const metadata = compactMetadata(pointMetadataDraft.value)
  const current = editingPointIndex.value >= 0 ? pointDraft.value[editingPointIndex.value] : null
  const point: PointItem = {
    id: current?.id ?? '',
    deviceId: selectedDevice.value?.id ?? '',
    pluginId: selectedDevice.value?.pluginId ?? '',
    name: pointForm.name,
    description: pointForm.description,
    address: pointForm.address,
    valueType: pointForm.valueType,
    unit: pointForm.unit,
    enabled: pointForm.enabled,
    tags: current?.tags ?? {},
    metadata,
  }
  if (editingPointIndex.value >= 0) {
    pointDraft.value.splice(editingPointIndex.value, 1, point)
  } else {
    pointDraft.value.push(point)
  }
  showPointModal.value = false
}

function applyPointAddressExample(value: string) {
  pointForm.address = value
}

function removePoint(index: number) {
  pointDraft.value.splice(index, 1)
}

async function savePoints() {
  await consoleStore.saveDevicePoints(pointDraft.value)
  message.success('点位表已保存')
}

function compactMetadata(values: Record<string, unknown>) {
  const metadata: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(values ?? {})) {
    if (value === null || value === undefined) {
      continue
    }
    if (typeof value === 'string' && value.trim() === '') {
      continue
    }
    metadata[key] = value
  }
  return metadata
}

function postConfigToIframe() {
  const currentDevice = configPage.value?.device ?? selectedDevice.value
  const deviceId = currentDevice?.id
  const runtimeApiBase = consoleStore.overview?.runtime.apiBase?.trim()
  iframeRef.value?.contentWindow?.postMessage({
    type: 'gcoll:init',
    payload: {
      apiBase: runtimeApiBase ? `${runtimeApiBase.replace(/\/+$/, '')}/api/v1` : getApiBasePath(),
      apis: toMessageData(buildPluginConfigApis(deviceId)),
      language: getCurrentLanguage(),
      config: toMessageData(configDraft.value),
      schema: toMessageData(configPage.value?.configSchema ?? {}),
      device: toMessageData(currentDevice),
      plugin: toMessageData(configPage.value?.plugin),
      readonly: false,
    },
  }, '*')
}

function handlePluginMessage(event: MessageEvent) {
  const data = event.data as { type?: string; payload?: unknown }
  if (!data?.type) return
  if (data.type === 'gcoll:page-ready' || data.type === 'gcoll:request-init') {
    postConfigToIframe()
  } else if (data.type === 'gcoll:config-change' && data.payload && typeof data.payload === 'object') {
    configDraft.value = { ...(data.payload as Record<string, unknown>) }
  } else if (data.type === 'gcoll:config-saved') {
    void consoleStore.loadSelectedDeviceDetails()
  } else if (data.type === 'gcoll:test-finished') {
    void consoleStore.loadDevicePluginConfigPage()
  }
}
</script>

<template>
  <div class="device-page">
    <PageHeader title="设备管理" description="左侧管理设备分组和设备，右侧维护设备配置与通用点位表。">
      <template #actions>
        <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadDevices">
          <template #icon><RefreshCw :size="16" /></template>
          刷新
        </NButton>
        <NButton secondary @click="openGroupModal">
          <template #icon><FolderPlus :size="16" /></template>
          添加分组
        </NButton>
        <NButton type="primary" :disabled="!pluginOptions.length || !groupOptions.length" @click="openDeviceModal()">
          <template #icon><CirclePlus :size="16" /></template>
          添加设备
        </NButton>
      </template>
    </PageHeader>

    <div class="device-console">
      <aside class="device-browser" @contextmenu="openContextMenu($event, { type: 'blank' })">
        <div class="panel-titlebar">
          <strong>设备菜单</strong>
          <div class="mini-actions">
            <NButton text size="small" @click.stop="openGroupModal"><FolderPlus :size="15" /></NButton>
            <NButton text size="small" :disabled="!pluginOptions.length || !groupOptions.length" @click.stop="openDeviceModal()"><FilePlus :size="15" /></NButton>
          </div>
        </div>
        <NInput v-model:value="searchText" class="device-search" size="small" clearable placeholder="搜索设备名称或插件" />
        <div class="device-groups">
          <div v-for="group in filteredGroups" :key="group.id" class="device-group" @contextmenu.stop="openContextMenu($event, { type: 'group', group })">
            <button class="group-title group-title--button" type="button">
              <span class="group-title__main">
                <ChevronDown :size="14" />
                {{ group.name }}
              </span>
              <span>{{ group.devices.length }}</span>
            </button>
            <div class="device-list">
              <button
                v-for="device in group.devices"
                :key="device.id"
                type="button"
                class="device-card"
                :class="{ 'device-card--active': device.id === consoleStore.selectedDeviceId }"
                @click="selectDevice(device.id)"
                @contextmenu.stop="openContextMenu($event, { type: 'device', device })"
              >
                <span class="device-dot" :class="`device-dot--${device.status}`" />
                <span class="device-card__main">
                  <strong>{{ device.name }}</strong>
                  <small>{{ device.pluginName }}</small>
                </span>
                <span class="device-card__status" :class="{ 'is-online': device.status === 'online' }">{{ statusLabel(device.status) }}</span>
              </button>
              <NEmpty v-if="!group.devices.length" class="device-empty" size="small" description="没有设备" />
            </div>
          </div>
        </div>
        <div class="device-count">共 {{ consoleStore.devices.length }} 个设备</div>
      </aside>

      <section class="device-detail-panel">
        <div v-if="selectedDevice" class="detail-titlebar detail-titlebar--header">
          <div class="detail-heading">
            <strong>{{ selectedDevice.name }}</strong>
            <small>{{ selectedDevice.pluginName }}</small>
          </div>
          <StatusBadge :label="statusLabel(selectedDevice.status)" :kind="deviceStatusKind(selectedDevice.status)" />
        </div>

        <div v-if="!selectedDevice" class="device-empty-shell">
          <NEmpty description="请选择左侧设备" />
        </div>

        <div v-else class="detail-tabs">
          <NTabs v-model:value="activeTab" class="detail-tabs__inner" type="line" animated>
            <NTabPane name="config" tab="设备配置">
              <div class="detail-body">
                <div class="config-form-shell">
                  <NAlert v-if="configPage && !configPage.configured" type="warning" :bordered="false">
                    当前设备还没有已保存的插件运行配置，保存配置后才能维护点位和启动采集。
                  </NAlert>

                  <div v-if="useCustomConfigPage" class="plugin-page-frame">
                    <iframe ref="iframeRef" title="插件设备配置页" sandbox="allow-scripts allow-forms allow-same-origin" :srcdoc="configSrcdoc" @load="postConfigToIframe" />
                  </div>

                  <NForm v-else class="device-protocol-form" label-placement="top">
                    <NFormItem v-for="field in configFields" :key="field.name" :label="field.title" :feedback="field.description">
                      <NInputNumber
                        v-if="field.type === 'number'"
                        class="full-input"
                        :value="fieldNumber(field.name, field.defaultValue)"
                        @update:value="(value) => setConfigField(field.name, value)"
                      />
                      <NSwitch
                        v-else-if="field.component === 'switch' || field.type === 'boolean'"
                        :value="fieldBool(field.name, field.defaultValue)"
                        @update:value="(value) => setConfigField(field.name, value)"
                      />
                      <NSelect
                        v-else-if="field.component === 'select' || field.options.length"
                        :value="fieldString(field.name, field.defaultValue)"
                        :options="field.options"
                        @update:value="(value) => setConfigField(field.name, value)"
                      />
                      <NInput
                        v-else
                        :type="field.component === 'textarea' ? 'textarea' : field.component === 'password' ? 'password' : 'text'"
                        :value="fieldString(field.name, field.defaultValue)"
                        @update:value="(value) => setConfigField(field.name, value)"
                      />
                    </NFormItem>
                  </NForm>

                  <div v-if="!useCustomConfigPage" class="form-actions">
                    <NButton secondary :loading="consoleStore.loading" @click="testConnection">
                      <template #icon><PlugZap :size="14" /></template>
                      测试连接
                    </NButton>
                    <NButton type="primary" :loading="consoleStore.loading" @click="saveConfig">
                      <template #icon><Save :size="14" /></template>
                      保存配置
                    </NButton>
                  </div>
                </div>

                <div class="test-panel">
                  <div class="test-header">
                    <strong>插件与运行状态</strong>
                    <NTag size="small" type="info" bordered>{{ configPage?.plugin.version }}</NTag>
                  </div>
                  <div class="test-meta">
                    <div><dt>插件 ID</dt><dd>{{ selectedDevice.pluginId }}</dd></div>
                    <div><dt>运行状态</dt><dd>{{ configPage?.plugin.status ?? '-' }}</dd></div>
                    <div><dt>权限</dt><dd>{{ configPage?.plugin.permissions.join(', ') || '-' }}</dd></div>
                    <div><dt>点位数量</dt><dd>{{ pointDraft.length }}</dd></div>
                  </div>
                  <NAlert class="device-plugin-alert" type="info" :bordered="false">
                    {{ useCustomConfigPage ? '插件自定义页面会直接调用主服务接口保存和测试，宿主只负责状态同步与结果展示。' : '系统简单配置页由宿主统一渲染字段、保存配置并测试连接。' }}
                  </NAlert>
                  <NDataTable class="events-table" size="small" :columns="eventColumns" :data="configPage?.recentEvents ?? []" :bordered="false" />
                </div>
              </div>
            </NTabPane>

            <NTabPane name="points" tab="点位表">
              <div class="points-panel points-panel--tab">
                <div class="points-titlebar">
                  <div class="points-title">
                    <strong>通用点位表</strong>
                    <NTag size="small" :type="configPage?.configured ? 'success' : 'warning'" bordered>
                      {{ configPage?.configured ? '配置已保存' : '请先保存设备配置' }}
                    </NTag>
                  </div>
                  <div class="points-actions">
                    <NButton size="small" secondary :disabled="!configPage?.configured || useCustomPointPage" @click="openPointModal()">
                      <template #icon><CirclePlus :size="14" /></template>
                      新增点位
                    </NButton>
                    <NButton size="small" type="primary" :disabled="!configPage?.configured || useCustomPointPage" :loading="consoleStore.loading" @click="savePoints">
                      <template #icon><Save :size="14" /></template>
                      保存点位
                    </NButton>
                    <NButton size="small" type="primary" secondary :disabled="!configPage?.configured || !pointDraft.length" @click="startTask">
                      <template #icon><Play :size="14" /></template>
                      启动采集
                    </NButton>
                  </div>
                </div>

                <div v-if="useCustomPointPage" class="plugin-page-frame plugin-page-frame--points">
                  <iframe title="插件点位配置页" sandbox="allow-scripts allow-forms" :srcdoc="pointSrcdoc" />
                </div>
                <NDataTable
                  v-else
                  class="points-table"
                  size="small"
                  :columns="pointColumns"
                  :data="pointDraft"
                  :bordered="false"
                  :single-line="false"
                  flex-height
                />
              </div>
            </NTabPane>
          </NTabs>
        </div>
      </section>
    </div>

    <NDropdown
      placement="bottom-start"
      trigger="manual"
      :x="contextX"
      :y="contextY"
      :show="contextVisible"
      :options="contextOptions"
      @clickoutside="contextVisible = false"
      @select="handleContextSelect"
    />

    <NModal v-model:show="showGroupModal" preset="dialog" title="添加设备分组" positive-text="创建" negative-text="取消" @positive-click="createGroup">
      <NForm label-placement="top">
        <NFormItem label="分组名称"><NInput v-model:value="groupForm.name" /></NFormItem>
      </NForm>
    </NModal>

    <NModal v-model:show="showDeviceModal" preset="dialog" title="添加设备" positive-text="创建" negative-text="取消" @positive-click="createDevice">
      <NForm label-placement="top">
        <NFormItem label="设备名称"><NInput v-model:value="deviceForm.name" /></NFormItem>
        <NFormItem label="设备分组"><NSelect v-model:value="deviceForm.groupId" :options="groupOptions" /></NFormItem>
        <NFormItem label="南向插件"><NSelect v-model:value="deviceForm.pluginId" :options="pluginOptions" /></NFormItem>
        <NFormItem label="说明"><NInput v-model:value="deviceForm.description" type="textarea" /></NFormItem>
      </NForm>
    </NModal>

    <NModal v-model:show="showPointModal" preset="dialog" :title="editingPointIndex >= 0 ? '编辑点位' : '新增点位'" positive-text="保存" negative-text="取消" @positive-click="savePointDraft">
      <NForm label-placement="top">
        <NFormItem label="点位名称"><NInput v-model:value="pointForm.name" /></NFormItem>
        <NFormItem :feedback="pointAddressDescription">
          <template #label>
            <span>地址</span>
          </template>
          <NInput v-model:value="pointForm.address">
            <template v-if="pointAddressExamples.length" #suffix>
              <NPopover trigger="click" placement="bottom-end" :show-arrow="false">
                <template #trigger>
                  <button type="button" class="point-example-trigger">点位示例</button>
                </template>
                <div class="point-example-popover">
                  <div class="point-example-title">设备支持的地址示例</div>
                  <div v-for="group in pointAddressExampleGroups" :key="group.label" class="point-example-group">
                    <strong class="point-example-group__title">{{ group.label }}</strong>
                    <div class="point-example-list">
                      <NTooltip v-for="example in group.items" :key="example.value" trigger="hover">
                        <template #trigger>
                          <button type="button" class="point-example-chip" @click="applyPointAddressExample(example.value)">
                            {{ example.value }}
                          </button>
                        </template>
                        {{ example.description || example.label }}
                      </NTooltip>
                    </div>
                  </div>
                </div>
              </NPopover>
            </template>
          </NInput>
        </NFormItem>
        <NFormItem label="值类型">
          <NSelect
            v-model:value="pointForm.valueType"
            :options="['bool', 'int', 'float', 'string', 'bytes', 'datetime', 'json'].map((value) => ({ label: value, value }))"
          />
        </NFormItem>
        <NFormItem label="单位"><NInput v-model:value="pointForm.unit" /></NFormItem>
        <NFormItem label="说明"><NInput v-model:value="pointForm.description" /></NFormItem>
        <NFormItem v-for="field in pointFields" :key="field.name" :label="field.title" :feedback="field.description">
          <NInputNumber
            v-if="field.type === 'number'"
            class="full-input"
            :value="pointFieldNumber(field.name, field.defaultValue)"
            @update:value="(value) => setPointField(field.name, value)"
          />
          <NSwitch
            v-else-if="field.component === 'switch' || field.type === 'boolean'"
            :value="pointFieldBool(field.name, field.defaultValue)"
            @update:value="(value) => setPointField(field.name, value)"
          />
          <NSelect
            v-else-if="field.component === 'select' || field.options.length"
            :value="pointFieldString(field.name, field.defaultValue)"
            :options="field.options"
            @update:value="(value) => setPointField(field.name, value)"
          />
          <NInput
            v-else
            :type="field.component === 'textarea' ? 'textarea' : field.component === 'password' ? 'password' : 'text'"
            :value="pointFieldString(field.name, field.defaultValue)"
            @update:value="(value) => setPointField(field.name, value)"
          />
        </NFormItem>
        <NFormItem label="启用"><NSwitch v-model:value="pointForm.enabled" /></NFormItem>
      </NForm>
    </NModal>
  </div>
</template>

<style scoped>
.point-example-popover {
  width: 360px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.point-example-title {
  font-size: 13px;
  color: inherit;
  opacity: 0.72;
}

.point-example-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.point-example-group__title {
  font-size: 12px;
  color: inherit;
}

.point-example-trigger {
  border: none;
  background: transparent;
  color: #0284c7;
  font-size: 12px;
  cursor: pointer;
  padding: 0;
}

.point-example-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.point-example-chip {
  border: 1px solid rgba(148, 163, 184, 0.35);
  background: transparent;
  color: inherit;
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
  line-height: 18px;
  cursor: pointer;
}

.point-example-chip:hover {
  border-color: #0ea5e9;
  color: #0284c7;
  background: rgba(14, 165, 233, 0.12);
}
</style>
