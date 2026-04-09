import { Routes, Route, Navigate } from 'react-router-dom'
import { useState, useEffect } from 'react'
import SignIn from './pages/SignIn'
import SignUp from './pages/SignUp'
import Dashboard from './pages/Dashboard'

export default function App() {
  const [authed, setAuthed] = useState(() => !!localStorage.getItem('authed'))

  const login = () => {
    localStorage.setItem('authed', '1')
    setAuthed(true)
  }

  const logout = () => {
    localStorage.removeItem('authed')
    setAuthed(false)
  }

  return (
    <Routes>
      <Route path="/sign-in" element={authed ? <Navigate to="/" /> : <SignIn onLogin={login} />} />
      <Route path="/sign-up" element={authed ? <Navigate to="/" /> : <SignUp onLogin={login} />} />
      <Route path="/*"       element={authed ? <Dashboard onLogout={logout} /> : <Navigate to="/sign-in" />} />
    </Routes>
  )
}
