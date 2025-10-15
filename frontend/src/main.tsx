import React from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from './components/app'

const rootEl = document.getElementById('root')

if (!rootEl) {
    throw new Error('Root element not found')
}

const root = createRoot(rootEl)

const qc = new QueryClient()

root.render(
  <React.StrictMode>
    <QueryClientProvider client={qc}>
      <App />
    </QueryClientProvider>
  </React.StrictMode>
)
