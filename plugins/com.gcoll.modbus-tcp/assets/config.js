(() => {
  const form = document.querySelector('#modbus-config-form')
  const testButton = document.querySelector('#test-button')
  const submitButton = form?.querySelector('button[type="submit"]')
  const statusElement = document.querySelector('#status-message')

  const numberFields = new Set([
    'port',
    'unitId',
    'timeoutMs',
    'pollIntervalMs',
    'maxRetryPerRequest',
  ])

  const state = {
    apiBase: '',
    deviceId: '',
    language: document.documentElement.lang || 'zh-CN',
    initialized: false,
    apis: {
      getConfig: '',
      saveConfig: '',
      testConfig: '',
    },
  }

  let pendingInitPromise = null
  let pendingInitResolve = null

  function readForm() {
    const data = new FormData(form)
    const config = {}
    for (const [key, value] of data.entries()) {
      config[key] = numberFields.has(key) ? Number(value) : String(value)
    }
    return config
  }

  function writeForm(config) {
    for (const element of form.elements) {
      if (!element.name) continue
      const value = config[element.name]
      if (value === undefined || value === null) continue
      element.value = String(value)
    }
  }

  function mergeSchemaDefaults(schema, config) {
    const properties = schema?.properties ?? {}
    const defaults = {}
    for (const [name, property] of Object.entries(properties)) {
      if (property && Object.prototype.hasOwnProperty.call(property, 'default')) {
        defaults[name] = property.default
      }
    }
    return { ...defaults, ...(config ?? {}) }
  }

  function setStatus(message, tone = '') {
    if (!statusElement) return
    statusElement.textContent = message || ''
    if (tone) {
      statusElement.dataset.tone = tone
    } else {
      delete statusElement.dataset.tone
    }
  }

  function setPending(pending) {
    if (testButton) testButton.disabled = pending
    if (submitButton) submitButton.disabled = pending
  }

  function postConfigChange(config) {
    window.parent.postMessage({ type: 'gcoll:config-change', payload: config }, '*')
  }

  function requestHostInit() {
    window.parent.postMessage({ type: 'gcoll:request-init' }, '*')
  }

  function resolvePendingInit() {
    if (pendingInitResolve) {
      pendingInitResolve()
      pendingInitResolve = null
      pendingInitPromise = null
    }
  }

  async function waitForHostInit(timeoutMs = 800) {
    if (state.initialized && (state.deviceId || state.apis.saveConfig || state.apis.testConfig)) {
      return
    }
    if (!pendingInitPromise) {
      pendingInitPromise = new Promise((resolve) => {
        pendingInitResolve = resolve
      })
    }
    requestHostInit()
    await Promise.race([
      pendingInitPromise,
      new Promise((resolve) => window.setTimeout(resolve, timeoutMs)),
    ])
  }

  function buildFallbackPath(kind) {
    if (!state.deviceId) return ''
    const basePath = `${state.apiBase}/devices/${state.deviceId}/protocol-config`
    if (kind === 'testConfig') {
      return `${basePath}/test`
    }
    return basePath
  }

  function resolveApiPath(kind) {
    const value = state.apis?.[kind]
    if (typeof value === 'string' && value.trim() !== '') {
      return value
    }
    return buildFallbackPath(kind)
  }

  async function request(url, options = {}) {
    if (!url) {
      throw new Error('宿主尚未提供可用的配置接口地址。')
    }
    const query = new URLSearchParams({ lang: state.language || 'zh-CN' }).toString()
    const separator = url.includes('?') ? '&' : '?'
    const response = await fetch(`${url}${separator}${query}`, {
      ...options,
      headers: {
        'Accept-Language': state.language || 'zh-CN',
        'Content-Type': 'application/json',
        ...(options.headers ?? {}),
      },
    })
    const result = await response.json().catch(() => null)
    if (!response.ok || !result || result.code !== 0) {
      throw new Error(result?.message || `请求失败: ${response.status}`)
    }
    return result.data
  }

  async function saveConfig() {
    const config = readForm()
    postConfigChange(config)
    setPending(true)
    setStatus('正在保存配置...', 'warning')
    try {
      let saveUrl = resolveApiPath('saveConfig')
      if (!saveUrl) {
        await waitForHostInit()
        saveUrl = resolveApiPath('saveConfig')
      }
      await request(saveUrl, {
        method: 'PUT',
        body: JSON.stringify({ config }),
      })
      setStatus('设备插件配置已保存。', 'success')
      window.parent.postMessage({ type: 'gcoll:config-saved', payload: config }, '*')
    } catch (error) {
      setStatus(error instanceof Error ? error.message : '保存配置失败。', 'error')
    } finally {
      setPending(false)
    }
  }

  async function testConnection() {
    const config = readForm()
    postConfigChange(config)
    setPending(true)
    setStatus('正在测试连接...', 'warning')
    try {
      let testUrl = resolveApiPath('testConfig')
      if (!testUrl) {
        await waitForHostInit()
        testUrl = resolveApiPath('testConfig')
      }
      const result = await request(testUrl, {
        method: 'POST',
        body: JSON.stringify({ config }),
      })
      const latency = typeof result?.latencyMs === 'number' ? `，耗时 ${result.latencyMs} ms` : ''
      setStatus(`${result?.message ?? '连接测试完成。'}${latency}`, result?.success ? 'success' : 'warning')
    } catch (error) {
      setStatus(error instanceof Error ? error.message : '测试连接失败。', 'error')
    } finally {
      window.parent.postMessage({ type: 'gcoll:test-finished' }, '*')
      setPending(false)
    }
  }

  form.addEventListener('input', () => {
    postConfigChange(readForm())
  })

  form.addEventListener('submit', (event) => {
    event.preventDefault()
    void saveConfig()
  })

  testButton.addEventListener('click', () => {
    void testConnection()
  })

  requestHostInit()
  window.setTimeout(requestHostInit, 120)
  window.setTimeout(requestHostInit, 360)
  window.setTimeout(() => {
    if (!state.initialized || (!state.deviceId && !state.apis.saveConfig && !state.apis.testConfig)) {
      setStatus('正在等待宿主初始化配置接口地址...', 'warning')
    }
  }, 500)

  window.addEventListener('message', (event) => {
    if (event.data?.type !== 'gcoll:init') return
    const payload = event.data.payload ?? {}
    state.apiBase = payload.apiBase || ''
    state.deviceId = payload.device?.id || ''
    state.language = payload.language || document.documentElement.lang || 'zh-CN'
    state.initialized = true
    state.apis = payload.apis ?? state.apis
    const config = mergeSchemaDefaults(payload.schema, payload.config)
    writeForm(config)
    postConfigChange(config)
    setStatus('', '')
    resolvePendingInit()
  })
})()
