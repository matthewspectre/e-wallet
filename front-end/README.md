# Go-Wallet Frontend

React + TypeScript + Vite frontend for the Go-Wallet API.

## Run locally

From the repository root, run the Go API. The backend listens on port `8080` by default.

Run the frontend in a separate terminal:

```powershell
cd D:\e-wallet\front-end
npm install
npm run dev
```

Open the Vite URL, normally `http://localhost:5173`. The frontend calls `http://localhost:8080` by default. To use another API URL, create `front-end/.env.local`:

```dotenv
VITE_API_BASE_URL=http://localhost:8080
```

The login form uses the account email and password. A successful login opens `/dashboard` and greets the user. Login state is kept in browser navigation state only; refreshing the dashboard returns the user to login. This prototype does not issue a token or create a persistent authenticated session, and wallet routes are not protected.

The dashboard's **Top Up** action opens `/dashboard/topup`. It loads the user's wallet balance, accepts preset or custom amounts (minimum Rp10.000), and creates a Midtrans Snap transaction through the API. The selected QRIS, virtual account, or card method is passed to Snap as the enabled payment method group. Completing a top up requires the backend's Midtrans notification URL to be publicly reachable. Top-up history is not shown yet because the API does not currently expose a history endpoint.
