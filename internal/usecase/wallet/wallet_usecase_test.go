package wallet

import (
	"errors"
	"testing"
	"time"

	entity "e-wallet/internal/entity/wallet"
)

type repositoryStub struct {
	topUp            *entity.TopUp
	topUps           []*entity.TopUp
	wallet           *entity.Wallet
	startDate        *time.Time
	endDate          *time.Time
	appliedStatus    string
	appliedPayment   string
	getTopUpCalled   bool
	applyStatusCalls int
}

func (r *repositoryStub) GetByUserID(int64) (*entity.Wallet, error) {
	return r.wallet, nil
}

func (r *repositoryStub) GetTopUpsByUserID(_ int64, startDate, endDate *time.Time) ([]*entity.TopUp, error) {
	r.startDate = startDate
	r.endDate = endDate
	return r.topUps, nil
}

func (r *repositoryStub) CreateTopUp(userID int64, data *entity.TopUp) error {
	data.ID = 10
	data.WalletID = 20
	r.topUp = data
	return nil
}

func (r *repositoryStub) UpdateTopUpPayment(orderID, token, redirectURL string) error {
	r.topUp.SnapToken = token
	r.topUp.RedirectURL = redirectURL
	return nil
}

func (r *repositoryStub) GetTopUpByOrderID(string) (*entity.TopUp, error) {
	r.getTopUpCalled = true
	return r.topUp, nil
}

func (r *repositoryStub) ApplyTopUpStatus(_, status, paymentType string) error {
	r.applyStatusCalls++
	r.appliedStatus = status
	r.appliedPayment = paymentType
	return nil
}

type paymentGatewayStub struct {
	validSignature bool
	status         *entity.PaymentStatus
}

func (p paymentGatewayStub) CreateSnapTransaction(string, int64, string) (string, string, error) {
	return "snap-token", "https://sandbox.example/checkout", nil
}

func (p paymentGatewayStub) VerifyNotification(string, string, string, string) bool {
	return p.validSignature
}

func (p paymentGatewayStub) GetTransactionStatus(string) (*entity.PaymentStatus, error) {
	return p.status, nil
}

func TestCreateTopUpPersistsPendingOrderAndReturnsSnapDetails(t *testing.T) {
	repository := &repositoryStub{}
	usecase := NewUsecase(repository, paymentGatewayStub{})

	result, err := usecase.CreateTopUp(7, 25000, "qris")
	if err != nil {
		t.Fatalf("CreateTopUp returned an error: %v", err)
	}
	if result.ID != 10 || result.WalletID != 20 {
		t.Errorf("unexpected persisted IDs: topup=%d wallet=%d", result.ID, result.WalletID)
	}
	if result.UserID != 7 {
		t.Errorf("expected user ID 7 from endpoint input, got %d", result.UserID)
	}
	if result.Amount != 25000 || result.Status != "pending" {
		t.Errorf("unexpected top-up amount/status: amount=%d status=%s", result.Amount, result.Status)
	}
	if result.SnapToken != "snap-token" || result.RedirectURL != "https://sandbox.example/checkout" {
		t.Errorf("Snap response details were not returned")
	}
	if result.OrderID == "" {
		t.Error("expected a generated order ID")
	}
}

func TestCreateTopUpRejectsUnsupportedPaymentMethod(t *testing.T) {
	usecase := NewUsecase(&repositoryStub{}, paymentGatewayStub{})

	_, err := usecase.CreateTopUp(7, 25000, "unknown")
	if !errors.Is(err, ErrInvalidPaymentMethod) {
		t.Fatalf("expected ErrInvalidPaymentMethod, got %v", err)
	}
}

func TestGetTopUpsByUserIDReturnsRepositoryResults(t *testing.T) {
	want := []*entity.TopUp{{ID: 1, UserID: 7, OrderID: "order-1"}}
	usecase := NewUsecase(&repositoryStub{topUps: want}, paymentGatewayStub{})

	got, err := usecase.GetTopUpsByUserID(7, nil, nil)
	if err != nil {
		t.Fatalf("GetTopUpsByUserID returned an error: %v", err)
	}
	if len(got) != 1 || got[0].UserID != 7 || got[0].OrderID != "order-1" {
		t.Fatalf("unexpected top-up list: %#v", got)
	}
}

func TestGetTopUpsByUserIDRejectsNonPositiveID(t *testing.T) {
	usecase := NewUsecase(&repositoryStub{}, paymentGatewayStub{})

	if _, err := usecase.GetTopUpsByUserID(0, nil, nil); err == nil {
		t.Fatal("expected non-positive user ID to be rejected")
	}
}

func TestGetTopUpsByUserIDPassesDateFiltersToRepository(t *testing.T) {
	startDate := time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, time.November, 5, 0, 0, 0, 0, time.UTC)
	repository := &repositoryStub{}
	usecase := NewUsecase(repository, paymentGatewayStub{})

	if _, err := usecase.GetTopUpsByUserID(7, &startDate, &endDate); err != nil {
		t.Fatalf("GetTopUpsByUserID returned an error: %v", err)
	}
	if repository.startDate == nil || !repository.startDate.Equal(startDate) {
		t.Errorf("expected start date %v to reach repository, got %v", startDate, repository.startDate)
	}
	if repository.endDate == nil || !repository.endDate.Equal(endDate) {
		t.Errorf("expected end date %v to reach repository, got %v", endDate, repository.endDate)
	}
}

func TestGetTopUpsByUserIDRejectsReversedDateRange(t *testing.T) {
	startDate := time.Date(2026, time.November, 5, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC)
	usecase := NewUsecase(&repositoryStub{}, paymentGatewayStub{})

	if _, err := usecase.GetTopUpsByUserID(7, &startDate, &endDate); err == nil {
		t.Fatal("expected reversed date range to be rejected")
	}
}

func TestGetByUserIDRefreshesPendingPaymentBeforeReturningWallet(t *testing.T) {
	repository := &repositoryStub{
		wallet: &entity.Wallet{UserID: 7, Balance: 25000},
		topUps: []*entity.TopUp{{
			UserID:  7,
			OrderID: "order-1",
			Amount:  25000,
			Status:  "pending",
		}},
	}
	gateway := paymentGatewayStub{
		status: &entity.PaymentStatus{
			OrderID:           "order-1",
			GrossAmount:       "25000.00",
			Currency:          "IDR",
			TransactionStatus: "settlement",
			PaymentType:       "qris",
		},
	}
	usecase := NewUsecase(repository, gateway)

	wallet, err := usecase.GetByUserID(7)
	if err != nil {
		t.Fatalf("GetByUserID returned an error: %v", err)
	}
	if wallet != repository.wallet {
		t.Fatal("expected wallet to be returned")
	}
	if repository.appliedStatus != "paid" || repository.appliedPayment != "qris" {
		t.Fatalf("expected settled payment to be applied, got status=%q payment=%q", repository.appliedStatus, repository.appliedPayment)
	}
}

func TestHandleNotificationRequiresValidSignatureAndMatchingSettlement(t *testing.T) {
	tests := []struct {
		name           string
		validSignature bool
		status         *entity.PaymentStatus
		wantError      error
		wantApplied    string
	}{
		{
			name:           "invalid signature",
			validSignature: false,
			status:         &entity.PaymentStatus{},
			wantError:      ErrInvalidNotification,
		},
		{
			name:           "mismatched amount",
			validSignature: true,
			status: &entity.PaymentStatus{
				OrderID:           "order-1",
				GrossAmount:       "25001.00",
				Currency:          "IDR",
				TransactionStatus: "settlement",
			},
			wantError: ErrUnexpectedPayment,
		},
		{
			name:           "settled payment",
			validSignature: true,
			status: &entity.PaymentStatus{
				OrderID:           "order-1",
				GrossAmount:       "25000.00",
				Currency:          "IDR",
				TransactionStatus: "settlement",
				PaymentType:       "bank_transfer",
			},
			wantApplied: "paid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &repositoryStub{
				topUp: &entity.TopUp{OrderID: "order-1", Amount: 25000},
			}
			gateway := paymentGatewayStub{validSignature: tt.validSignature, status: tt.status}
			usecase := NewUsecase(repository, gateway)

			err := usecase.HandleNotification("order-1", "200", "25000.00", "signature")
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("expected error %v, got %v", tt.wantError, err)
			}
			if tt.wantError == nil && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if repository.appliedStatus != tt.wantApplied {
				t.Errorf("expected applied status %q, got %q", tt.wantApplied, repository.appliedStatus)
			}
		})
	}
}
