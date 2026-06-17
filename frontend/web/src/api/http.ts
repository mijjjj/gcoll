import { getCurrentLanguage, translate } from '../i18n'

interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

interface HttpRequestOptions extends RequestInit {
  showError?: boolean
  fallbackMessageKey?: string
}

interface HttpRequestErrorOptions {
  code?: number
  status?: number
  path: string
  cause?: unknown
}

type HttpErrorNotifier = (message: string) => void

let httpErrorNotifier: HttpErrorNotifier | null = null

export class HttpRequestError extends Error {
  code: number
  status: number
  path: string
  cause?: unknown

  constructor(message: string, options: HttpRequestErrorOptions) {
    super(message)
    this.name = 'HttpRequestError'
    this.code = options.code ?? -1
    this.status = options.status ?? 0
    this.path = options.path
    this.cause = options.cause
  }
}

export function setHttpErrorNotifier(notifier: HttpErrorNotifier | null) {
  httpErrorNotifier = notifier
}

function notifyHttpError(message: string) {
  httpErrorNotifier?.(message)
}

function buildRequestUrl(path: string, language: string) {
  const [pathname, rawQuery = ''] = path.split('?', 2)
  const searchParams = new URLSearchParams(rawQuery)
  searchParams.set('lang', language)
  return `/api/v1${pathname}?${searchParams.toString()}`
}

async function parseJsonResponse<T>(response: Response): Promise<ApiResponse<T> | null> {
  const contentType = response.headers.get('content-type') ?? ''
  if (!contentType.includes('application/json')) {
    return null
  }
  try {
    return (await response.json()) as ApiResponse<T>
  } catch {
    return null
  }
}

export async function request<T>(path: string, options: HttpRequestOptions = {}): Promise<T> {
  const {
    showError = true,
    fallbackMessageKey = 'api.requestFailed',
    headers,
    ...fetchOptions
  } = options

  const language = getCurrentLanguage()
  const requestUrl = buildRequestUrl(path, language)
  const fallbackMessage = translate(fallbackMessageKey, language)

  let response: Response
  try {
    response = await fetch(requestUrl, {
      ...fetchOptions,
      headers: {
        'Accept-Language': language,
        ...(fetchOptions.body ? { 'Content-Type': 'application/json' } : {}),
        ...headers,
      },
    })
  } catch (error) {
    const requestError = new HttpRequestError(translate('api.networkFailed', language), {
      path: requestUrl,
      cause: error,
    })
    if (showError) {
      notifyHttpError(requestError.message)
    }
    throw requestError
  }

  const result = await parseJsonResponse<T>(response)

  if (!response.ok) {
    const requestError = new HttpRequestError(result?.message || `${fallbackMessage}: ${response.status}`, {
      code: result?.code ?? response.status,
      status: response.status,
      path: requestUrl,
    })
    if (showError) {
      notifyHttpError(requestError.message)
    }
    throw requestError
  }

  if (!result) {
    const requestError = new HttpRequestError(translate('api.responseParseFailed', language), {
      status: response.status,
      path: requestUrl,
    })
    if (showError) {
      notifyHttpError(requestError.message)
    }
    throw requestError
  }

  if (result.code !== 0) {
    const requestError = new HttpRequestError(result.message || fallbackMessage, {
      code: result.code,
      status: response.status,
      path: requestUrl,
    })
    if (showError) {
      notifyHttpError(requestError.message)
    }
    throw requestError
  }

  return result.data
}
