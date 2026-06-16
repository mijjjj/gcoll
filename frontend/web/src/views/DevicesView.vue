<script setup lang="ts">
import { computed, h, onMounted, ref } from 'vue'
import { NButton, NDataTable, NDescriptions, NDescriptionsItem, NInput, NTag, type DataTableColumns } from 'naive-ui'
import { ChevronDown, Download, MoreHorizontal, Play, Plus, RefreshCw, Search, Upload } from '@lucide/vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'
import type { DeviceItem, PointItem } from '../api/console'

const consoleStore = useConsoleStore()
const keyword = ref('')

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
            <span>设备详情 - {{ consoleStore.selectedDevice?.name ?? '待选择' }}</span>
          </div>
        </div>

        <div v-if="consoleStore.selectedDevice" class="detail-body detail-body--single">
          <section class="config-form-shell">
            <NDescriptions label-placement="left" :column="2" size="small" bordered>
              <NDescriptionsItem label="设备编号">{{ (consoleStore.selectedDevice as DeviceItem).code }}</NDescriptionsItem>
              <NDescriptionsItem label="状态">
                <StatusBadge :label="statusLabel((consoleStore.selectedDevice as DeviceItem).status)" :kind="statusKind((consoleStore.selectedDevice as DeviceItem).status)" />
              </NDescriptionsItem>
              <NDescriptionsItem label="南向插件">{{ (consoleStore.selectedDevice as DeviceItem).pluginName }}</NDescriptionsItem>
              <NDescriptionsItem label="上报模式">
                {{ (consoleStore.selectedDevice as DeviceItem).reportMode === 'change' ? '变化上报' : '立即上报' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="点位数量">{{ (consoleStore.selectedDevice as DeviceItem).pointCount }}</NDescriptionsItem>
              <NDescriptionsItem label="最后在线">{{ (consoleStore.selectedDevice as DeviceItem).lastSeenAt }}</NDescriptionsItem>
              <NDescriptionsItem label="描述">{{ (consoleStore.selectedDevice as DeviceItem).description }}</NDescriptionsItem>
            </NDescriptions>
            <div class="form-actions">
              <NButton>取消</NButton>
              <NButton type="primary">保存配置</NButton>
            </div>
          </section>

          <aside class="test-panel">
            <div class="test-header">
              <strong>连接测试 / 调试</strong>
              <div class="test-actions">
                <NButton size="small" secondary type="primary">测试连接</NButton>
                <NButton size="small" secondary>调试模式</NButton>
              </div>
            </div>
            <div class="success-box">
              <StatusBadge label="插件宿主待接入" kind="warning" />
              <span>当前页面已完成设备配置主框架，真实测试动作后续接入插件宿主。</span>
            </div>
            <dl class="test-meta">
              <div><dt>插件 ID</dt><dd>{{ (consoleStore.selectedDevice as DeviceItem).pluginId }}</dd></div>
              <div><dt>设备 ID</dt><dd>{{ (consoleStore.selectedDevice as DeviceItem).id }}</dd></div>
              <div><dt>权限状态</dt><dd>{{ (consoleStore.selectedDevice as DeviceItem).enabled ? '可编辑' : '只读' }}</dd></div>
            </dl>
          </aside>
        </div>
      </section>
    </div>

    <section class="points-panel">
      <div class="points-titlebar">
        <div class="points-title">
          <strong>通用点位表（共 {{ consoleStore.points.length }} 个点）</strong>
          <div class="points-actions">
            <NButton size="small" secondary type="primary">
              <template #icon><Upload :size="14" /></template>
              导入点位
            </NButton>
            <NButton size="small" secondary type="primary">
              <template #icon><Download :size="14" /></template>
              导出点位
            </NButton>
            <NButton size="small" secondary :loading="consoleStore.loading" @click="consoleStore.loadPointsForSelectedDevice">
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

      <NDataTable size="small" :columns="pointColumns" :data="consoleStore.points" :bordered="false" :single-line="false" class="points-table" />
    </section>
  </div>
</template>
