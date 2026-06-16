import { getCurrentLanguage } from '../i18n'

interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

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

export interface ModbusTcpDeviceConfig {
  host: string
  port: number
  unitId: number
  timeoutMs: number
  pollIntervalMs: number
  reportMode: string
  debugEnabled: boolean
  maxCoilBatch: number
  maxRegisterBatch: number
  lowLatencyMs: number
  highLatencyMs: number
}

export interface ModbusTcpReadBlock {
  area: string
  start: number
  quantity: number
  pointIds: string[]
  latencyMs: number
}

export interface ModbusTcpPoint {
  id: string
  name: string
  area: string
  address: number
  quantity: number
  valueType: string
  mode: string
  reportMode: string
  enabled: boolean
  byteOrder: string
  wordOrder: string
  scale: string
  current: string
  quality: string
  lastReadAt: string
  description: string
}

export interface ModbusTcpDebugLog {
  time: string
  level: string
  message: string
  traceId: string
  area: string
  address: string
  costMs: number
  rawHex: string
}

export interface ModbusTcpOperation {
  key: string
  label: string
  description: string
  enabled: boolean
}

export interface ModbusTcpDeviceConfigPage {
  plugin: PluginItem
  device: DeviceItem
  config: ModbusTcpDeviceConfig
  readPlan: ModbusTcpReadBlock[]
  points: ModbusTcpPoint[]
  debugLogs: ModbusTcpDebugLog[]
  operations: ModbusTcpOperation[]
  warnings: string[]
}

export interface DeviceGroup {
  id: string
  name: string
  count: number
}

export interface DeviceItem {
  id: string
  name: string
  code: string
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

async function request<T>(path: string): Promise<T> {
  const language = getCurrentLanguage()
  const searchParams = new URLSearchParams({ lang: language })
  const response = await fetch(`/api/v1${path}?${searchParams.toString()}`, {
    headers: {
      'Accept-Language': language,
    },
  })
  if (!response.ok) {
    throw new Error(`请求失败: ${response.status}`)
  }
  const result = (await response.json()) as ApiResponse<T>
  if (result.code !== 0) {
    throw new Error(result.message || '请求失败')
  }
  return result.data
}

export const consoleApi = {
  getOverview: () => request<OverviewData>('/runtime/overview'),
  getPlugins: () => request<{ items: PluginItem[] }>('/plugins'),
  getDevices: () => request<{ groups: DeviceGroup[]; items: DeviceItem[] }>('/devices'),
  getModbusTcpDeviceConfigPage: (deviceId: string) => request<ModbusTcpDeviceConfigPage>(`/devices/${deviceId}/protocol-config`),
  getDevicePoints: (deviceId: string) => request<{ items: PointItem[] }>(`/devices/${deviceId}/points`),
  getTasks: () => request<{ items: TaskSummary[] }>('/tasks'),
  getPointCache: () => request<{ items: PointCacheItem[] }>('/point-cache'),
  getPipelineRules: () => request<{ items: PipelineRuleItem[] }>('/pipeline/rules'),
  getTargets: () => request<{ items: ForwardTargetItem[] }>('/targets'),
  getLogs: () => request<{ items: RuntimeEvent[] }>('/logs'),
}
