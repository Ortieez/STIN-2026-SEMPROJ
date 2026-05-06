import { useState, useEffect } from 'react'
/* @ts-ignore */
import './App.css'
import Login from './components/Login'
import Dashboard from './components/Dashboard'

function App() {
  const [token, setToken] = useState<string | null>(localStorage.getItem('auth_token'))

  useEffect(() => {
    if (token) {
      localStorage.setItem('auth_token', token)
    } else {
      localStorage.removeItem('auth_token')
    }
  }, [token])

  const handleLogin = (newToken: string) => {
    setToken(newToken)
  }

  const handleLogout = () => {
    setToken(null)
  }

  return (
    <div className="app-container">
      {!token ? (
        <Login onLogin={handleLogin} />
      ) : (
        <Dashboard token={token} onLogout={handleLogout} />
      )}
    </div>
  )
}

export default App
