package wallet

import "time"

type WalletResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TopUpCreateRequest struct {
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	PaymentMethod string `json:"payment_method"`
}

type TopUpResponse struct {
	ID          int64     `json:"id"`
	OrderID     string    `json:"order_id"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	SnapToken   string    `json:"snap_token"`
	RedirectURL string    `json:"redirect_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type MidtransNotificationRequest struct {
	OrderID      string `json:"order_id" binding:"required"`
	StatusCode   string `json:"status_code" binding:"required"`
	GrossAmount  string `json:"gross_amount" binding:"required"`
	SignatureKey string `json:"signature_key" binding:"required"`
}
