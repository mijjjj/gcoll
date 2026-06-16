import { defineStore } from 'pinia'
import { getCurrentLanguage } from '../i18n'
import { fetchRuntimeHealth, type RuntimeHealth } from '../api/runtime'

interface RuntimeState {
  health: RuntimeHealth | null
  loading: boolean
  error: string
}

export const useRuntimeStore = defineStore('runtime', {
  state: (): RuntimeState => ({
    health: null,
    loading: false,
    error: '',
  }),
  actions: {
    async refreshHealth() {
      this.loading = true
      this.error = ''
      try {
        this.health = await fetchRuntimeHealth()
      } catch (error) {
        const fallback = getCurrentLanguage() === 'en' ? 'Runtime health check failed' : '运行时健康检查失败'
        this.error = error instanceof Error ? error.message : fallback
      } finally {
        this.loading = false
      }
    },
  },
})
