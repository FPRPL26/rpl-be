package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	CreateSnapTransaction(transactionID string, amount int64, customerName, customerEmail string) (string, error)
	VerifySignatureKey(orderID, statusCode, grossAmount, signatureKey string) bool
}

type midtransService struct {
	serverKey string
}

func NewMidtransService() MidtransService {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	isProd := os.Getenv("MIDTRANS_IS_PROD") == "true"

	midtrans.ServerKey = serverKey
	if isProd {
		midtrans.Environment = midtrans.Production
	} else {
		midtrans.Environment = midtrans.Sandbox
	}

	return &midtransService{
		serverKey: serverKey,
	}
}

func (s *midtransService) CreateSnapTransaction(transactionID string, amount int64, customerName, customerEmail string) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionID,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
	}

	snapResp, err := snap.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return snapResp.RedirectURL, nil
}

func (s *midtransService) VerifySignatureKey(orderID, statusCode, grossAmount, signatureKey string) bool {
	payload := orderID + statusCode + grossAmount + s.serverKey
	h := sha512.New()
	h.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return expectedSignature == signatureKey
}
