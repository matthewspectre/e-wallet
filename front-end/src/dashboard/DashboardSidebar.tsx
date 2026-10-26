import { useNavigate } from 'react-router-dom'

type User = {
  id: number
  name: string
  email: string
}

type MenuIconName = 'dashboard' | 'topup' | 'transfer' | 'history' | 'settings'

type DashboardSidebarProps = {
  user: User
  activePage: 'dashboard' | 'topup'
}

const menuItems: { label: string; icon: MenuIconName; page?: 'dashboard' | 'topup' }[] = [
  { label: 'Dashboard', icon: 'dashboard', page: 'dashboard' },
  { label: 'Top Up', icon: 'topup', page: 'topup' },
  { label: 'Transfer', icon: 'transfer' },
  { label: 'Riwayat Transaksi', icon: 'history' },
  { label: 'Pengaturan', icon: 'settings' },
]

function MenuIcon({ name }: { name: MenuIconName }) {
  const commonProps = {
    width: 18,
    height: 18,
    viewBox: '0 0 24 24',
    fill: 'none',
    stroke: 'currentColor',
    strokeWidth: 1.8,
    strokeLinecap: 'round' as const,
    strokeLinejoin: 'round' as const,
    'aria-hidden': true as const,
  }

  switch (name) {
    case 'dashboard':
      return <svg {...commonProps}><path d="m3 10 9-7 9 7v10a1 1 0 0 1-1 1h-6v-7h-4v7H4a1 1 0 0 1-1-1z" /></svg>
    case 'topup':
      return <svg {...commonProps}><circle cx="12" cy="12" r="9" fill="currentColor" stroke="none" /><path d="M12 7v10M7 12h10" stroke="#064d3d" strokeWidth="2" /></svg>
    case 'transfer':
      return <svg {...commonProps}><path d="m21 3-7.2 18-3.8-7-7-3.8z" /><path d="M21 3 10 14" /></svg>
    case 'history':
      return <svg {...commonProps}><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2M3 12H1m2-5H1m2 10H1" /></svg>
    case 'settings':
      return <svg {...commonProps}><circle cx="12" cy="12" r="3" /><path d="m19.4 15 .1.1 1.4 1.1-1.4 2.4-1.7-.6a8 8 0 0 1-1.6.9l-.3 1.8h-2.8l-.3-1.8a8 8 0 0 1-1.6-.9l-1.7.6-1.4-2.4 1.4-1.1A7 7 0 0 1 7 13.2l-1.8-.3v-2.8l1.8-.3a7 7 0 0 1 .5-1.7L6.1 7l1.4-2.4 1.7.6a8 8 0 0 1 1.6-.9l.3-1.8h2.8l.3 1.8a8 8 0 0 1 1.6.9l1.7-.6L21 7l-1.4 1.1a7 7 0 0 1 .5 1.7l1.8.3v2.8l-1.8.3a7 7 0 0 1-.7 1.8Z" /></svg>
  }
}

function DashboardSidebar({ user, activePage }: DashboardSidebarProps) {
  const navigate = useNavigate()

  return (
    <aside className="dashboard-sidebar" aria-label="Navigasi utama">
      <div className="sidebar-brand">
        <span className="sidebar-brand-mark" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none">
            <path d="M4 7.5h15a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7.5a2 2 0 0 1 1-1.7l12-3.3" />
            <path d="M3 9h18" />
            <circle cx="17" cy="13.5" r="1" fill="currentColor" stroke="none" />
          </svg>
        </span>
        <span>Go<span className="sidebar-brand-accent">Wallet</span></span>
      </div>

      <nav className="sidebar-nav">
        {menuItems.map((item) => {
          const active = item.page === activePage
          const target = item.page === 'topup' ? '/dashboard/topup' : '/dashboard'
          return (
            <button
              className={`sidebar-nav-item${active ? ' active' : ''}`}
              type="button"
              key={item.label}
              disabled={!item.page || active}
              aria-current={active ? 'page' : undefined}
              title={item.page ? undefined : 'Segera hadir'}
              onClick={() => navigate(target, { state: { user } })}
            >
              <MenuIcon name={item.icon} />
              <span>{item.label}</span>
            </button>
          )
        })}
      </nav>

      <section className="sidebar-promo" aria-label="Informasi Go-Wallet">
        <div className="promo-wallet-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none">
            <path d="M4 7h15a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h13" />
            <path d="M3 9h18m-6 5h6" />
            <circle cx="17" cy="14" r=".8" fill="currentColor" stroke="none" />
          </svg>
        </div>
        <p>Kelola keuangan lebih mudah dengan GoWallet</p>
        <button className="promo-link" type="button" disabled>
          Pelajari lebih lanjut <span aria-hidden="true">→</span>
        </button>
      </section>
    </aside>
  )
}

export default DashboardSidebar
