import { request } from './http'

export interface RuntimeHealth {
  status: string
  service: string
  version: string
  mode: string
  checkedAt: string
}

export async function fetchRuntimeHealth(): Promise<RuntimeHealth> {
  return request<RuntimeHealth>('/runtime/health', {
    fallbackMessageKey: 'api.runtimeHealthFailed',
  })
}
