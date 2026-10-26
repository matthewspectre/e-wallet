import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import './LoginPage.css'

type User = {
  id: number
  name: string
  email: string
}

type LoginResponse = { user: User }

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080').replace(/\/$/, '')

function LoginPage() {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await fetch(`${apiBaseUrl}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email.trim(), password }),
      })
      const data: LoginResponse | { error?: string } = await response.json()

      if (!response.ok) {
        throw new Error('error' in data && data.error ? data.error : 'Login gagal. Silakan coba lagi.')
      }

      const login = data as LoginResponse
      navigate('/dashboard', { state: { user: login.user }, replace: true })
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Tidak dapat terhubung ke server.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="page-shell">
      <section className="login-layout" aria-label="Login Go-Wallet">
        <div className="intro-panel">
          <div className="brand">
            <div className="brand-mark" aria-hidden="true">W</div>
            <span>Go-Wallet</span>
          </div>
          <div className="intro-copy">
            <span className="eyebrow">DOMPET DIGITALMU</span>
            <h1>Atur uangmu, dengan lebih mudah.</h1>
            <p>Satu tempat untuk melihat saldo dan mengelola aktivitas wallet kamu.</p>
          </div>
          <div className="decorative-wallet" aria-hidden="true">
            <div className="decorative-chip" />
            <div className="decorative-line" />
            <div className="decorative-circle" />
          </div>
          <p className="intro-caption">Sederhana, aman, dan selalu dalam kendalimu.</p>
        </div>

        <div className="form-panel">
          <div className="form-heading">
            <span className="eyebrow">SELAMAT DATANG KEMBALI</span>
            <h2>Masuk ke akunmu</h2>
            <p>Masukkan email dan password yang kamu gunakan saat mendaftar.</p>
          </div>

          <form className="login-form" onSubmit={handleSubmit}>
            <label htmlFor="email">Email</label>
            <input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="nama@email.com"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              required
            />

            <div className="password-label">
              <label htmlFor="password">Password</label>
            </div>
            <input
              id="password"
              type="password"
              autoComplete="current-password"
              placeholder="Masukkan password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
            />

            {error && <p className="form-error" role="alert">{error}</p>}

            <button className="primary-button" type="submit" disabled={loading}>
              {loading ? 'Memproses...' : 'Masuk'}
              {!loading && <span aria-hidden="true">→</span>}
            </button>
          </form>

          <p className="form-footnote">Belum punya akun? Buat akun melalui layanan registrasi Go-Wallet.</p>
        </div>
      </section>
      <footer className="page-footer">© 2026 Go-Wallet <span>•</span> Portfolio project</footer>
    </main>
  )
}

export default LoginPage
