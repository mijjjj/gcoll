import { defineStore } from 'pinia'
import {
  consoleApi,
  type DeviceGroup,
  type DeviceItem,
  type ForwardTargetItem,
  type ModbusTcpDeviceConfigPage,
  type OverviewData,
  type PipelineRuleItem,
  type PluginItem,
  type PointCacheItem,
  type PointItem,
  type RuntimeEvent,
  type TaskSummary,
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
  modbusTcpDeviceConfigPage: ModbusTcpDeviceConfigPage | null
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
    modbusTcpDeviceConfigPage: null,
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
        if (!this.selectedDeviceId && result.items.length > 0) {
          this.selectedDeviceId = result.items[0].id
        }
        await this.fetchSelectedDeviceDetails()
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
    async loadModbusTcpDeviceConfigPage() {
      await this.run(async () => {
        await this.fetchModbusTcpDeviceConfigPage()
      })
    },
    async fetchSelectedDeviceDetails() {
      await this.loadPointsForSelectedDevice()
      await this.fetchModbusTcpDeviceConfigPage()
    },
    async fetchModbusTcpDeviceConfigPage() {
      const device = this.devices.find((item) => item.id === this.selectedDeviceId)
      if (!device || device.pluginId !== 'com.gcoll.modbus-tcp') {
        this.modbusTcpDeviceConfigPage = null
        return
      }
      this.modbusTcpDeviceConfigPage = await consoleApi.getModbusTcpDeviceConfigPage(device.id)
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
        this.error = error instanceof Error ? error.message : '请求失败'
      } finally {
        this.loading = false
      }
    },
  },
})
