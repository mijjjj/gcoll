import { request } from './http'

export interface MetricItem {
  key: string
  label: string
  value: string
  hint: string
  tone: 'primary' | 'success' | 'warning' | 'error' | 'neutral'
}

export interface RuntimeStatus {
  status: string
  service: string
  version: string
  mode: string
  checkedAt: string
  apiBase: string
  database: string
}

export interface DataPlaneStatus {
  queueUsagePercent: number
  ruleHitPercent: number
  forwardPercent: number
  throughput: string
  latency: string
  backpressure: string
}

export interface PluginSummary {
  running: number
  total: number
}

export interface RuntimeDependency {
  name: string
  status: string
  detail: string
}

export interface RuntimeEvent {
  id: string
  time: string
  level: string
  source: string
  pluginId: string
  deviceId: string
  taskId: string
  message: string
  traceId: string
}

export interface PluginItem {
  id: string
  name: string
  type: string
  version: string
  runtime: string
  protocol: string
  status: string
  permissions: string[]
  updatedAt: string
}

export interface PluginOperation {
  key: string
  label: string
  description: string
  enabled: boolean
}

export interface DevicePluginConfigPage {
  plugin: PluginItem
  device: DeviceItem
  config: Record<string, unknown>
  configSchema: Record<string, unknown>
  customConfigPage: PluginAssetPage
  customPointPage: PluginAssetPage
  configured: boolean
  points: PointItem[]
  recentEvents: RuntimeEvent[]
  operations: PluginOperation[]
  warnings: string[]
}

export interface PluginAssetPage {
  enabled: boolean
  entry: string
  script: string
  html: string
  js: string
}

export interface TestConnectionResult {
  success: boolean
  message: string
  latencyMs: number
  traceId: string
}

export interface DeviceGroup {
  id: string
  name: string
  count: number
}

export interface DeviceItem {
  id: string
  name: string
  groupId: string
  pluginId: string
  pluginName: string
  status: string
  enabled: boolean
  pointCount: number
  reportMode: string
  lastSeenAt: string
  description: string
}

export interface PointItem {
  id: string
  deviceId: string
  pluginId: string
  name: string
  description: string
  address: string
  valueType: string
  unit: string
  enabled: boolean
  tags: Record<string, string>
  metadata?: Record<string, unknown>
}

export interface TaskSummary {
  id: string
  name: string
  deviceId: string
  deviceName: string
  southPluginName: string
  pointCount: number
  reportMode: string
  status: string
  rate: string
  ruleHitRate: string
  lastCollectedAt: string
}

export interface PointCacheItem {
  pointId: string
  deviceId: string
  pointName: string
  value: string
  quality: string
  changed: boolean
  updatedAt: string
}

export interface PipelineRuleItem {
  id: string
  name: string
  enabled: boolean
  expression: string
  matched: number
  targetCount: number
  updatedAt: string
}

export interface ForwardTargetItem {
  id: string
  name: string
  pluginName: string
  status: string
  endpoint: string
  lastError: string
  updatedAt: string
}

export interface OverviewData {
  metrics: MetricItem[]
  runtime: RuntimeStatus
  dataPlane: DataPlaneStatus
  tasks: TaskSummary[]
  recentEvents: RuntimeEvent[]
  pluginSummary: PluginSummary
  network: RuntimeDependency
}

export const consoleApi = {
  getOverview: () => request<OverviewData>('/runtime/overview'),
  getPlugins: () => request<{ items: PluginItem[] }>('/plugins'),
  createDeviceGroup: (payload: { id?: string; name: string }) =>
    request<{ group: DeviceGroup }>('/device-groups', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  deleteDeviceGroup: (groupId: string) =>
    request<{ groupId: string }>(`/device-groups/${groupId}`, {
      method: 'DELETE',
    }),
  createDevice: (payload: Partial<DeviceItem> & { groupId: string; pluginId: string; reportMode: string; config?: Record<string, unknown> }) =>
    request<{ device: DeviceItem }>('/devices', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  getDevices: () => request<{ groups: DeviceGroup[]; items: DeviceItem[] }>('/devices'),
  moveDeviceToGroup: (deviceId: string, groupId: string) =>
    request<{ device: DeviceItem }>(`/devices/${deviceId}/group`, {
      method: 'PATCH',
      body: JSON.stringify({ groupId }),
    }),
  deleteDevice: (deviceId: string) =>
    request<{ deviceId: string }>(`/devices/${deviceId}`, {
      method: 'DELETE',
    }),
  getDevicePluginConfigPage: (deviceId: string) => request<DevicePluginConfigPage>(`/devices/${deviceId}/protocol-config`),
  updateDevicePluginConfig: (deviceId: string, config: Record<string, unknown>) =>
    request<{ config: Record<string, unknown> }>(`/devices/${deviceId}/protocol-config`, {
      method: 'PUT',
      body: JSON.stringify({ config }),
    }),
  testDevicePluginConnection: (deviceId: string, config?: Record<string, unknown>) =>
    request<TestConnectionResult>(`/devices/${deviceId}/protocol-config/test`, {
      method: 'POST',
      body: JSON.stringify({ config }),
    }),
  getDevicePoints: (deviceId: string) => request<{ items: PointItem[] }>(`/devices/${deviceId}/points`),
  createDevicePoint: (deviceId: string, payload: Partial<PointItem> & { pluginId: string; name: string; address: string; valueType: string; metadata?: Record<string, unknown> }) =>
    request<{ point: PointItem }>(`/devices/${deviceId}/points`, {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  updateDevicePoints: (deviceId: string, items: PointItem[]) =>
    request<{ items: PointItem[] }>(`/devices/${deviceId}/points`, {
      method: 'PUT',
      body: JSON.stringify({ items }),
    }),
  getTasks: () => request<{ items: TaskSummary[] }>('/tasks'),
  startDeviceTask: (deviceId: string) =>
    request<{ task: TaskSummary }>(`/devices/${deviceId}/tasks/start`, {
      method: 'POST',
    }),
  startTask: (taskId: string) =>
    request<{ task: TaskSummary }>(`/tasks/${taskId}/start`, {
      method: 'POST',
    }),
  stopTask: (taskId: string) =>
    request<{ task: TaskSummary }>(`/tasks/${taskId}/stop`, {
      method: 'POST',
    }),
  getPointCache: () => request<{ items: PointCacheItem[] }>('/point-cache'),
  getPipelineRules: () => request<{ items: PipelineRuleItem[] }>('/pipeline/rules'),
  getTargets: () => request<{ items: ForwardTargetItem[] }>('/targets'),
  getLogs: () => request<{ items: RuntimeEvent[] }>('/logs'),
}
