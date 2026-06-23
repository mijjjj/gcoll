import { defineStore } from 'pinia'
import { translate } from '../i18n'
import {
  consoleApi,
  type DevicePluginConfigPage,
  type DeviceGroup,
  type DeviceItem,
  type ForwardTargetItem,
  type OverviewData,
  type PipelineRuleItem,
  type PluginItem,
  type PointCacheItem,
  type PointItem,
  type RuntimeEvent,
  type TaskSummary,
  type TestConnectionResult,
} from '../api/console'

interface ConsoleState {
  overview: OverviewData | null
  deviceGroups: DeviceGroup[]
  devices: DeviceItem[]
  selectedDeviceId: string
  points: PointItem[]
  plugins: PluginItem[]
  tasks: TaskSummary[]
  pointCache: PointCacheItem[]
  pipelineRules: PipelineRuleItem[]
  targets: ForwardTargetItem[]
  logs: RuntimeEvent[]
  devicePluginConfigPage: DevicePluginConfigPage | null
  loading: boolean
  error: string
}

export const useConsoleStore = defineStore('console', {
  state: (): ConsoleState => ({
    overview: null,
    deviceGroups: [],
    devices: [],
    selectedDeviceId: '',
    points: [],
    plugins: [],
    tasks: [],
    pointCache: [],
    pipelineRules: [],
    targets: [],
    logs: [],
    devicePluginConfigPage: null,
    loading: false,
    error: '',
  }),
  getters: {
    selectedDevice: (state) => state.devices.find((device) => device.id === state.selectedDeviceId),
  },
  actions: {
    async loadOverview() {
      await this.run(async () => {
        this.overview = await consoleApi.getOverview()
        this.tasks = this.overview.tasks
        this.logs = this.overview.recentEvents
      })
    },
    async loadDevices() {
      await this.run(async () => {
        const result = await consoleApi.getDevices()
        this.deviceGroups = result.groups
        this.devices = result.items
        const selectedExists = result.items.some((device) => device.id === this.selectedDeviceId)
        if (!selectedExists) {
          this.selectedDeviceId = result.items[0]?.id ?? ''
        }
        await this.fetchSelectedDeviceDetails()
      })
    },
    async createDeviceGroup(payload: { name: string }) {
      await this.run(async () => {
        await consoleApi.createDeviceGroup(payload)
        await this.loadDevices()
      })
    },
    async deleteDeviceGroup(groupId: string) {
      await this.run(async () => {
        await consoleApi.deleteDeviceGroup(groupId)
        await this.loadDevices()
      })
    },
    async createDevice(payload: { name: string; groupId: string; pluginId: string; description: string }) {
      await this.run(async () => {
        const result = await consoleApi.createDevice({
          ...payload,
          enabled: false,
          reportMode: 'change',
        })
        this.selectedDeviceId = result.device.id
        await this.loadDevices()
      })
    },
    async moveDeviceToGroup(deviceId: string, groupId: string) {
      await this.run(async () => {
        await consoleApi.moveDeviceToGroup(deviceId, groupId)
        this.selectedDeviceId = deviceId
        await this.loadDevices()
      })
    },
    async deleteDevice(deviceId: string) {
      await this.run(async () => {
        await consoleApi.deleteDevice(deviceId)
        if (this.selectedDeviceId === deviceId) {
          this.selectedDeviceId = ''
        }
        await this.loadDevices()
      })
    },
    async loadPointsForSelectedDevice() {
      if (!this.selectedDeviceId) {
        this.points = []
        return
      }
      const result = await consoleApi.getDevicePoints(this.selectedDeviceId)
      this.points = result.items
    },
    async selectDevice(deviceId: string) {
      this.selectedDeviceId = deviceId
      await this.run(async () => {
        await this.fetchSelectedDeviceDetails()
      })
    },
    async loadSelectedDeviceDetails() {
      await this.run(async () => {
        await this.fetchSelectedDeviceDetails()
      })
    },
    async loadPlugins() {
      await this.run(async () => {
        this.plugins = (await consoleApi.getPlugins()).items
      })
    },
    async loadDevicePluginConfigPage() {
      await this.run(async () => {
        await this.fetchDevicePluginConfigPage()
      })
    },
    async saveDevicePluginConfig(config: Record<string, unknown>) {
      if (!this.selectedDeviceId) return
      await this.run(async () => {
        await consoleApi.updateDevicePluginConfig(this.selectedDeviceId, config)
        await this.fetchSelectedDeviceDetails()
      })
    },
    async testDevicePluginConnection(config?: Record<string, unknown>): Promise<TestConnectionResult | null> {
      if (!this.selectedDeviceId) return null
      let result: TestConnectionResult | null = null
      await this.run(async () => {
        result = await consoleApi.testDevicePluginConnection(this.selectedDeviceId, config)
        await this.fetchDevicePluginConfigPage()
      })
      return result
    },
    async startSelectedDeviceTask() {
      if (!this.selectedDeviceId) return
      await this.run(async () => {
        await consoleApi.startDeviceTask(this.selectedDeviceId)
        await this.loadTasks()
        await this.fetchSelectedDeviceDetails()
      })
    },
    async createPoint(payload: { name: string; description: string; address: string; valueType: string; unit: string; metadata: Record<string, unknown> }) {
      if (!this.selectedDeviceId) return
      const device = this.devices.find((item) => item.id === this.selectedDeviceId)
      if (!device) return
      await this.run(async () => {
        await consoleApi.createDevicePoint(this.selectedDeviceId, {
          pluginId: device.pluginId,
          name: payload.name,
          description: payload.description,
          address: payload.address,
          valueType: payload.valueType,
          unit: payload.unit,
          enabled: true,
          tags: {},
          metadata: payload.metadata,
        })
        await this.fetchSelectedDeviceDetails()
      })
    },
    async saveDevicePoints(items: PointItem[]) {
      if (!this.selectedDeviceId) return
      await this.run(async () => {
        const result = await consoleApi.updateDevicePoints(this.selectedDeviceId, items)
        this.points = result.items
        await this.fetchDevicePluginConfigPage()
        await this.loadDevices()
      })
    },
    async startTask(taskId: string) {
      await this.run(async () => {
        await consoleApi.startTask(taskId)
        await this.loadTasks()
      })
    },
    async stopTask(taskId: string) {
      await this.run(async () => {
        await consoleApi.stopTask(taskId)
        await this.loadTasks()
      })
    },
    async fetchSelectedDeviceDetails() {
      if (!this.selectedDeviceId) {
        this.points = []
        this.devicePluginConfigPage = null
        return
      }
      await this.loadPointsForSelectedDevice()
      await this.fetchDevicePluginConfigPage()
    },
    async fetchDevicePluginConfigPage() {
      const device = this.devices.find((item) => item.id === this.selectedDeviceId)
      if (!device) {
        this.devicePluginConfigPage = null
        return
      }
      this.devicePluginConfigPage = await consoleApi.getDevicePluginConfigPage(device.id)
    },
    async loadTasks() {
      await this.run(async () => {
        this.tasks = (await consoleApi.getTasks()).items
      })
    },
    async loadPointCache() {
      await this.run(async () => {
        this.pointCache = (await consoleApi.getPointCache()).items
      })
    },
    async loadPipelineRules() {
      await this.run(async () => {
        this.pipelineRules = (await consoleApi.getPipelineRules()).items
      })
    },
    async loadTargets() {
      await this.run(async () => {
        this.targets = (await consoleApi.getTargets()).items
      })
    },
    async loadLogs() {
      await this.run(async () => {
        this.logs = (await consoleApi.getLogs()).items
      })
    },
    async run(action: () => Promise<void>) {
      this.loading = true
      this.error = ''
      try {
        await action()
      } catch (error) {
        this.error = error instanceof Error ? error.message : translate('api.requestFailed')
      } finally {
        this.loading = false
      }
    },
  },
})
