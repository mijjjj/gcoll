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
    'maxCoilBatch',
    'maxRegisterBatch',
    'debugLogLimit',
    'maxRetryPerRequest',
  ])

  const state = {
    apiBase: '/api/v1',
    deviceId: '',
    language: document.documentElement.lang || 'zh-CN',
  }

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

  async function request(path, options = {}) {
    const query = new URLSearchParams({ lang: state.language || 'zh-CN' }).toString()
    const response = await fetch(`${state.apiBase}${path}?${query}`, {
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
      await request(`/devices/${state.deviceId}/protocol-config`, {
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
      const result = await request(`/devices/${state.deviceId}/protocol-config/test`, {
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

  window.addEventListener('message', (event) => {
    if (event.data?.type !== 'gcoll:init') return
    const payload = event.data.payload ?? {}
    state.apiBase = payload.apiBase || '/api/v1'
    state.deviceId = payload.device?.id || ''
    state.language = payload.language || document.documentElement.lang || 'zh-CN'
    const config = payload.config ?? {}
    writeForm(config)
    postConfigChange(config)
    setStatus('', '')
  })
})()
