import { getCurrentLanguage } from '../i18n'

export interface RuntimeHealth {
  status: string
  service: string
  version: string
  mode: string
  checkedAt: string
}

interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export async function fetchRuntimeHealth(): Promise<RuntimeHealth> {
  const language = getCurrentLanguage()
  const fallbackMessage = language === 'en' ? 'Runtime health check failed' : '运行时健康检查失败'
  const searchParams = new URLSearchParams({ lang: language })
  const response = await fetch(`/api/v1/runtime/health?${searchParams.toString()}`, {
    headers: {
      'Accept-Language': language,
    },
  })
  if (!response.ok) {
    throw new Error(`${fallbackMessage}: ${response.status}`)
  }
  const result = (await response.json()) as ApiResponse<RuntimeHealth>
  if (result.code !== 0) {
    throw new Error(result.message || fallbackMessage)
  }
  return result.data
}
