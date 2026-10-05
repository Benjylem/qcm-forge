import { useEffect, useState } from 'react'

type HealthStatus = 'loading' | 'ok' | 'down'

function App() {
  const [status, setStatus] = useState<HealthStatus>('loading')

  useEffect(() => {
    fetch('/health')
      .then((res) => setStatus(res.ok ? 'ok' : 'down'))
      .catch(() => setStatus('down'))
  }, [])

  return (
    <main>
      <h1>qcm-forge</h1>
      <p>API : {status}</p>
    </main>
  )
}

export default App
