import { useEffect, useState } from 'react'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'
import type { FormEvent } from 'react'
import DashboardSidebar from '../DashboardSidebar'
import '../DashboardPage.css'
import './TopUpPage.css'

type User = {
  id: number
  name: string
  email: string
}

type Wallet = {
  balance: number
  currency: string
}

type TopUpHistoryItem = {
  id: number
  user_id: number
  amount: number
  currency: string
  status: string
  payment_type: string
  redirect_url?: string
  created_at: string
}

type TopUpHistoryResponse = {
  data: TopUpHistoryItem[]
  total: number
}

type TopUpResponse = {
  redirect_url: string
}

type TopUpLocationState = {
  user?: User
}

type PaymentMethod = 'qris' | 'virtual_account' | 'credit_card'

const amountOptions = [25000, 50000, 100000, 200000, 500000, 1000000]
const paymentMethods: { id: PaymentMethod; title: string; description: string; icon: 'qris' | 'bank' | 'card' }[] = [
  {
    id: 'qris',
    title: 'QRIS',
    description: 'GoPay, OVO, DANA, ShopeePay, dan lainnya',
    icon: 'qris',
  },
  {
    id: 'virtual_account',
    title: 'Virtual Account',
    description: 'BCA, BNI, BRI, Mandiri, dan bank lainnya',
    icon: 'bank',
  },
  {
    id: 'credit_card',
    title: 'Kartu Debit / Kredit',
    description: 'Bayar menggunakan kartu debit atau kredit',
    icon: 'card',
  },
]

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080').replace(/\/$/, '')

function formatRupiah(amount: number, currency = 'IDR') {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency,
    maximumFractionDigits: 0,
  }).format(amount)
}

function formatTransactionDate(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Tanggal tidak tersedia'

  const formattedDate = new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(date)
  const formattedTime = new Intl.DateTimeFormat('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
    hourCycle: 'h23',
  }).format(date)

  return `${formattedDate}, ${formattedTime}`
}

function getStatusDisplay(status: string) {
  switch (status.toLowerCase()) {
    case 'paid':
    case 'settlement':
    case 'capture':
      return { label: 'Berhasil', className: 'success' }
    case 'pending':
      return { label: 'Diproses', className: 'pending' }
    case 'expired':
    case 'expire':
    case 'cancelled':
    case 'cancel':
    case 'failed':
    case 'failure':
      return { label: 'Gagal', className: 'failed' }
    default:
      return { label: status || 'Tidak diketahui', className: 'pending' }
  }
}

function HistoryPaymentIcon({ paymentType }: { paymentType: string }) {
  if (paymentType.toLowerCase().includes('bank') || paymentType.toLowerCase().includes('va')) {
    return (
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path d="m3 9 9-6 9 6H3Zm2 2v7m5-7v7m4-7v7m5-7v7M3 21h18M2 18h20" />
      </svg>
    )
  }
  if (paymentType.toLowerCase().includes('card')) {
    return (
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <rect x="2.5" y="5" width="19" height="14" rx="2" />
        <path d="M3 9h18M6 15h4" />
      </svg>
    )
  }
  return (
    <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path d="M3 3h7v7H3zM14 3h7v7h-7zM3 14h7v7H3z" />
      <path d="M14 14h3v3h-3zm5 0h2m-7 5h3m2-2v4m-5-1h2" />
    </svg>
  )
}

function PaymentMethodIcon({ icon }: { icon: 'qris' | 'bank' | 'card' }) {
  if (icon === 'qris') {
    return (
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path d="M3 3h7v7H3zM14 3h7v7h-7zM3 14h7v7H3z" />
        <path d="M14 14h3v3h-3zm5 0h2m-7 5h3m2-2v4m-5-1h2" />
      </svg>
    )
  }
  if (icon === 'bank') {
    return (
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path d="m3 9 9-6 9 6H3Zm2 2v7m5-7v7m4-7v7m5-7v7M3 21h18M2 18h20" />
      </svg>
    )
  }
  return (
    <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <rect x="2.5" y="5" width="19" height="14" rx="2" />
      <path d="M3 9h18M6 15h4" />
    </svg>
  )
}

function TopUpPage() {
  const location = useLocation()
  const navigate = useNavigate()
  const user = (location.state as TopUpLocationState | null)?.user
  const userId = user?.id
  const [wallet, setWallet] = useState<Wallet | null>(null)
  const [topUps, setTopUps] = useState<TopUpHistoryItem[]>([])
  const [historyLoading, setHistoryLoading] = useState(true)
  const [historyError, setHistoryError] = useState('')
  const [historyActionError, setHistoryActionError] = useState('')
  const [amount, setAmount] = useState(25000)
  const [customAmount, setCustomAmount] = useState('')
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('qris')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (userId === undefined) return
    const controller = new AbortController()

    async function loadWallet() {
      try {
        const response = await fetch(`${apiBaseUrl}/wallets/user/${userId}`, { signal: controller.signal })
        const data: Wallet | { error?: string } = await response.json()
        if (!response.ok) {
          throw new Error('error' in data && data.error ? data.error : 'Saldo belum dapat dimuat.')
        }
        setWallet(data as Wallet)
      } catch (requestError) {
        if (controller.signal.aborted) return
        setError(requestError instanceof Error ? requestError.message : 'Saldo belum dapat dimuat.')
      }
    }

    async function loadTopUpHistory() {
      setHistoryLoading(true)
      setHistoryError('')
      try {
        const response = await fetch(`${apiBaseUrl}/wallets/user/${userId}/topups`, {
          signal: controller.signal,
        })
        const data: TopUpHistoryResponse | { error?: string } = await response.json()
        if (!response.ok) {
          throw new Error('error' in data && data.error ? data.error : 'Riwayat top up belum dapat dimuat.')
        }
        if (!('data' in data) || !Array.isArray(data.data)) {
          throw new Error('Format riwayat top up dari server tidak valid.')
        }
        setTopUps(data.data.slice(0, 4))
      } catch (requestError) {
        if (controller.signal.aborted) return
        setHistoryError(requestError instanceof Error ? requestError.message : 'Riwayat top up belum dapat dimuat.')
      } finally {
        if (!controller.signal.aborted) setHistoryLoading(false)
      }
    }

    void loadWallet()
    void loadTopUpHistory()
    return () => controller.abort()
  }, [userId])

  if (!user) {
    return <Navigate to="/login" replace />
  }

  const currentUser = user
  const selectedAmount = customAmount === '' ? amount : Number(customAmount)

  function chooseAmount(value: number) {
    setAmount(value)
    setCustomAmount('')
    setError('')
  }

  function openPendingPayment(redirectUrl: string) {
    setHistoryActionError('')
    if (!redirectUrl) {
      setHistoryActionError('Link pembayaran tidak tersedia. Muat ulang riwayat atau buat transaksi top up baru.')
      return
    }

    let paymentUrl: URL
    try {
      paymentUrl = new URL(redirectUrl)
    } catch {
      setHistoryActionError('URL pembayaran tidak valid.')
      return
    }

    const midtransHosts = ['app.midtrans.com', 'app.sandbox.midtrans.com']
    if (paymentUrl.protocol !== 'https:' || !midtransHosts.includes(paymentUrl.hostname)) {
      setHistoryActionError('URL pembayaran Midtrans tidak valid.')
      return
    }

    const paymentWindow = window.open('about:blank', '_blank')
    if (!paymentWindow) {
      setHistoryActionError('Tab pembayaran diblokir browser. Izinkan pop-up untuk situs ini lalu coba lagi.')
      return
    }
    paymentWindow.opener = null
    paymentWindow.location.replace(paymentUrl.toString())
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')

    if (!Number.isSafeInteger(selectedAmount) || selectedAmount < 10000) {
      setError('Nominal top up minimal Rp10.000.')
      return
    }

    const paymentWindow = window.open('about:blank', '_blank')
    if (!paymentWindow) {
      setError('Tab pembayaran diblokir browser. Izinkan pop-up untuk situs ini lalu coba lagi.')
      return
    }
    paymentWindow.opener = null
    paymentWindow.document.title = 'Menyiapkan pembayaran'
    paymentWindow.document.body.textContent = 'Menyiapkan halaman pembayaran Midtrans...'

    setLoading(true)
    try {
      const response = await fetch(`${apiBaseUrl}/wallets/user/${currentUser.id}/topups`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ amount: selectedAmount, payment_method: paymentMethod }),
      })
      const data: TopUpResponse | { error?: string } = await response.json()

      if (!response.ok) {
        throw new Error('error' in data && data.error ? data.error : 'Transaksi top up gagal dibuat.')
      }
      if (!('redirect_url' in data) || !data.redirect_url) {
        throw new Error('URL pembayaran tidak diterima dari server.')
      }

      const paymentUrl = new URL(data.redirect_url)
      if (paymentUrl.protocol !== 'https:' && paymentUrl.protocol !== 'http:') {
        throw new Error('URL pembayaran dari server tidak valid.')
      }
      paymentWindow.location.replace(paymentUrl.toString())
    } catch (requestError) {
      paymentWindow.close()
      setError(requestError instanceof Error ? requestError.message : 'Tidak dapat terhubung ke server.')
      setLoading(false)
    }
  }

  return (
    <main className="dashboard-shell">
      <DashboardSidebar user={user} activePage="topup" />
      <section className="dashboard-content topup-content">
        <header className="dashboard-topbar topup-topbar">
          <span className="dashboard-page-label">Top Up</span>
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

        <div className="topup-layout">
          <section className="topup-main-column">
            <button className="topup-back" type="button" onClick={() => navigate('/dashboard', { state: { user } })}>
              <span aria-hidden="true">←</span> Dashboard
            </button>
            <div className="topup-heading">
              <h1>Top Up Saldo</h1>
              <p>Tambahkan saldo ke GoWallet kamu dengan mudah dan aman melalui berbagai metode pembayaran yang tersedia.</p>
            </div>

            <form className="topup-form-card" onSubmit={handleSubmit}>
              <section className="topup-step" aria-labelledby="amount-heading">
                <div className="step-heading">
                  <span className="step-number">1</span>
                  <div>
                    <h2 id="amount-heading">Pilih Nominal</h2>
                    <p>Pilih nominal top up yang kamu inginkan.</p>
                  </div>
                </div>
                <div className="amount-grid">
                  {amountOptions.map((option) => (
                    <button
                      className={`amount-option${customAmount === '' && amount === option ? ' selected' : ''}`}
                      type="button"
                      key={option}
                      aria-pressed={customAmount === '' && amount === option}
                      onClick={() => chooseAmount(option)}
                    >
                      {formatRupiah(option)}
                    </button>
                  ))}
                </div>
                <label className="custom-amount-label" htmlFor="custom-amount">Nominal lain</label>
                <div className="custom-amount-field">
                  <input
                    id="custom-amount"
                    type="number"
                    min="10000"
                    step="1000"
                    inputMode="numeric"
                    placeholder="Masukkan nominal (min. Rp10.000)"
                    value={customAmount}
                    onChange={(event) => {
                      setCustomAmount(event.target.value)
                      setError('')
                    }}
                  />
                  <span>Rp</span>
                </div>
              </section>

              <section className="topup-step payment-step" aria-labelledby="payment-heading">
                <div className="step-heading">
                  <span className="step-number">2</span>
                  <div>
                    <h2 id="payment-heading">Metode Pembayaran</h2>
                    <p>Pilih metode pembayaran yang kamu gunakan.</p>
                  </div>
                </div>
                <div className="payment-method-list">
                  {paymentMethods.map((method) => (
                    <button
                      className={`payment-method-option${paymentMethod === method.id ? ' selected' : ''}`}
                      type="button"
                      key={method.id}
                      aria-pressed={paymentMethod === method.id}
                      onClick={() => setPaymentMethod(method.id)}
                    >
                      <span className="payment-method-icon"><PaymentMethodIcon icon={method.icon} /></span>
                      <span className="payment-method-copy">
                        <strong>{method.title}</strong>
                        <small>{method.description}</small>
                      </span>
                      <span className="payment-method-check" aria-hidden="true">
                        {paymentMethod === method.id ? '✓' : ''}
                      </span>
                    </button>
                  ))}
                </div>
              </section>

              {error && <p className="topup-error" role="alert">{error}</p>}
              <button className="topup-submit" type="submit" disabled={loading}>
                {loading ? 'Menyiapkan pembayaran...' : 'Lanjutkan ke Pembayaran'}
                {!loading && <span aria-hidden="true">→</span>}
              </button>
              <p className="topup-secure-note">Pembayaran diproses dengan aman melalui Midtrans.</p>
            </form>
          </section>

          <aside className="topup-info-column">
            <section className="topup-balance-card">
              <span className="topup-balance-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none">
                  <path d="M4 7h15a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h13" />
                  <path d="M3 9h18m-6 5h6" />
                  <circle cx="17" cy="14" r=".8" fill="currentColor" stroke="none" />
                </svg>
              </span>
              <div>
                <span>Saldo GoWallet</span>
                <strong>{wallet ? formatRupiah(wallet.balance, wallet.currency) : 'Memuat saldo...'}</strong>
              </div>
            </section>

            <section className="topup-info-card">
              <h2><span aria-hidden="true">ⓘ</span> Informasi Top Up</h2>
              <div className="topup-info-item">
                <span className="info-icon" aria-hidden="true">ϟ</span>
                <div><strong>Proses instan</strong><p>Saldo akan masuk setelah pembayaran berhasil diverifikasi.</p></div>
              </div>
              <div className="topup-info-item">
                <span className="info-icon" aria-hidden="true">✓</span>
                <div><strong>Aman &amp; terpercaya</strong><p>Pembayaran diproses melalui Midtrans.</p></div>
              </div>
              <div className="topup-info-item">
                <span className="info-icon" aria-hidden="true">▣</span>
                <div><strong>Pilihan metode</strong><p>QRIS, virtual account, dan kartu debit/kredit.</p></div>
              </div>
            </section>

            <section className="topup-history-card">
              <div className="topup-history-heading">
                <h2>Riwayat Top Up Terakhir</h2>
              </div>
              {historyLoading ? (
                <div className="topup-history-empty">Memuat riwayat transaksi...</div>
              ) : historyError ? (
                <div className="topup-history-empty topup-history-error" role="alert">{historyError}</div>
              ) : topUps.length === 0 ? (
                <div className="topup-history-empty">Riwayat transaksi belum tersedia.</div>
              ) : (
                <ul className="topup-history-list">
                  {topUps.map((topUp) => {
                    const status = getStatusDisplay(topUp.status)
                    const isPending = topUp.status.toLowerCase() === 'pending'
                    return (
                      <li className="topup-history-item" key={topUp.id}>
                        {isPending ? (
                          <button
                            className="topup-history-open"
                            type="button"
                            aria-label={`Lanjutkan pembayaran top up ${formatRupiah(topUp.amount, topUp.currency || 'IDR')}`}
                            onClick={() => openPendingPayment(topUp.redirect_url ?? '')}
                          >
                            <span className="topup-history-icon"><HistoryPaymentIcon paymentType={topUp.payment_type} /></span>
                            <span className="topup-history-details">
                              <strong>{formatRupiah(topUp.amount, topUp.currency || 'IDR')}</strong>
                              <small>{formatTransactionDate(topUp.created_at)}</small>
                            </span>
                            <span className={`topup-history-status ${status.className}`}>{status.label}</span>
                            <span className="topup-history-arrow" aria-hidden="true">›</span>
                          </button>
                        ) : (
                          <>
                            <span className="topup-history-icon"><HistoryPaymentIcon paymentType={topUp.payment_type} /></span>
                            <span className="topup-history-details">
                              <strong>{formatRupiah(topUp.amount, topUp.currency || 'IDR')}</strong>
                              <small>{formatTransactionDate(topUp.created_at)}</small>
                            </span>
                            <span className={`topup-history-status ${status.className}`}>{status.label}</span>
                            <span className="topup-history-arrow" aria-hidden="true">›</span>
                          </>
                        )}
                      </li>
                    )
                  })}
                </ul>
              )}
              {historyActionError && <p className="topup-history-error" role="alert">{historyActionError}</p>}
              <p className="topup-history-note">Status pembayaran akan diperbarui setelah Midtrans mengirim notifikasi.</p>
            </section>
          </aside>
        </div>
      </section>
    </main>
  )
}

export default TopUpPage
