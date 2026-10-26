import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import DashboardPage from './dashboard/DashboardPage'
import TopUpPage from './dashboard/topup/TopUpPage'
import LoginPage from './login/LoginPage'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/dashboard/topup" element={<TopUpPage />} />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
