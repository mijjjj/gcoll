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
  NSelect,
  NSwitch,
  NTabPane,
  NTabs,
  NTag,
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

interface ConfigField {
  name: string
  title: string
  description: string
  type: string
  component: string
  defaultValue: unknown
  options: Array<{ label: string; value: string }>
}

type ContextTarget = { type: 'blank' } | { type: 'group'; group: DeviceGroup } | { type: 'device'; device: DeviceItem }

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
  enabled: true,
  metadataText: '{}',
})

const selectedDevice = computed(() => consoleStore.selectedDevice)
const configPage = computed(() => consoleStore.devicePluginConfigPage)
const customConfigPage = computed(() => configPage.value?.customConfigPage)
const customPointPage = computed(() => configPage.value?.customPointPage)
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
      return `${device.name} ${device.code} ${device.pluginName}`.toLowerCase().includes(keyword)
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
  const properties = (schema.properties ?? {}) as Record<string, Record<string, unknown>>
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
})

const configSrcdoc = computed(() => {
  const page = customConfigPage.value
  if (!page?.html) return ''
  const script = page.js ? `<script>${page.js}<\/script>` : ''
  return `${page.html}\n${script}`
})

const pointSrcdoc = computed(() => {
  const page = customPointPage.value
  if (!page?.html) return ''
  const script = page.js ? `<script>${page.js}<\/script>` : ''
  return `${page.html}\n${script}`
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
    title: 'metadata',
    key: 'metadata',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render: (row) => JSON.stringify(row.metadata ?? {}),
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

const eventColumns: DataTableColumns<RuntimeEvent> = [
  { title: '时间', key: 'time', minWidth: 160 },
  { title: '级别', key: 'level', width: 80 },
  { title: '来源', key: 'source', width: 100 },
  { title: '消息', key: 'message', minWidth: 240, ellipsis: { tooltip: true } },
]

watch(
  () => configPage.value?.config,
  (config) => {
    const draft: Record<string, unknown> = { ...(config ?? {}) }
    for (const field of configFields.value) {
      if (draft[field.name] === undefined && field.defaultValue !== undefined) {
        draft[field.name] = field.defaultValue
      }
    }
    configDraft.value = draft
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

function setConfigField(name: string, value: unknown) {
  configDraft.value = { ...configDraft.value, [name]: value }
  postConfigToIframe()
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
  deviceForm.code = ''
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

function openPointModal(index = -1) {
  editingPointIndex.value = index
  const point = index >= 0 ? pointDraft.value[index] : null
  pointForm.name = point?.name ?? ''
  pointForm.description = point?.description ?? ''
  pointForm.address = point?.address ?? ''
  pointForm.valueType = point?.valueType ?? 'float'
  pointForm.unit = point?.unit ?? ''
  pointForm.enabled = point?.enabled ?? true
  pointForm.metadataText = JSON.stringify(point?.metadata ?? {}, null, 2)
  showPointModal.value = true
}

function savePointDraft() {
  const metadata = parseJson(pointForm.metadataText)
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

function removePoint(index: number) {
  pointDraft.value.splice(index, 1)
}

async function savePoints() {
  await consoleStore.saveDevicePoints(pointDraft.value)
  message.success('点位表已保存')
}

function parseJson(value: string): Record<string, unknown> {
  const parsed = JSON.parse(value || '{}')
  if (!parsed || Array.isArray(parsed) || typeof parsed !== 'object') {
    throw new Error('JSON 必须是对象')
  }
  return parsed as Record<string, unknown>
}

function postConfigToIframe() {
  iframeRef.value?.contentWindow?.postMessage({
    type: 'gcoll:init',
    payload: {
      config: configDraft.value,
      schema: configPage.value?.configSchema ?? {},
      device: selectedDevice.value,
      plugin: configPage.value?.plugin,
      readonly: false,
    },
  }, '*')
}

function handlePluginMessage(event: MessageEvent) {
  const data = event.data as { type?: string; payload?: unknown }
  if (!data?.type) return
  if (data.type === 'gcoll:config-change' && data.payload && typeof data.payload === 'object') {
    configDraft.value = { ...(data.payload as Record<string, unknown>) }
  } else if (data.type === 'gcoll:save-config') {
    void saveConfig()
  } else if (data.type === 'gcoll:test-connection') {
    void testConnection()
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
        <NInput v-model:value="searchText" class="device-search" size="small" clearable placeholder="搜索设备名称、编码或插件" />
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
                  <small>{{ device.code }} / {{ device.pluginName }}</small>
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
            <small>{{ selectedDevice.code }} / {{ selectedDevice.pluginName }}</small>
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
                    <iframe ref="iframeRef" title="插件设备配置页" sandbox="allow-scripts allow-forms" :srcdoc="configSrcdoc" @load="postConfigToIframe" />
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

                  <div class="form-actions">
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
                    插件自定义页面只能通过宿主消息通道保存和测试，最终配置仍保存在主程序数据库。
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
        <NFormItem label="设备编码"><NInput v-model:value="deviceForm.code" /></NFormItem>
        <NFormItem label="设备分组"><NSelect v-model:value="deviceForm.groupId" :options="groupOptions" /></NFormItem>
        <NFormItem label="南向插件"><NSelect v-model:value="deviceForm.pluginId" :options="pluginOptions" /></NFormItem>
        <NFormItem label="说明"><NInput v-model:value="deviceForm.description" type="textarea" /></NFormItem>
      </NForm>
    </NModal>

    <NModal v-model:show="showPointModal" preset="dialog" :title="editingPointIndex >= 0 ? '编辑点位' : '新增点位'" positive-text="保存" negative-text="取消" @positive-click="savePointDraft">
      <NForm label-placement="top">
        <NFormItem label="点位名称"><NInput v-model:value="pointForm.name" /></NFormItem>
        <NFormItem label="地址"><NInput v-model:value="pointForm.address" /></NFormItem>
        <NFormItem label="值类型">
          <NSelect
            v-model:value="pointForm.valueType"
            :options="['bool', 'int', 'float', 'string', 'bytes', 'datetime', 'json'].map((value) => ({ label: value, value }))"
          />
        </NFormItem>
        <NFormItem label="单位"><NInput v-model:value="pointForm.unit" /></NFormItem>
        <NFormItem label="说明"><NInput v-model:value="pointForm.description" /></NFormItem>
        <NFormItem label="metadata JSON"><NInput v-model:value="pointForm.metadataText" type="textarea" :autosize="{ minRows: 4, maxRows: 8 }" /></NFormItem>
        <NFormItem label="启用"><NSwitch v-model:value="pointForm.enabled" /></NFormItem>
      </NForm>
    </NModal>
  </div>
</template>
