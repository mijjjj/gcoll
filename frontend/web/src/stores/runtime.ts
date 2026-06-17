import { defineStore } from 'pinia'
import { translate } from '../i18n'
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
        this.error = error instanceof Error ? error.message : translate('api.runtimeHealthFailed')
      } finally {
        this.loading = false
      }
    },
  },
})
