import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../api'

export default function SignUp({ onLogin }) {
  const [form, setForm]   = useState({ name: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const set = (k) => (e) => setForm(f => ({ ...f, [k]: e.target.value }))

  const submit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await api.signUp(form)
      await api.signIn({ email: form.email, password: form.password })
      onLogin()
      navigate('/')
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-left">
        <div className="auth-logo">MoodLog</div>
        <div className="auth-quote">
          Begin with<br />one honest<br /><span>moment.</span>
        </div>
        <div className="auth-tagline">Track your mood. Discover your patterns.</div>
      </div>

      <div className="auth-right">
        <div className="auth-box">
          <h1>Create account</h1>
          <p>Already have one? <Link to="/sign-in">Sign in</Link></p>

          {error && <div className="error-msg">{error}</div>}

          <form onSubmit={submit}>
            <div className="form-group">
              <label>Name</label>
              <input type="text" value={form.name} onChange={set('name')} placeholder="Your name" required minLength={3} />
            </div>
            <div className="form-group">
              <label>Email</label>
              <input type="email" value={form.email} onChange={set('email')} placeholder="you@example.com" required />
            </div>
            <div className="form-group">
              <label>Password</label>
              <input type="password" value={form.password} onChange={set('password')} placeholder="Min 6 characters" required minLength={6} />
            </div>
            <button className="btn-primary" type="submit" disabled={loading}>
              {loading ? 'Creating...' : 'Create account'}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}
