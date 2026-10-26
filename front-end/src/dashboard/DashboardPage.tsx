import { useEffect, useState } from 'react'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'
import DashboardSidebar from './DashboardSidebar'
import './DashboardPage.css'

type User = {
  id: number
  name: string
  email: string
}

type DashboardLocationState = {
  user?: User
}

type Wallet = {
  balance: number
  currency: string
}

type QuickActionName = 'Top Up' | 'Transfer' | 'Riwayat' | 'Lainnya'
type QuickActionIconName = 'topup' | 'transfer' | 'history' | 'more'

const quickActions: { label: QuickActionName; icon: QuickActionIconName }[] = [
  { label: 'Top Up', icon: 'topup' },
  { label: 'Transfer', icon: 'transfer' },
  { label: 'Riwayat', icon: 'history' },
  { label: 'Lainnya', icon: 'more' },
]

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080').replace(/\/$/, '')

function QuickActionIcon({ name }: { name: QuickActionIconName }) {
  if (name === 'more') {
    return (
      <svg className="quick-action-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <rect x="4" y="4" width="6" height="6" rx="1.5" />
        <rect x="14" y="4" width="6" height="6" rx="1.5" />
        <rect x="4" y="14" width="6" height="6" rx="1.5" />
        <rect x="14" y="14" width="6" height="6" rx="1.5" />
      </svg>
    )
  }
  if (name === 'topup') {
    return <svg className="quick-action-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="12" r="9" fill="currentColor" stroke="none" /><path d="M12 7v10M7 12h10" stroke="#064d3d" strokeWidth="2" /></svg>
  }
  if (name === 'transfer') {
    return <svg className="quick-action-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="m21 3-7.2 18-3.8-7-7-3.8z" /><path d="M21 3 10 14" /></svg>
  }
  return <svg className="quick-action-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></svg>
}

function DashboardPage() {
  const location = useLocation()
  const navigate = useNavigate()
  const user = (location.state as DashboardLocationState | null)?.user
  const userId = user?.id
  const [wallet, setWallet] = useState<Wallet | null>(null)
  const [walletError, setWalletError] = useState('')
  const [isBalanceVisible, setIsBalanceVisible] = useState(true)
  const [actionMessage, setActionMessage] = useState('')

  useEffect(() => {
    if (userId === undefined) return
    const controller = new AbortController()

    async function loadWallet() {
      try {
        const response = await fetch(`${apiBaseUrl}/wallets/user/${userId}`, {
          signal: controller.signal,
        })
        const data: Wallet | { error?: string } = await response.json()
        if (!response.ok) {
          throw new Error('error' in data && data.error ? data.error : 'Saldo belum dapat dimuat.')
        }
        setWallet(data as Wallet)
      } catch (error) {
        if (controller.signal.aborted) return
        setWalletError(error instanceof Error ? error.message : 'Saldo belum dapat dimuat.')
      }
    }

    void loadWallet()
    return () => controller.abort()
  }, [userId])

  if (!user) {
    return <Navigate to="/login" replace />
  }

  const formattedBalance = wallet
    ? new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: wallet.currency || 'IDR',
        maximumFractionDigits: 0,
      }).format(wallet.balance)
    : null

  return (
    <main className="dashboard-shell">
      <DashboardSidebar user={user} activePage="dashboard" />

      <section className="dashboard-content">
        <header className="dashboard-topbar">
          <span className="dashboard-page-label">Dashboard</span>
          <div className="dashboard-account">
            <div className="dashboard-avatar" aria-hidden="true">{user.name.charAt(0).toUpperCase()}</div>
            <div className="dashboard-account-copy">
              <span>{user.name}</span>
              <small>{user.email}</small>
            </div>
            <button
              className="dashboard-logout"
              type="button"
              onClick={() => navigate('/login', { replace: true })}
            >
              Keluar
            </button>
          </div>
        </header>
        <section className="dashboard-welcome" aria-labelledby="dashboard-title">
          <p className="welcome-kicker">Selamat datang kembali,</p>
          <h1 id="dashboard-title">{user.name} <span aria-label="melambaikan tangan">👋</span></h1>
          <p className="welcome-subtitle">Kelola saldo dan transaksi kamu dengan mudah dan aman.</p>

          <section className="wallet-summary" aria-label="Ringkasan wallet">
            <div className="total-balance">
              <div className="balance-label-row">
                <svg className="summary-wallet-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path d="M4 7h15a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h13" />
                  <path d="M3 9h18m-6 5h6" />
                  <circle cx="17" cy="14" r=".8" fill="currentColor" stroke="none" />
                </svg>
                <span>Total Saldo</span>
                <button
                  className="balance-visibility"
                  type="button"
                  onClick={() => setIsBalanceVisible((visible) => !visible)}
                  aria-label={isBalanceVisible ? 'Sembunyikan saldo' : 'Tampilkan saldo'}
                  title={isBalanceVisible ? 'Sembunyikan saldo' : 'Tampilkan saldo'}
                >
                  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
                    {isBalanceVisible
                      ? <><path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6S2 12 2 12Z" /><circle cx="12" cy="12" r="2.5" /></>
                      : <><path d="m3 3 18 18M10.6 6.2A10 10 0 0 1 12 6c6.5 0 10 6 10 6a16 16 0 0 1-3.1 3.6M6.2 6.2C3.5 7.8 2 12 2 12s3.5 6 10 6a10 10 0 0 0 3.1-.5" /></>
                    }
                  </svg>
                </button>
              </div>
              <strong className="balance-amount">
                {isBalanceVisible ? formattedBalance ?? (walletError || 'Memuat saldo...') : '••••••••'}
              </strong>
              <span className="balance-caption">GoWallet Balance</span>
            </div>

            <div className="summary-divider" />
            <div className="transaction-summary income">
              <span className="transaction-symbol" aria-hidden="true">↗</span>
              <div>
                <span className="transaction-label">Pemasukan</span>
                <strong>Belum tersedia</strong>
              </div>
            </div>
            <div className="transaction-summary expense">
              <span className="transaction-symbol" aria-hidden="true">↓</span>
              <div>
                <span className="transaction-label">Pengeluaran</span>
                <strong>Belum tersedia</strong>
              </div>
            </div>
          </section>

          <nav className="quick-actions" aria-label="Aksi wallet">
            {quickActions.map((action) => (
              <button
                className="quick-action"
                type="button"
                key={action.label}
                onClick={() => {
                  if (action.label === 'Top Up') {
                    navigate('/dashboard/topup', { state: { user } })
                    return
                  }
                  setActionMessage(`${action.label} akan segera hadir.`)
                }}
              >
                <span className={`quick-action-icon-wrap${action.label === 'Top Up' ? ' primary' : ''}`}>
                  <QuickActionIcon name={action.icon} />
                </span>
                <span>{action.label}</span>
              </button>
            ))}
          </nav>
          {actionMessage && <p className="action-message" role="status">{actionMessage}</p>}
        </section>
      </section>
    </main>
  )
}

export default DashboardPage
