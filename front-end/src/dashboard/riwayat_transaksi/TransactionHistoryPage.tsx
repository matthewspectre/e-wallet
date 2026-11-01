import { useEffect, useState } from 'react'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'
import DashboardSidebar from '../DashboardSidebar'
import '../DashboardPage.css'
import './TransactionHistoryPage.css'

type User = {
  id: number
  name: string
  email: string
}

type HistoryItem = {
  id: number
  amount: number
  currency: string
  status: string
  payment_type: string
  created_at: string
}

type HistoryResponse = {
  data: HistoryItem[]
  total: number
}

type LocationState = {
  user?: User
}

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080').replace(/\/$/, '')

function formatRupiah(amount: number, currency = 'IDR') {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency,
    maximumFractionDigits: 0,
  }).format(amount)
}

function formatDate(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Tanggal tidak tersedia'

  return new Intl.DateTimeFormat('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}

function statusDisplay(status: string) {
  switch (status.toLowerCase()) {
    case 'paid':
    case 'settlement':
    case 'capture':
      return { label: 'Berhasil', className: 'success' }
    case 'pending':
      return { label: 'Diproses', className: 'pending' }
    case 'expired':
    case 'expire':
      return { label: 'Kedaluwarsa', className: 'failed' }
    case 'cancelled':
    case 'cancel':
      return { label: 'Dibatalkan', className: 'failed' }
    case 'failed':
    case 'failure':
      return { label: 'Gagal', className: 'failed' }
    default:
      return { label: status || 'Tidak diketahui', className: 'pending' }
  }
}

function TransactionHistoryPage() {
  const location = useLocation()
  const navigate = useNavigate()
  const user = (location.state as LocationState | null)?.user
  const userId = user?.id
  const [transactions, setTransactions] = useState<HistoryItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (userId === undefined) return
    const controller = new AbortController()

    async function loadTransactions() {
      try {
        const response = await fetch(`${apiBaseUrl}/wallets/user/${userId}/topups`, {
          signal: controller.signal,
        })
        const data: HistoryResponse | { error?: string } = await response.json()
        if (!response.ok) {
          throw new Error('error' in data && data.error ? data.error : 'Riwayat transaksi belum dapat dimuat.')
        }
        if (!('data' in data) || !Array.isArray(data.data)) {
          throw new Error('Format riwayat transaksi dari server tidak valid.')
        }

        const sortedTransactions = [...data.data].sort((a, b) => {
          const dateDifference = new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
          return dateDifference || b.id - a.id
        })
        setTransactions(sortedTransactions)
      } catch (requestError) {
        if (controller.signal.aborted) return
        setError(requestError instanceof Error ? requestError.message : 'Riwayat transaksi belum dapat dimuat.')
      } finally {
        if (!controller.signal.aborted) setLoading(false)
      }
    }

    void loadTransactions()
    return () => controller.abort()
  }, [userId])

  if (!user) {
    return <Navigate to="/login" replace />
  }

  return (
    <main className="dashboard-shell">
      <DashboardSidebar user={user} activePage="history" />
      <section className="dashboard-content history-content">
        <header className="dashboard-topbar">
          <span className="dashboard-page-label">Riwayat Transaksi</span>
          <div className="dashboard-account">
            <div className="dashboard-avatar" aria-hidden="true">{user.name.charAt(0).toUpperCase()}</div>
            <div className="dashboard-account-copy">
              <span>{user.name}</span>
              <small>{user.email}</small>
            </div>
            <button className="dashboard-logout" type="button" onClick={() => navigate('/login', { replace: true })}>
              Keluar
            </button>
          </div>
        </header>

        <div className="history-page-layout">
          <button className="history-back" type="button" onClick={() => navigate('/dashboard', { state: { user } })}>
            <span aria-hidden="true">←</span> Dashboard
          </button>
          <div className="history-page-heading">
            <h1>Riwayat Transaksi</h1>
            <p>Lihat status dan detail transaksi top up akun GoWallet kamu.</p>
          </div>

          <section className="history-card" aria-label="Daftar riwayat transaksi">
            {loading ? (
              <p className="history-message">Memuat riwayat transaksi...</p>
            ) : error ? (
              <p className="history-message error" role="alert">{error}</p>
            ) : transactions.length === 0 ? (
              <p className="history-message">Belum ada transaksi.</p>
            ) : (
              <div className="history-table-wrap">
                <table className="history-table">
                  <thead>
                    <tr>
                      <th scope="col">Transaksi</th>
                      <th scope="col">Metode Pembayaran</th>
                      <th scope="col">Tanggal</th>
                      <th scope="col">Status</th>
                      <th scope="col">Nominal</th>
                    </tr>
                  </thead>
                  <tbody>
                    {transactions.map((transaction) => {
                      const status = statusDisplay(transaction.status)
                      return (
                        <tr key={transaction.id}>
                          <td data-label="Transaksi">Top Up</td>
                          <td data-label="Metode Pembayaran">{transaction.payment_type || '—'}</td>
                          <td data-label="Tanggal">{formatDate(transaction.created_at)}</td>
                          <td data-label="Status">
                            <span className={`history-status ${status.className}`}>{status.label}</span>
                          </td>
                          <td className="history-amount" data-label="Nominal">
                            {formatRupiah(transaction.amount, transaction.currency || 'IDR')}
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </section>
        </div>
      </section>
    </main>
  )
}

export default TransactionHistoryPage
