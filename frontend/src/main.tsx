import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import './leaflet-fix.ts' // Import marker fix
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
