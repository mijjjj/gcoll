(() => {
  const form = document.querySelector('#modbus-config-form')
  const testButton = document.querySelector('#test-button')
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

  form.addEventListener('input', () => {
    window.parent.postMessage({ type: 'gcoll:config-change', payload: readForm() }, '*')
  })

  form.addEventListener('submit', (event) => {
    event.preventDefault()
    window.parent.postMessage({ type: 'gcoll:config-change', payload: readForm() }, '*')
    window.parent.postMessage({ type: 'gcoll:save-config' }, '*')
  })

  testButton.addEventListener('click', () => {
    window.parent.postMessage({ type: 'gcoll:config-change', payload: readForm() }, '*')
    window.parent.postMessage({ type: 'gcoll:test-connection' }, '*')
  })

  window.addEventListener('message', (event) => {
    if (event.data?.type !== 'gcoll:init') return
    writeForm(event.data.payload?.config ?? {})
  })
})()
